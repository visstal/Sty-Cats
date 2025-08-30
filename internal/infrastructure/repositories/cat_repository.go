package repositories

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/domain/interfaces"
	"spy-cat-agency/internal/infrastructure/database"
)

type CatRepository struct {
	db *database.DB
}

func NewCatRepository(db *database.DB) interfaces.CatRepository {
	return &CatRepository{
		db: db,
	}
}

func (r *CatRepository) WithTx(tx *gorm.DB) interfaces.CatRepository {
	return &CatRepository{
		db: &database.DB{DB: tx},
	}
}

func (r *CatRepository) Create(ctx context.Context, spyCat *entities.SpyCat) (*entities.SpyCat, error) {
	if err := r.db.WithContext(ctx).Create(spyCat).Error; err != nil {
		return nil, fmt.Errorf("failed to create cat: %w", err)
	}
	return spyCat, nil
}

func (r *CatRepository) GetByID(ctx context.Context, id int32) (*entities.SpyCat, error) {
	var spyCat entities.SpyCat
	if err := r.db.WithContext(ctx).First(&spyCat, id).Error; err != nil {
		return nil, fmt.Errorf("failed to get cat: %w", err)
	}
	return &spyCat, nil
}

func (r *CatRepository) List(ctx context.Context, limit, offset int32) ([]*entities.SpyCat, error) {
	var spyCats []*entities.SpyCat
	if err := r.db.WithContext(ctx).
		Order("id").
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&spyCats).Error; err != nil {
		return nil, fmt.Errorf("failed to list cats: %w", err)
	}
	return spyCats, nil
}

func (r *CatRepository) UpdateSalary(ctx context.Context, id int32, salary float64) (*entities.SpyCat, error) {
	var spyCat entities.SpyCat
	if err := r.db.WithContext(ctx).First(&spyCat, id).Error; err != nil {
		return nil, fmt.Errorf("failed to find cat: %w", err)
	}

	spyCat.UpdateSalary(salary)
	if err := r.db.WithContext(ctx).Save(&spyCat).Error; err != nil {
		return nil, fmt.Errorf("failed to update cat salary: %w", err)
	}

	return &spyCat, nil
}

func (r *CatRepository) Delete(ctx context.Context, id int32) error {
	var spyCat entities.SpyCat
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&spyCat).Error; err != nil {
		return fmt.Errorf("failed to find cat: %w", err)
	}

	if spyCat.MissionID != nil {
		return fmt.Errorf("cannot delete cat: cat is currently assigned to mission ID %d", *spyCat.MissionID)
	}

	if err := r.db.WithContext(ctx).Delete(&entities.SpyCat{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete cat: %w", err)
	}
	return nil
}

func (r *CatRepository) AssignToMission(ctx context.Context, catID, missionID int32) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&entities.SpyCat{}).
		Where("id = ?", catID).
		Updates(map[string]interface{}{
			"mission_id": missionID,
			"updated_at": now,
		}).Error; err != nil {
		return fmt.Errorf("failed to assign cat to mission: %w", err)
	}
	return nil
}

func (r *CatRepository) UnassignFromMission(ctx context.Context, catID int32) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&entities.SpyCat{}).
		Where("id = ?", catID).
		Updates(map[string]interface{}{
			"mission_id": nil,
			"updated_at": now,
		}).Error; err != nil {
		return fmt.Errorf("failed to unassign cat from mission: %w", err)
	}
	return nil
}
