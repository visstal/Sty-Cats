package repositories

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/domain/interfaces"
)

type MissionRepository struct {
	db *gorm.DB
}

func NewMissionRepository(db *gorm.DB) interfaces.MissionRepository {
	return &MissionRepository{db: db}
}

func (r *MissionRepository) WithTx(tx *gorm.DB) interfaces.MissionRepository {
	return &MissionRepository{db: tx}
}

func (r *MissionRepository) Create(mission *entities.Mission) (*entities.Mission, error) {
	if err := mission.Validate(); err != nil {
		return nil, err
	}

	missionCopy := *mission
	missionCopy.Targets = nil 

	if err := r.db.Create(&missionCopy).Error; err != nil {
		return nil, err
	}

	mission.ID = missionCopy.ID
	mission.CreatedAt = missionCopy.CreatedAt
	mission.UpdatedAt = missionCopy.UpdatedAt

	return mission, nil
}

func (r *MissionRepository) GetAll() ([]*entities.Mission, error) {
	var missions []*entities.Mission

	if err := r.db.
		Preload("Cat").
		Preload("Targets", func(db *gorm.DB) *gorm.DB {
			return db.Order("targets.created_at ASC") 
		}).
		Find(&missions).Error; err != nil {
		return nil, err
	}

	return missions, nil
}

func (r *MissionRepository) GetByID(id int32) (*entities.Mission, error) {
	var mission entities.Mission
	if err := r.db.
		Preload("Cat").
		Preload("Targets", func(db *gorm.DB) *gorm.DB {
			return db.Order("targets.created_at ASC")
		}).
		First(&mission, id).Error; err != nil {
		return nil, err
	}
	return &mission, nil
}

func (r *MissionRepository) Delete(id int32) error {
	var mission entities.Mission
	if err := r.db.Preload("Targets").First(&mission, id).Error; err != nil {
		return err
	}

	if mission.CatID != nil {
		return errors.New("cannot delete mission with assigned cat")
	}

	if err := r.db.Delete(&entities.Mission{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete mission: %w", err)
	}

	return nil
}

func (r *MissionRepository) CheckMissionExists(id int32) (bool, error) {
	var count int64
	if err := r.db.Model(&entities.Mission{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *MissionRepository) Update(mission *entities.Mission) (*entities.Mission, error) {
	if err := r.db.Save(mission).Error; err != nil {
		return nil, err
	}
	return mission, nil
}

func (r *MissionRepository) AssignCatToMission(missionID, catID int32) error {
	now := time.Now()
	if err := r.db.Model(&entities.Mission{}).
		Where("id = ?", missionID).
		Updates(map[string]interface{}{
			"cat_id":     catID,
			"start_date": now,
			"updated_at": now,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (r *MissionRepository) UnassignCatFromMission(catID int32) error {
	return errors.New("UnassignCatFromMission should be handled by cat repository")
}

func (r *MissionRepository) GetFreeCats() ([]*entities.SpyCat, error) {
	var cats []*entities.SpyCat
	if err := r.db.Where("mission_id IS NULL").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}
