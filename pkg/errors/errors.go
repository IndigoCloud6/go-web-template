package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType int

const (
	// ValidationError represents validation errors (400)
	ValidationErrorType ErrorType = iota
	// NotFoundError represents resource not found errors (404)
	NotFoundErrorType
	// UnauthorizedError represents authentication errors (401)
	UnauthorizedErrorType
	// ForbiddenError represents authorization errors (403)
	ForbiddenErrorType
	// ConflictError represents conflict errors (409)
	ConflictErrorType
	// InternalError represents internal server errors (500)
	InternalErrorType
)

// AppError is a custom error type that provides more context
type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatusCode returns the appropriate HTTP status code for the error type
func (e *AppError) HTTPStatusCode() int {
	switch e.Type {
	case ValidationErrorType:
		return http.StatusBadRequest
	case NotFoundErrorType:
		return http.StatusNotFound
	case UnauthorizedErrorType:
		return http.StatusUnauthorized
	case ForbiddenErrorType:
		return http.StatusForbidden
	case ConflictErrorType:
		return http.StatusConflict
	case InternalErrorType:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *AppError {
	return &AppError{
		Type:    ValidationErrorType,
		Message: message,
	}
}

// NewValidationErrorWithCause creates a new validation error with underlying cause
func NewValidationErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    ValidationErrorType,
		Message: message,
		Err:     err,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(message string) *AppError {
	return &AppError{
		Type:    NotFoundErrorType,
		Message: message,
	}
}

// NewNotFoundErrorWithCause creates a new not found error with underlying cause
func NewNotFoundErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    NotFoundErrorType,
		Message: message,
		Err:     err,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:    UnauthorizedErrorType,
		Message: message,
	}
}

// NewUnauthorizedErrorWithCause creates a new unauthorized error with underlying cause
func NewUnauthorizedErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    UnauthorizedErrorType,
		Message: message,
		Err:     err,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError(message string) *AppError {
	return &AppError{
		Type:    ForbiddenErrorType,
		Message: message,
	}
}

// NewForbiddenErrorWithCause creates a new forbidden error with underlying cause
func NewForbiddenErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    ForbiddenErrorType,
		Message: message,
		Err:     err,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(message string) *AppError {
	return &AppError{
		Type:    ConflictErrorType,
		Message: message,
	}
}

// NewConflictErrorWithCause creates a new conflict error with underlying cause
func NewConflictErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    ConflictErrorType,
		Message: message,
		Err:     err,
	}
}

// NewInternalError creates a new internal server error
func NewInternalError(message string) *AppError {
	return &AppError{
		Type:    InternalErrorType,
		Message: message,
	}
}

// NewInternalErrorWithCause creates a new internal server error with underlying cause
func NewInternalErrorWithCause(message string, err error) *AppError {
	return &AppError{
		Type:    InternalErrorType,
		Message: message,
		Err:     err,
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ValidationErrorType
	}
	return false
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == NotFoundErrorType
	}
	return false
}

// IsUnauthorizedError checks if the error is an unauthorized error
func IsUnauthorizedError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == UnauthorizedErrorType
	}
	return false
}

// IsForbiddenError checks if the error is a forbidden error
func IsForbiddenError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ForbiddenErrorType
	}
	return false
}

// IsConflictError checks if the error is a conflict error
func IsConflictError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ConflictErrorType
	}
	return false
}

// IsInternalError checks if the error is an internal error
func IsInternalError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == InternalErrorType
	}
	return false
}

// GetHTTPStatusCode returns the HTTP status code for an error
// If the error is not an AppError, it returns 500
func GetHTTPStatusCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatusCode()
	}
	return http.StatusInternalServerError
}

// GetErrorMessage returns the error message for display
// It avoids exposing internal error details
func GetErrorMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return "internal server error"
}
