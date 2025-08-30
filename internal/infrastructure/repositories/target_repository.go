package repositories

import (
	"context"
	"gorm.io/gorm"

	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/domain/interfaces"
)

type TargetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) interfaces.TargetRepository {
	return &TargetRepository{db: db}
}

func (r *TargetRepository) WithTx(tx *gorm.DB) interfaces.TargetRepository {
	return &TargetRepository{db: tx}
}

func (r *TargetRepository) Create(ctx context.Context, target *entities.Target) (*entities.Target, error) {
	if err := r.db.WithContext(ctx).Create(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (r *TargetRepository) CreateMany(ctx context.Context, targets []*entities.Target) error {
	for _, target := range targets {
		if err := r.db.WithContext(ctx).Create(target).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *TargetRepository) GetByMissionID(ctx context.Context, missionID int32) ([]*entities.Target, error) {
	var targets []*entities.Target
	if err := r.db.WithContext(ctx).Where("mission_id = ?", missionID).Find(&targets).Error; err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *TargetRepository) GetByID(ctx context.Context, id int32) (*entities.Target, error) {
	var target entities.Target
	if err := r.db.WithContext(ctx).First(&target, id).Error; err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) Update(ctx context.Context, target *entities.Target) (*entities.Target, error) {
	if err := r.db.WithContext(ctx).Save(target).Error; err != nil {
		return nil, err
	}
	return target, nil
}

func (r *TargetRepository) Delete(ctx context.Context, id int32) error {
	return r.db.WithContext(ctx).Delete(&entities.Target{}, id).Error
}

func (r *TargetRepository) DeleteByMissionID(ctx context.Context, missionID int32) error {
	return r.db.WithContext(ctx).Unscoped().Where("mission_id = ?", missionID).Delete(&entities.Target{}).Error
}

func (r *TargetRepository) UpdateStatus(ctx context.Context, id int32, status entities.TargetStatus) error {
	return r.db.WithContext(ctx).Model(&entities.Target{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TargetRepository) UpdateNotes(ctx context.Context, id int32, notes string) error {
	return r.db.WithContext(ctx).Model(&entities.Target{}).Where("id = ?", id).Update("notes", notes).Error
}
