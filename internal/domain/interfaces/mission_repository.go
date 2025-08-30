package interfaces

import (
	"spy-cat-agency/internal/domain/entities"

	"gorm.io/gorm"
)

type MissionRepository interface {
	Create(mission *entities.Mission) (*entities.Mission, error)
	GetAll() ([]*entities.Mission, error)
	GetByID(id int32) (*entities.Mission, error)
	Delete(id int32) error
	CheckMissionExists(id int32) (bool, error)
	Update(mission *entities.Mission) (*entities.Mission, error)
	AssignCatToMission(missionID, catID int32) error
	UnassignCatFromMission(catID int32) error
	GetFreeCats() ([]*entities.SpyCat, error)
	WithTx(tx *gorm.DB) MissionRepository
}
