package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return cv.FormatValidationErrors(err)
	}
	return nil
}

func (cv *CustomValidator) FormatValidationErrors(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var messages []string

		for _, fe := range ve {
			message := cv.getErrorMessage(fe)
			messages = append(messages, message)
		}

		return fmt.Errorf("validation failed: %s", strings.Join(messages, "; "))
	}

	return fmt.Errorf("validation error: %w", err)
}

func (cv *CustomValidator) getErrorMessage(fe validator.FieldError) string {
	field := cv.getFieldDisplayName(fe.Field())

	switch fe.Tag() {
	case "required":
		if fe.Type().Kind() >= reflect.Int && fe.Type().Kind() <= reflect.Uint64 && fe.Value() == reflect.Zero(fe.Type()).Interface() {
			return fmt.Sprintf("%s must be at least 1", field)
		}
		return fmt.Sprintf("%s is required", field)
	case "min":
		if fe.Type().Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters long", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		if fe.Type().Kind() == reflect.String {
			return fmt.Sprintf("%s must not exceed %s characters", field, fe.Param())
		}
		return fmt.Sprintf("%s must not exceed %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, fe.Param())
	case "dive":
		return fmt.Sprintf("invalid %s item", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

func (cv *CustomValidator) getFieldDisplayName(field string) string {
	switch field {
	case "name":
		return "Name"
	case "breed":
		return "Breed"
	case "years_of_experience":
		return "Years of experience"
	case "salary":
		return "Salary"
	case "description":
		return "Description"
	case "start_date":
		return "Start date"
	case "end_date":
		return "End date"
	case "targets":
		return "Targets"
	case "country":
		return "Country"
	case "status":
		return "Status"
	case "notes":
		return "Notes"
	case "cat_id":
		return "Cat ID"
	default:
		parts := strings.Split(field, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			}
		}
		return strings.Join(parts, " ")
	}
}
