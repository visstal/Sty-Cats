package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type ValidationService struct{}

func NewValidationService() *ValidationService {
	return &ValidationService{}
}

func (v *ValidationService) ValidateID(idStr, paramName string) (int32, error) {
	if idStr == "" {
		return 0, fmt.Errorf("missing %s parameter", paramName)
	}

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: must be a valid integer", paramName)
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid %s: must be a positive integer", paramName)
	}

	return int32(id), nil
}

func (v *ValidationService) ValidateCatID(c echo.Context) (int32, error) {
	return v.ValidateID(c.Param("id"), "cat ID")
}

func (v *ValidationService) ValidateMissionID(c echo.Context) (int32, error) {
	return v.ValidateID(c.Param("id"), "mission ID")
}

func (v *ValidationService) ValidateMissionIDFromParam(c echo.Context, paramName string) (int32, error) {
	return v.ValidateID(c.Param(paramName), "mission ID")
}

func (v *ValidationService) ValidateTargetID(c echo.Context) (int32, error) {
	return v.ValidateID(c.Param("targetId"), "target ID")
}

func (v *ValidationService) ValidatePaginationParams(c echo.Context) (limit, offset int32, err error) {
	limit = 10
	offset = 0

	limitStr := c.QueryParam("limit")
	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter: must be a valid integer")
		} else {
			if l <= 0 {
				return 0, 0, fmt.Errorf("invalid limit parameter: must be a positive integer")
			}
			if l > 100 {
				return 0, 0, fmt.Errorf("invalid limit parameter: maximum allowed is 100")
			}
			limit = int32(l)
		}
	}

	offsetStr := c.QueryParam("offset")
	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter: must be a valid integer")
		} else {
			if o < 0 {
				return 0, 0, fmt.Errorf("invalid offset parameter: cannot be negative")
			}
			offset = int32(o)
		}
	}

	return limit, offset, nil
}

func (v *ValidationService) ValidateJSONBinding(c echo.Context, req interface{}, entityName string) error {
	if err := c.Bind(req); err != nil {
		return fmt.Errorf("invalid request body format for %s: please check your JSON format and field types", entityName)
	}
	return nil
}

func (v *ValidationService) ValidateEchoStruct(c echo.Context, req interface{}, entityName string) error {
	if err := c.Validate(req); err != nil {
		return fmt.Errorf("validation failed for %s: %s", entityName, err.Error())
	}
	return nil
}

func (v *ValidationService) ValidateTargetStatus(status string) error {
	validStatuses := []string{"init", "in_progress", "completed"}
	status = strings.ToLower(strings.TrimSpace(status))

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return nil
		}
	}

	return fmt.Errorf("invalid target status: must be one of %v", validStatuses)
}

func (v *ValidationService) CreateErrorResponse(message, details string) map[string]interface{} {
	response := map[string]interface{}{
		"error": message,
	}

	if details != "" {
		response["details"] = details
	}

	return response
}

func validateMissionName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("mission name is required")
	}
	if len(name) > 100 {
		return fmt.Errorf("mission name must be less than 100 characters")
	}
	return nil
}

func validateMissionDescription(description string) error {
	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Errorf("mission description is required")
	}
	if len(description) > 500 {
		return fmt.Errorf("mission description must be less than 500 characters")
	}
	return nil
}

func validateMissionDates(startDate, endDate time.Time) error {
	if startDate.IsZero() || endDate.IsZero() {
		return nil
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("mission end date must be after start date")
	}
	return nil
}
