package entities

import (
	"time"

	"gorm.io/gorm"
)

type TargetStatus string

const (
	TargetStatusInit       TargetStatus = "init"
	TargetStatusInProgress TargetStatus = "in_progress"
	TargetStatusCompleted  TargetStatus = "completed"
)

type Target struct {
	ID        int32          `gorm:"primaryKey;autoIncrement" json:"id"`
	MissionID int32          `gorm:"not null;index" json:"mission_id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Country   string         `gorm:"size:100;not null" json:"country"`
	Notes     *string        `gorm:"type:text" json:"notes"`
	Status    TargetStatus   `gorm:"size:20;not null;default:'init'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Target) TableName() string {
	return "targets"
}
