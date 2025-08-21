package response

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Title  string       `json:"title,omitempty"`
	Detail string       `json:"detail,omitempty"`
	Errors []FieldError `json:"errors,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ParseValidationErrors(err error) []FieldError {
	var verrs validator.ValidationErrors
	if errors.As(err, &verrs) {
		fieldErrors := make([]FieldError, 0, len(verrs))
		for _, fe := range verrs {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   fe.Field(),
				Code:    fe.Tag(),
				Message: fe.Error(),
			})
		}
		return fieldErrors
	}
	return nil
}

func ValidationError(errors []FieldError) *APIResponse {
	return Err(http.StatusBadRequest, &APIError{
		Title:  "Validation failed",
		Errors: errors,
	})
}

func BadRequest(detail string) *APIResponse {
	return Err(http.StatusBadRequest, &APIError{
		Title:  "Bad Request",
		Detail: detail,
	})
}

func NotFound(resource string) *APIResponse {
	return Err(http.StatusNotFound, &APIError{
		Title:  "Not Found",
		Detail: resource + " not found",
	})
}

func InternalServerError(detail string) *APIResponse {
	return Err(http.StatusInternalServerError, &APIError{
		Title:  "Internal Server Error",
		Detail: detail,
	})
}

func ServiceUnavailable(detail string) *APIResponse {
	return Err(http.StatusServiceUnavailable, &APIError{
		Title:  "Service Unavailable",
		Detail: detail,
	})
}
