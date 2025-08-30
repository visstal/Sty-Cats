package interfaces

import (
	"context"
	"spy-cat-agency/internal/domain/entities"
	"gorm.io/gorm"
)

type CatRepository interface {
	Create(ctx context.Context, cat *entities.SpyCat) (*entities.SpyCat, error)
	GetByID(ctx context.Context, id int32) (*entities.SpyCat, error)
	List(ctx context.Context, limit, offset int32) ([]*entities.SpyCat, error)
	UpdateSalary(ctx context.Context, id int32, salary float64) (*entities.SpyCat, error)
	Delete(ctx context.Context, id int32) error
	UnassignFromMission(ctx context.Context, catID int32) error
	AssignToMission(ctx context.Context, catID, missionID int32) error
	WithTx(tx *gorm.DB) CatRepository
}
