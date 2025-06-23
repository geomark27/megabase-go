package utils

import (
	"fmt"
	"net/http"
)

// APIError representa un error con código HTTP específico
type APIError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
	Details    string `json:"details,omitempty"`
}

// Error implementa la interfaz error
func (e *APIError) Error() string {
	return e.Message
}

// GetStatusCode retorna el código HTTP del error
func (e *APIError) GetStatusCode() int {
	return e.StatusCode
}

// Funciones de creación de errores específicos

// NewBadRequestError crea un error 400
func NewBadRequestError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewUnauthorizedError crea un error 401
func NewUnauthorizedError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError crea un error 403
func NewForbiddenError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewNotFoundError crea un error 404
func NewNotFoundError(resource string) *APIError {
	return &APIError{
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewConflictError crea un error 409
func NewConflictError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewInternalServerError crea un error 500
func NewInternalServerError(message string) *APIError {
	return &APIError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
	}
}

// NewValidationError crea un error de validación con detalles
func NewValidationError(details string) *APIError {
	return &APIError{
		Message:    "Invalid request data",
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

// Helper para verificar si un error es APIError
func IsAPIError(err error) (*APIError, bool) {
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}