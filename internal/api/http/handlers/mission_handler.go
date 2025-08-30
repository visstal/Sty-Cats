package handlers

import (
	"net/http"
	"strconv"

	"spy-cat-agency/internal/application/dto"
	"spy-cat-agency/internal/application/services"

	"github.com/labstack/echo/v4"
)

type MissionHandler struct {
	missionService services.MissionService
}

func NewMissionHandler(missionService services.MissionService) *MissionHandler {
	return &MissionHandler{
		missionService: missionService,
	}
}

// CreateMission creates a new mission
// @Summary Create a new mission
// @Description Create a new spy mission
// @Tags missions
// @Accept json
// @Produce json
// @Param mission body dto.CreateMissionRequest true "Mission data"
// @Success 201 {object} dto.MissionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions [post]
func (h *MissionHandler) CreateMission(c echo.Context) error {
	var req dto.CreateMissionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body format",
			"details": "Please check your JSON format and field types",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	mission, err := h.missionService.CreateMission(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Failed to create mission",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, mission)
}

// ListMissions returns all missions
// @Summary List all missions
// @Description Get all spy missions
// @Tags missions
// @Produce json
// @Success 200 {array} dto.MissionResponse
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions [get]
func (h *MissionHandler) ListMissions(c echo.Context) error {
	missions, err := h.missionService.ListMissions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch missions",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, missions)
}

// GetMission returns a specific mission
// @Summary Get mission by ID
// @Description Get a specific spy mission by its ID
// @Tags missions
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} dto.MissionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{id} [get]
func (h *MissionHandler) GetMission(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid mission ID",
		})
	}

	mission, err := h.missionService.GetMission(int32(id))
	if err != nil {
		if err.Error() == "failed to get mission: record not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Mission not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to fetch mission",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mission)
}

// DeleteMission deletes a mission
// @Summary Delete mission
// @Description Delete a spy mission by its ID
// @Tags missions
// @Param id path int true "Mission ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{id} [delete]
func (h *MissionHandler) DeleteMission(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid mission ID",
		})
	}

	if err := h.missionService.DeleteMission(int32(id)); err != nil {
		if err.Error() == "mission with id "+strconv.Itoa(int(id))+" not found" {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Mission not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to delete mission",
			"details": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// AssignCatToMission assigns a cat to a mission
// @Summary Assign a cat to a mission
// @Description Assign a spy cat to a mission and set start date
// @Tags missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param cat_id body object{cat_id:int} true "Cat assignment data"
// @Success 200 {object} dto.MissionResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /missions/{id}/assign [post]
func (h *MissionHandler) AssignCatToMission(c echo.Context) error {
	missionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid mission ID",
			"details": "Mission ID must be a positive integer",
		})
	}

	var req dto.AssignCatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body format",
			"details": "Please check your JSON format and field types",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	mission, err := h.missionService.AssignCatToMission(int32(missionID), req.CatID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to assign cat to mission",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, mission)
}

// GetFreeCats returns all cats that are not assigned to any mission
// @Summary Get free cats
// @Description Get all cats that are available for assignment
// @Tags cats
// @Produce json
// @Success 200 {array} dto.CatResponse
// @Failure 500 {object} map[string]interface{}
// @Router /missions/free-cats [get]
func (h *MissionHandler) GetFreeCats(c echo.Context) error {
	cats, err := h.missionService.GetFreeCats()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get free cats",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cats)
}

// AddTargetToMission adds a new target to a mission
// @Summary Add a target to a mission
// @Description Add a new target to an existing mission (up to 3 targets total)
// @Tags missions
// @Accept json
// @Produce json
// @Param missionId path int true "Mission ID"
// @Param request body dto.AddTargetRequest true "Add target request"
// @Success 201 {object} dto.TargetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{missionId}/targets [post]
func (h *MissionHandler) AddTargetToMission(c echo.Context) error {
	missionIDStr := c.Param("id")
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

	target, err := h.missionService.AddTargetToMission(int32(missionID), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Failed to add target",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, target)
}

// DeleteTargetFromMission deletes a target from a mission
// @Summary Delete a target from a mission
// @Description Delete a target from a mission (only if status is 'init')
// @Tags missions
// @Param missionId path int true "Mission ID"
// @Param targetId path int true "Target ID"
// @Success 204 "Target deleted successfully"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/agency/missions/{missionId}/targets/{targetId} [delete]
func (h *MissionHandler) DeleteTargetFromMission(c echo.Context) error {
	missionIDStr := c.Param("id")
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

	err = h.missionService.DeleteTargetFromMission(int32(missionID), int32(targetID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Failed to delete target",
			"details": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetCatMission gets the current mission for a specific cat
// @Summary Get cat's current mission
// @Description Get the current mission assigned to a specific cat with all targets
// @Tags spy-cats
// @Produce json
// @Param catId path int true "Cat ID"
// @Success 200 {object} dto.MissionResponse
// @Success 204 "Cat has no assigned mission"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/spy-cats/{catId}/mission [get]
func (h *MissionHandler) GetCatMission(c echo.Context) error {
	catIDStr := c.Param("catId")
	catID, err := strconv.ParseInt(catIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid cat ID",
			"details": "Cat ID must be a valid integer",
		})
	}

	mission, err := h.missionService.GetCatMission(int32(catID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error":   "Failed to get cat mission",
			"details": err.Error(),
		})
	}

	if mission == nil {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, mission)
}

// UpdateTargetStatus updates the status of a target (spy cat functionality)
// @Summary Update target status
// @Description Update the status of a target (only by assigned cat)
// @Tags spy-cats
// @Accept json
// @Produce json
// @Param catId path int true "Cat ID"
// @Param targetId path int true "Target ID"
// @Param request body object{status:string} true "Status update request"
// @Success 200 {object} dto.TargetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/spy-cats/{catId}/mission/targets/{targetId}/status [put]
func (h *MissionHandler) UpdateTargetStatus(c echo.Context) error {
	catIDStr := c.Param("catId")
	catID, err := strconv.ParseInt(catIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid cat ID",
			"details": "Cat ID must be a valid integer",
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

	var req struct {
		Status string `json:"status" validate:"required,oneof=init in_progress completed"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	target, err := h.missionService.UpdateTargetStatus(int32(catID), int32(targetID), req.Status)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Failed to update target status",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, target)
}

// UpdateTargetNotes updates the notes of a target (spy cat functionality)
// @Summary Update target notes
// @Description Update the notes of a target (only by assigned cat)
// @Tags spy-cats
// @Accept json
// @Produce json
// @Param catId path int true "Cat ID"
// @Param targetId path int true "Target ID"
// @Param request body object{notes:string} true "Notes update request"
// @Success 200 {object} dto.TargetResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/spy-cats/{catId}/mission/targets/{targetId}/notes [put]
func (h *MissionHandler) UpdateTargetNotes(c echo.Context) error {
	catIDStr := c.Param("catId")
	catID, err := strconv.ParseInt(catIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid cat ID",
			"details": "Cat ID must be a valid integer",
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

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	target, err := h.missionService.UpdateTargetNotes(int32(catID), int32(targetID), req.Notes)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error":   "Failed to update target notes",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, target)
}
