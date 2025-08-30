package services

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"spy-cat-agency/internal/application/dto"
	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/infrastructure/database"
	"spy-cat-agency/internal/infrastructure/repositories"
)

type MissionService interface {
	CreateMission(req dto.CreateMissionRequest) (*dto.MissionResponse, error)
	ListMissions() ([]*dto.MissionResponse, error)
	GetMission(id int32) (*dto.MissionResponse, error)
	DeleteMission(id int32) error
	AssignCatToMission(missionID, catID int32) (*dto.MissionResponse, error)
	GetFreeCats() ([]*dto.CatResponse, error)
	AddTargetToMission(missionID int32, req dto.AddTargetRequest) (*dto.TargetResponse, error)
	DeleteTargetFromMission(missionID, targetID int32) error
	GetCatMission(catID int32) (*dto.MissionResponse, error)
	UpdateTargetStatus(catID, targetID int32, status string) (*dto.TargetResponse, error)
	UpdateTargetNotes(catID, targetID int32, notes string) (*dto.TargetResponse, error)
}

type missionService struct {
	db          *database.DB
	missionRepo *repositories.MissionRepository
	targetRepo  *repositories.TargetRepository
	catRepo     *repositories.CatRepository
}

func NewMissionService(db *database.DB, missionRepo *repositories.MissionRepository, targetRepo *repositories.TargetRepository, catRepo *repositories.CatRepository) MissionService {
	return &missionService{
		db:          db,
		missionRepo: missionRepo,
		targetRepo:  targetRepo,
		catRepo:     catRepo,
	}
}

