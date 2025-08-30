package handlers

import (
	"net/http"
	"strconv"

	"spy-cat-agency/internal/application/dto"
	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/domain/interfaces"

	"github.com/labstack/echo/v4"
)

type TargetHandler struct {
	targetRepo  interfaces.TargetRepository
	missionRepo interfaces.MissionRepository
}

func NewTargetHandler(targetRepo interfaces.TargetRepository, missionRepo interfaces.MissionRepository) *TargetHandler {
	return &TargetHandler{
		targetRepo:  targetRepo,
		missionRepo: missionRepo,
	}
}

// AddTarget adds a new target to a mission
// @Summary Add a target to a mission
// @Description Add a new target to an existing mission (up to 3 targets total)
// @Tags targets
// @Accept json
// @Produce json
// @Param missionId path int true "Mission ID"
// @Param request body dto.AddTargetRequest true "Add target request"
// @Success 201 {object} dto.TargetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{missionId}/targets [post]
func (h *TargetHandler) AddTarget(c echo.Context) error {
	missionIDStr := c.Param("missionId")
	missionID, err := strconv.ParseInt(missionIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid mission ID",
			"details": "Mission ID must be a valid integer",
		})
	}

	var req dto.AddTargetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	mission, err := h.missionRepo.GetByID(int32(missionID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Mission not found",
			"details": "The specified mission does not exist",
		})
	}

	if !mission.CanAddTarget() {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Cannot add target",
			"details": "Mission already has the maximum number of targets (3)",
		})
	}

	for _, target := range mission.Targets {
		if target.Name == req.Name {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":   "Duplicate target name",
				"details": "A target with this name already exists in the mission",
			})
		}
	}

	target := req.ToTargetModel(int32(missionID))
	createdTarget, err := h.targetRepo.Create(c.Request().Context(), target)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to create target",
			"details": err.Error(),
		})
	}

	response := dto.TargetResponse{
		ID:        createdTarget.ID,
		MissionID: createdTarget.MissionID,
		Name:      createdTarget.Name,
		Country:   createdTarget.Country,
		Notes:     createdTarget.Notes,
		Status:    string(createdTarget.Status),
		CreatedAt: createdTarget.CreatedAt,
		UpdatedAt: createdTarget.UpdatedAt,
	}

	return c.JSON(http.StatusCreated, response)
}

// DeleteTarget deletes a target from a mission
// @Summary Delete a target from a mission
// @Description Delete a target from a mission (only if status is not 'init')
// @Tags targets
// @Param missionId path int true "Mission ID"
// @Param targetId path int true "Target ID"
// @Success 204 "Target deleted successfully"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{missionId}/targets/{targetId} [delete]
func (h *TargetHandler) DeleteTarget(c echo.Context) error {
	missionIDStr := c.Param("missionId")
	missionID, err := strconv.ParseInt(missionIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid mission ID",
			"details": "Mission ID must be a valid integer",
		})
	}

	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid target ID",
			"details": "Target ID must be a valid integer",
		})
	}

	target, err := h.targetRepo.GetByID(c.Request().Context(), int32(targetID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Target not found",
			"details": "The specified target does not exist",
		})
	}

	if target.MissionID != int32(missionID) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Target does not belong to mission",
			"details": "The target does not belong to the specified mission",
		})
	}

	if target.Status == entities.TargetStatusInit {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Cannot delete target",
			"details": "Targets in 'init' status cannot be deleted",
		})
	}

	mission, err := h.missionRepo.GetByID(int32(missionID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Mission not found",
			"details": "The specified mission does not exist",
		})
	}

	if len(mission.Targets) <= entities.MinTargetsRequired {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Cannot delete target",
			"details": "Mission must have at least one target",
		})
	}

	if err := h.targetRepo.Delete(c.Request().Context(), int32(targetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to delete target",
			"details": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// UpdateTarget updates a target
// @Summary Update a target
// @Description Update target information
// @Tags targets
// @Accept json
// @Produce json
// @Param missionId path int true "Mission ID"
// @Param targetId path int true "Target ID"
// @Param request body dto.UpdateTargetRequest true "Update target request"
// @Success 200 {object} dto.TargetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{missionId}/targets/{targetId} [put]
func (h *TargetHandler) UpdateTarget(c echo.Context) error {
	missionIDStr := c.Param("missionId")
	missionID, err := strconv.ParseInt(missionIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid mission ID",
			"details": "Mission ID must be a valid integer",
		})
	}

	targetIDStr := c.Param("targetId")
	targetID, err := strconv.ParseInt(targetIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid target ID",
			"details": "Target ID must be a valid integer",
		})
	}

	var req dto.UpdateTargetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	target, err := h.targetRepo.GetByID(c.Request().Context(), int32(targetID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"error":   "Target not found",
			"details": "The specified target does not exist",
		})
	}

	if target.MissionID != int32(missionID) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Target does not belong to mission",
			"details": "The target does not belong to the specified mission",
		})
	}

	if req.Name != nil {
		target.Name = *req.Name
	}
	if req.Country != nil {
		target.Country = *req.Country
	}
	if req.Notes != nil {
		target.Notes = req.Notes
	}
	if req.Status != nil {
		target.Status = entities.TargetStatus(*req.Status)
	}

	updatedTarget, err := h.targetRepo.Update(c.Request().Context(), target)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to update target",
			"details": err.Error(),
		})
	}

	response := dto.TargetResponse{
		ID:        updatedTarget.ID,
		MissionID: updatedTarget.MissionID,
		Name:      updatedTarget.Name,
		Country:   updatedTarget.Country,
		Notes:     updatedTarget.Notes,
		Status:    string(updatedTarget.Status),
		CreatedAt: updatedTarget.CreatedAt,
		UpdatedAt: updatedTarget.UpdatedAt,
	}

	return c.JSON(http.StatusOK, response)
}
