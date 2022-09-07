package helpers

import (
	"database/sql"
	"net/http"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/translator"
	customValidator "github.com/abdivasiyev/project_template/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

func ToCustomError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrNotFound
	}

	return err
}

func ConvertErrorToErrorResponse(statusCode int, err error) models.ErrorResponse {
	if err == nil {
		return models.ErrorResponse{
			ErrorCode:    http.StatusOK,
			ErrorMessage: "",
		}
	}

	result, ok := convertCustomErrors(err)
	if ok {
		return result
	}

	result, ok = convertValidationErrors(err)
	if ok {
		return result
	}

	result.ErrorCode = statusCode
	return result
}

func convertCustomErrors(err error) (models.ErrorResponse, bool) {
	var (
		customValidationError customValidator.ValidationError
	)

	if !errors.As(err, &customValidationError) {
		return models.ErrorResponse{}, false
	}

	result := models.ErrorResponse{
		ErrorCode:    http.StatusBadRequest,
		ErrorMessage: http.StatusText(http.StatusBadRequest),
	}

	result.Validations = []models.ValidationResponse{
		{
			Field: customValidationError.Field(),
			Error: customValidationError.Message(),
		},
	}

	return result, true
}

func convertValidationErrors(err error) (models.ErrorResponse, bool) {
	var (
		validationErrors validator.ValidationErrors
	)

	if !errors.As(err, &validationErrors) {
		return models.ErrorResponse{}, false
	}

	result := models.ErrorResponse{
		ErrorCode:    http.StatusBadRequest,
		ErrorMessage: http.StatusText(http.StatusBadRequest),
	}

	result.Validations = make([]models.ValidationResponse, len(validationErrors))

	for i, validationErr := range validationErrors {
		translatedErr := validationErr.Translate(translator.Get())

		result.Validations[i] = models.ValidationResponse{
			Field: validationErr.Field(),
			Error: translatedErr,
		}
	}

	return result, true
}