func (s *missionService) CreateMission(req dto.CreateMissionRequest) (*dto.MissionResponse, error) {
	if err := validateCreateMissionRequest(req); err != nil {
		return nil, err
	}

	mission := req.ToModel()

	var createdMission *entities.Mission

	err := s.db.RunTransaction(func(tx *gorm.DB) error {
		txMissionRepo := s.missionRepo.WithTx(tx)
		var err error
		createdMission, err = txMissionRepo.Create(mission)
		if err != nil {
			return fmt.Errorf("failed to create mission: %w", err)
		}

		if len(mission.Targets) > 0 {
			txTargetRepo := s.targetRepo.WithTx(tx)
			for i := range mission.Targets {
				mission.Targets[i].MissionID = createdMission.ID
				mission.Targets[i].Status = entities.TargetStatusInit
				mission.Targets[i].ID = 0

				_, err := txTargetRepo.Create(context.TODO(), &mission.Targets[i])
				if err != nil {
					return fmt.Errorf("failed to create target: %w", err)
				}
			}
			createdMission.Targets = mission.Targets
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return dto.MissionFromModel(createdMission), nil
}

func (s *missionService) ListMissions() ([]*dto.MissionResponse, error) {
	missions, err := s.missionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to list missions: %w", err)
	}

	responses := make([]*dto.MissionResponse, len(missions))
	for i, mission := range missions {
		responses[i] = dto.MissionFromModel(mission)
	}

	return responses, nil
}

func (s *missionService) GetMission(id int32) (*dto.MissionResponse, error) {
	mission, err := s.missionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get mission: %w", err)
	}

	return dto.MissionFromModel(mission), nil
}

func (s *missionService) DeleteMission(id int32) error {
	exists, err := s.missionRepo.CheckMissionExists(id)
	if err != nil {
		return fmt.Errorf("failed to check mission existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("mission with id %d not found", id)
	}

	return s.db.RunTransaction(func(tx *gorm.DB) error {
		txTargetRepo := s.targetRepo.WithTx(tx)
		if err := txTargetRepo.DeleteByMissionID(context.TODO(), id); err != nil {
			return fmt.Errorf("failed to delete mission targets: %w", err)
		}

		txMissionRepo := s.missionRepo.WithTx(tx)
		if err := txMissionRepo.Delete(id); err != nil {
			return fmt.Errorf("failed to delete mission: %w", err)
		}

		return nil
	})
}

func validateCreateMissionRequest(req dto.CreateMissionRequest) error {
	if err := validateMissionName(req.Name); err != nil {
		return err
	}

	if err := validateMissionDescription(req.Description); err != nil {
		return err
	}

	if req.StartDate != nil && req.EndDate != nil {
		if err := validateMissionDates(*req.StartDate, *req.EndDate); err != nil {
			return err
		}
	}

	return nil
}

func (s *missionService) AssignCatToMission(missionID, catID int32) (*dto.MissionResponse, error) {
	err := s.db.RunTransaction(func(tx *gorm.DB) error {
		ctx := context.Background()

		txMissionRepo := s.missionRepo.WithTx(tx)
		if err := txMissionRepo.AssignCatToMission(missionID, catID); err != nil {
			return fmt.Errorf("failed to assign cat in mission table: %w", err)
		}

		txCatRepo := s.catRepo.WithTx(tx)
		if err := txCatRepo.AssignToMission(ctx, catID, missionID); err != nil {
			return fmt.Errorf("failed to assign mission to cat: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated mission: %w", err)
	}

	return dto.MissionFromModel(mission), nil
}

func (s *missionService) GetFreeCats() ([]*dto.CatResponse, error) {
	cats, err := s.missionRepo.GetFreeCats()
	if err != nil {
		return nil, fmt.Errorf("failed to get free cats: %w", err)
	}

	var responses []*dto.CatResponse
	for _, cat := range cats {
		responses = append(responses, &dto.CatResponse{
			ID:                cat.ID,
			Name:              cat.Name,
			YearsOfExperience: cat.YearsOfExperience,
			Breed:             cat.Breed,
			Salary:            cat.Salary,
			MissionID:         cat.MissionID,
			CreatedAt:         cat.CreatedAt,
			UpdatedAt:         cat.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *missionService) AddTargetToMission(missionID int32, req dto.AddTargetRequest) (*dto.TargetResponse, error) {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if !mission.CanAddTarget() {
		return nil, fmt.Errorf("mission already has the maximum number of targets (%d)", len(mission.Targets))
	}

	for _, target := range mission.Targets {
		if target.Name == req.Name {
			return nil, fmt.Errorf("target with name '%s' already exists in this mission", req.Name)
		}
	}

	target := req.ToTargetModel(missionID)
	createdTarget, err := s.targetRepo.Create(context.TODO(), target)
	if err != nil {
		return nil, fmt.Errorf("failed to create target: %w", err)
	}

	return &dto.TargetResponse{
		ID:        createdTarget.ID,
		MissionID: createdTarget.MissionID,
		Name:      createdTarget.Name,
		Country:   createdTarget.Country,
		Notes:     createdTarget.Notes,
		Status:    string(createdTarget.Status),
		CreatedAt: createdTarget.CreatedAt,
		UpdatedAt: createdTarget.UpdatedAt,
	}, nil
}

func (s *missionService) DeleteTargetFromMission(missionID, targetID int32) error {
	target, err := s.targetRepo.GetByID(context.TODO(), targetID)
	if err != nil {
		return fmt.Errorf("target not found: %w", err)
	}

	if target.MissionID != missionID {
		return fmt.Errorf("target does not belong to the specified mission")
	}

	if target.Status != "init" {
		return fmt.Errorf("only targets in 'init' status can be deleted")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("mission not found: %w", err)
	}

	if !mission.HasMinimumTargets() || len(mission.Targets) <= 1 {
		return fmt.Errorf("mission must have at least one target")
	}

	if err := s.targetRepo.Delete(context.TODO(), targetID); err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}

	return nil
}

func (s *missionService) GetCatMission(catID int32) (*dto.MissionResponse, error) {
	missions, err := s.missionRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get missions: %w", err)
	}

	for _, mission := range missions {
		if mission.CatID != nil && *mission.CatID == catID {
			return dto.MissionFromModel(mission), nil
		}
	}

	return nil, nil
}

func (s *missionService) UpdateTargetStatus(catID, targetID int32, status string) (*dto.TargetResponse, error) {
	target, err := s.targetRepo.GetByID(context.TODO(), targetID)
	if err != nil {
		return nil, fmt.Errorf("target not found: %w", err)
	}

	mission, err := s.missionRepo.GetByID(target.MissionID)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if mission.CatID == nil || *mission.CatID != catID {
		return nil, fmt.Errorf("target does not belong to the cat's mission")
	}

	if target.Status == "completed" {
		return nil, fmt.Errorf("target status is final and cannot be changed")
	}

	if err := s.targetRepo.UpdateStatus(context.TODO(), targetID, entities.TargetStatus(status)); err != nil {
		return nil, fmt.Errorf("failed to update target status: %w", err)
	}

	updatedTarget, err := s.targetRepo.GetByID(context.TODO(), targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated target: %w", err)
	}

	if status == "completed" {
		if err := s.checkAndCompleteMission(target.MissionID); err != nil {
			fmt.Printf("Error checking mission completion: %v\n", err)
		}
	}

	return &dto.TargetResponse{
		ID:        updatedTarget.ID,
		MissionID: updatedTarget.MissionID,
		Name:      updatedTarget.Name,
		Country:   updatedTarget.Country,
		Notes:     updatedTarget.Notes,
		Status:    string(updatedTarget.Status),
		CreatedAt: updatedTarget.CreatedAt,
		UpdatedAt: updatedTarget.UpdatedAt,
	}, nil
}

func (s *missionService) UpdateTargetNotes(catID, targetID int32, notes string) (*dto.TargetResponse, error) {
	target, err := s.targetRepo.GetByID(context.TODO(), targetID)
	if err != nil {
		return nil, fmt.Errorf("target not found: %w", err)
	}

	mission, err := s.missionRepo.GetByID(target.MissionID)
	if err != nil {
		return nil, fmt.Errorf("mission not found: %w", err)
	}

	if mission.CatID == nil || *mission.CatID != catID {
		return nil, fmt.Errorf("target does not belong to the cat's mission")
	}

	if target.Status == "completed" {
		return nil, fmt.Errorf("target is final and cannot be modified")
	}

	if err := s.targetRepo.UpdateNotes(context.TODO(), targetID, notes); err != nil {
		return nil, fmt.Errorf("failed to update target notes: %w", err)
	}

	updatedTarget, err := s.targetRepo.GetByID(context.TODO(), targetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated target: %w", err)
	}

	return &dto.TargetResponse{
		ID:        updatedTarget.ID,
		MissionID: updatedTarget.MissionID,
		Name:      updatedTarget.Name,
		Country:   updatedTarget.Country,
		Notes:     updatedTarget.Notes,
		Status:    string(updatedTarget.Status),
		CreatedAt: updatedTarget.CreatedAt,
		UpdatedAt: updatedTarget.UpdatedAt,
	}, nil
}

func (s *missionService) checkAndCompleteMission(missionID int32) error {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return fmt.Errorf("failed to get mission: %w", err)
	}

	allCompleted := true
	for _, target := range mission.Targets {
		if target.Status != "completed" {
			allCompleted = false
			break
		}
	}

	if allCompleted && !mission.IsCompleted {
		now := time.Now()

		mission.IsCompleted = true
		mission.CompletedAt = &now

		if mission.EndDate.IsZero() {
			mission.EndDate = now
		}

		assignedCatID := mission.CatID

		mission.CatID = nil

		_, err := s.missionRepo.Update(mission)
		if err != nil {
			return fmt.Errorf("failed to complete mission: %w", err)
		}

		if assignedCatID != nil {
			ctx := context.Background()
			if err := s.catRepo.UnassignFromMission(ctx, *assignedCatID); err != nil {
				return fmt.Errorf("failed to unassign cat from completed mission: %w", err)
			}
		}
	}

	return nil
}
