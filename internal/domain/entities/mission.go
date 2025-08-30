package entities

import (
	"errors"
	"time"
)

const (
	MinTargetsRequired = 1
	MaxTargetsAllowed  = 3
)

type Mission struct {
	ID          int32      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string     `json:"name" gorm:"not null;size:100"`
	Description string     `json:"description" gorm:"not null;size:500"`
	StartDate   time.Time  `json:"start_date" gorm:"not null"`
	EndDate     time.Time  `json:"end_date" gorm:"not null"`
	CatID       *int32     `json:"cat_id" gorm:"index"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	Cat     *SpyCat  `json:"cat,omitempty" gorm:"foreignKey:CatID;references:ID"`
	Targets []Target `json:"targets,omitempty" gorm:"foreignKey:MissionID"`
}

func (m *Mission) Validate() error {
	if err := m.ValidateTargets(); err != nil {
		return err
	}

	if m.Name == "" {
		return errors.New("mission name is required")
	}

	if m.Description == "" {
		return errors.New("mission description is required")
	}

	return nil
}

func (m *Mission) ValidateTargets() error {
	targetCount := len(m.Targets)

	if targetCount < MinTargetsRequired {
		return errors.New("mission must have at least one target")
	}

	if targetCount > MaxTargetsAllowed {
		return errors.New("mission cannot have more than 3 targets")
	}

	for i, target := range m.Targets {
		if target.Name == "" {
			return errors.New("target name is required")
		}
		if target.Country == "" {
			return errors.New("target country is required")
		}
		if len(target.Name) > 100 {
			return errors.New("target name cannot exceed 100 characters")
		}
		if len(target.Country) > 100 {
			return errors.New("target country cannot exceed 100 characters")
		}

		for j := i + 1; j < len(m.Targets); j++ {
			if m.Targets[j].Name == target.Name {
				return errors.New("duplicate target names are not allowed in the same mission")
			}
		}
	}

	return nil
}

func (m *Mission) CanAddTarget() bool {
	return len(m.Targets) < MaxTargetsAllowed
}

func (m *Mission) HasMinimumTargets() bool {
	return len(m.Targets) >= MinTargetsRequired
}

func (Mission) TableName() string {
	return "missions"
}
