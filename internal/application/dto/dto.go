package dto

import (
	"time"

	"spy-cat-agency/internal/domain/entities"
)

type CreateCatRequest struct {
	Name              string  `json:"name" validate:"required,min=1,max=100"`
	YearsOfExperience int32   `json:"years_of_experience" validate:"required,min=0"`
	Breed             string  `json:"breed" validate:"required,min=1,max=100"`
	Salary            float64 `json:"salary" validate:"required,min=0"`
}

type UpdateCatSalaryRequest struct {
	Salary float64 `json:"salary" validate:"required,min=0"`
}

type CatResponse struct {
	ID                int32     `json:"id"`
	Name              string    `json:"name"`
	YearsOfExperience int32     `json:"years_of_experience"`
	Breed             string    `json:"breed"`
	Salary            float64   `json:"salary"`
	MissionID         *int32    `json:"mission_id,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CatListResponse struct {
	Cats   []CatResponse `json:"cats"`
	Breeds []string      `json:"breeds"`
	Total  int64         `json:"total"`
	Limit  int32         `json:"limit"`
	Offset int32         `json:"offset"`
}

type CreateTargetRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=100"`
	Country string `json:"country" validate:"required,min=1,max=100"`
}

type CreateMissionRequest struct {
	Name        string                `json:"name" validate:"required,min=1,max=100"`
	Description string                `json:"description" validate:"required,min=1,max=500"`
	StartDate   *time.Time            `json:"start_date,omitempty"`
	EndDate     *time.Time            `json:"end_date,omitempty"`
	Targets     []CreateTargetRequest `json:"targets" validate:"required,min=1,max=3,dive"`
}

func (r *CreateMissionRequest) ToModel() *entities.Mission {
	mission := &entities.Mission{
		Name:        r.Name,
		Description: r.Description,
		IsCompleted: false,
	}

	if r.StartDate != nil {
		mission.StartDate = *r.StartDate
	}

	if r.EndDate != nil {
		mission.EndDate = *r.EndDate
	}

	if len(r.Targets) > 0 {
		mission.Targets = make([]entities.Target, len(r.Targets))
		for i, targetReq := range r.Targets {
			mission.Targets[i] = entities.Target{
				Name:    targetReq.Name,
				Country: targetReq.Country,
				Status:  entities.TargetStatusInit,
			}
		}
	}

	return mission
}

type TargetResponse struct {
	ID        int32     `json:"id"`
	MissionID int32     `json:"mission_id"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Notes     *string   `json:"notes"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MissionResponse struct {
	ID          int32            `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	StartDate   *time.Time       `json:"start_date,omitempty"`
	EndDate     *time.Time       `json:"end_date,omitempty"`
	CatID       *int32           `json:"cat_id"`
	IsCompleted bool             `json:"is_completed"`
	CompletedAt *time.Time       `json:"completed_at,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Cat         *CatResponse     `json:"cat,omitempty"`
	Targets     []TargetResponse `json:"targets,omitempty"`
}

func MissionFromModel(mission *entities.Mission) *MissionResponse {
	response := &MissionResponse{
		ID:          mission.ID,
		Name:        mission.Name,
		Description: mission.Description,
		CatID:       mission.CatID,
		IsCompleted: mission.IsCompleted,
		CompletedAt: mission.CompletedAt,
		CreatedAt:   mission.CreatedAt,
		UpdatedAt:   mission.UpdatedAt,
	}

	if !mission.StartDate.IsZero() {
		response.StartDate = &mission.StartDate
	}

	if !mission.EndDate.IsZero() {
		response.EndDate = &mission.EndDate
	}

	if mission.Cat != nil {
		response.Cat = &CatResponse{
			ID:                mission.Cat.ID,
			Name:              mission.Cat.Name,
			YearsOfExperience: mission.Cat.YearsOfExperience,
			Breed:             mission.Cat.Breed,
			Salary:            mission.Cat.Salary,
			CreatedAt:         mission.Cat.CreatedAt,
			UpdatedAt:         mission.Cat.UpdatedAt,
		}
	}

	if len(mission.Targets) > 0 {
		response.Targets = make([]TargetResponse, len(mission.Targets))
		for i, target := range mission.Targets {
			response.Targets[i] = TargetResponse{
				ID:        target.ID,
				MissionID: target.MissionID,
				Name:      target.Name,
				Country:   target.Country,
				Notes:     target.Notes,
				Status:    string(target.Status),
				CreatedAt: target.CreatedAt,
				UpdatedAt: target.UpdatedAt,
			}
		}
	}

	return response
}

type MissionListResponse struct {
	Missions []MissionResponse `json:"missions"`
	Total    int64             `json:"total"`
	Limit    int32             `json:"limit"`
	Offset   int32             `json:"offset"`
}

type AddTargetRequest struct {
	Name    string  `json:"name" validate:"required,min=1,max=100"`
	Country string  `json:"country" validate:"required,min=1,max=100"`
	Notes   *string `json:"notes,omitempty"`
}

type UpdateTargetRequest struct {
	Name    *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Country *string `json:"country,omitempty" validate:"omitempty,min=1,max=100"`
	Notes   *string `json:"notes,omitempty"`
	Status  *string `json:"status,omitempty" validate:"omitempty,oneof=init in_progress completed"`
}

func (r *AddTargetRequest) ToTargetModel(missionID int32) *entities.Target {
	return &entities.Target{
		MissionID: missionID,
		Name:      r.Name,
		Country:   r.Country,
		Notes:     r.Notes,
		Status:    entities.TargetStatusInit,
	}
}

type AssignCatRequest struct {
	CatID int32 `json:"cat_id" validate:"required,min=1"`
}
