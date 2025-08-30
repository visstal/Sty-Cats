package entities

import (
	"time"
)

type SpyCat struct {
	ID                int32     `gorm:"primaryKey;autoIncrement"`
	Name              string    `gorm:"type:varchar(100);not null"`
	YearsOfExperience int32     `gorm:"not null;check:years_of_experience >= 0"`
	Breed             string    `gorm:"type:varchar(100);not null"`
	Salary            float64   `gorm:"type:numeric(12,2);not null;check:salary >= 0"`
	MissionID         *int32    `gorm:"type:integer;default:null;index"`
	CreatedAt         time.Time `gorm:"not null;default:now()"`
	UpdatedAt         time.Time `gorm:"not null;default:now()"`
}

func (SpyCat) TableName() string {
	return "spy_cats"
}

func NewSpyCat(name, breed string, yearsOfExperience int32, salary float64) *SpyCat {
	return &SpyCat{
		Name:              name,
		YearsOfExperience: yearsOfExperience,
		Breed:             breed,
		Salary:            salary,
	}
}

func (c *SpyCat) UpdateSalary(newSalary float64) {
	c.Salary = newSalary
	c.UpdatedAt = time.Now()
}
