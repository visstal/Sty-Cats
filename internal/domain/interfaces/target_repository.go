package interfaces

import (
	"context"
	"spy-cat-agency/internal/domain/entities"

	"gorm.io/gorm"
)

type TargetRepository interface {
	Create(ctx context.Context, target *entities.Target) (*entities.Target, error)
	CreateMany(ctx context.Context, targets []*entities.Target) error
	GetByMissionID(ctx context.Context, missionID int32) ([]*entities.Target, error)
	GetByID(ctx context.Context, id int32) (*entities.Target, error)
	Update(ctx context.Context, target *entities.Target) (*entities.Target, error)
	Delete(ctx context.Context, id int32) error
	DeleteByMissionID(ctx context.Context, missionID int32) error
	UpdateStatus(ctx context.Context, id int32, status entities.TargetStatus) error
	UpdateNotes(ctx context.Context, id int32, notes string) error
	WithTx(tx *gorm.DB) TargetRepository
}
