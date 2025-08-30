package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"spy-cat-agency/internal/application/dto"
	"spy-cat-agency/internal/application/services"
	"spy-cat-agency/internal/domain/entities"
	"spy-cat-agency/internal/domain/interfaces"
	"spy-cat-agency/internal/infrastructure/external"

	"github.com/labstack/echo/v4"
)

type CatHandler struct {
	catRepo           interfaces.CatRepository
	breedService      *external.BreedService
	validationService *services.ValidationService
}

func NewCatHandler(catRepo interfaces.CatRepository) *CatHandler {
	return &CatHandler{
		catRepo:           catRepo,
		breedService:      external.NewBreedService(),
		validationService: services.NewValidationService(),
	}
}

// CreateCat creates a new spy cat
// @Summary Create a new spy cat
// @Description Create a new spy cat with the provided information
// @Tags cats
// @Accept json
// @Produce json
// @Param cat body dto.CreateCatRequest true "Cat creation request"
// @Success 201 {object} dto.CatResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cats [post]
func (h *CatHandler) CreateCat(c echo.Context) error {
	var req dto.CreateCatRequest

	if err := h.validationService.ValidateJSONBinding(c, &req, "cat"); err != nil {
		return c.JSON(http.StatusBadRequest, h.validationService.CreateErrorResponse(
			"Invalid request body format",
			"Please check your JSON format and field types",
		))
	}

	if err := h.validationService.ValidateEchoStruct(c, &req, "cat"); err != nil {
		return c.JSON(http.StatusBadRequest, h.validationService.CreateErrorResponse(
			"Validation failed",
			"Name, breed, experience and salary are required",
		))
	}

	isValidBreed, err := h.breedService.ValidateBreed(req.Breed)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.validationService.CreateErrorResponse(
			"Failed to validate breed",
			"Unable to connect to breed validation service",
		))
	}
	if !isValidBreed {
		return c.JSON(http.StatusBadRequest, h.validationService.CreateErrorResponse(
			"Invalid breed",
			"Please use a valid cat breed from TheCatAPI. Available breeds can be fetched from /api/v1/breeds endpoint",
		))
	}

	spyCat := entities.NewSpyCat(req.Name, req.Breed, req.YearsOfExperience, req.Salary)

	created, err := h.catRepo.Create(c.Request().Context(), spyCat)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create cat",
			"details": "Database error occurred while creating the spy cat",
		})
	}

	response := h.toResponseDTO(created)
	return c.JSON(http.StatusCreated, response)
}

// GetCat retrieves a cat by ID
// @Summary Get a spy cat by ID
// @Description Get a spy cat by its ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} dto.CatResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/cats/{id} [get]
func (h *CatHandler) GetCat(c echo.Context) error {
	id, err := h.validationService.ValidateCatID(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.validationService.CreateErrorResponse("Invalid cat ID", err.Error()))
	}

	spyCat, err := h.catRepo.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, h.validationService.CreateErrorResponse("Cat not found", ""))
	}

	response := h.toResponseDTO(spyCat)
	return c.JSON(http.StatusOK, response)
}

// ListCats retrieves cats with pagination
// @Summary List spy cats
// @Description Get a paginated list of spy cats
// @Tags cats
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.CatListResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/cats [get]
func (h *CatHandler) ListCats(c echo.Context) error {
	limit, offset, err := h.validationService.ValidatePaginationParams(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.validationService.CreateErrorResponse("Invalid pagination parameters", err.Error()))
	}

	spyCats, err := h.catRepo.List(c.Request().Context(), limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.validationService.CreateErrorResponse("Failed to list cats", ""))
	}

	breeds, err := h.breedService.GetBreedNames()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.validationService.CreateErrorResponse("Failed to fetch breeds", ""))
	}

	catResponses := make([]dto.CatResponse, len(spyCats))
	for i, spyCat := range spyCats {
		catResponses[i] = *h.toResponseDTO(spyCat)
	}

	response := dto.CatListResponse{
		Cats:   catResponses,
		Breeds: breeds,
		Total:  int64(len(catResponses)), // simplified - in real app would be actual count
		Limit:  limit,
		Offset: offset,
	}

	return c.JSON(http.StatusOK, response)
}

// UpdateCatSalary updates a cat's salary
// @Summary Update spy cat salary
// @Description Update the salary of a spy cat
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Param salary body dto.UpdateCatSalaryRequest true "Salary update request"
// @Success 200 {object} dto.CatResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/cats/{id}/salary [put]
func (h *CatHandler) UpdateCatSalary(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid cat ID",
			"details": "Cat ID must be a positive integer",
		})
	}

	var req dto.UpdateCatSalaryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Invalid request body format",
			"details": "Please check your JSON format and field types",
		})
	}

	// Validate request body
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	updated, err := h.catRepo.UpdateSalary(c.Request().Context(), int32(id), req.Salary)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error":   "Cat not found",
			"details": "No spy cat exists with the provided ID",
		})
	}

	response := h.toResponseDTO(updated)
	return c.JSON(http.StatusOK, response)
}

// DeleteCat deletes a cat by ID
// @Summary Delete a spy cat
// @Description Delete a spy cat by its ID
// @Tags cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/cats/{id} [delete]
func (h *CatHandler) DeleteCat(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid cat ID"})
	}

	err = h.catRepo.Delete(c.Request().Context(), int32(id))
	if err != nil {
		if strings.Contains(err.Error(), "cannot delete cat: cat is currently assigned to mission") {
			return c.JSON(http.StatusConflict, map[string]string{
				"error":   "Cannot delete cat",
				"details": "This spy cat is currently assigned to a mission. Please unassign the cat from the mission before deletion.",
			})
		}
		if strings.Contains(err.Error(), "failed to find cat") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Cat not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete cat"})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// GetBreeds returns the list of valid cat breeds from TheCatAPI
// @Summary Get valid cat breeds
// @Description Get a list of valid cat breeds for cat creation
// @Tags cats
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/cats/breeds [get]
func (h *CatHandler) GetBreeds(c echo.Context) error {
	breeds, err := h.breedService.GetBreedNames()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch cat breeds"})
	}

	return c.JSON(http.StatusOK, map[string][]string{"breeds": breeds})
}

func (h *CatHandler) toResponseDTO(spyCat *entities.SpyCat) *dto.CatResponse {
	return &dto.CatResponse{
		ID:                spyCat.ID,
		Name:              spyCat.Name,
		YearsOfExperience: spyCat.YearsOfExperience,
		Breed:             spyCat.Breed,
		Salary:            spyCat.Salary,
		MissionID:         spyCat.MissionID,
		CreatedAt:         spyCat.CreatedAt,
		UpdatedAt:         spyCat.UpdatedAt,
	}
}
