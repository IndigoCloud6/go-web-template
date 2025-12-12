package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input")
	if err.Type != ValidationErrorType {
		t.Errorf("expected ValidationErrorType, got %v", err.Type)
	}
	if err.Message != "invalid input" {
		t.Errorf("expected 'invalid input', got %v", err.Message)
	}
	if err.HTTPStatusCode() != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, err.HTTPStatusCode())
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("user not found")
	if err.Type != NotFoundErrorType {
		t.Errorf("expected NotFoundErrorType, got %v", err.Type)
	}
	if err.HTTPStatusCode() != http.StatusNotFound {
		t.Errorf("expected %d, got %d", http.StatusNotFound, err.HTTPStatusCode())
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	err := NewUnauthorizedError("unauthorized")
	if err.Type != UnauthorizedErrorType {
		t.Errorf("expected UnauthorizedErrorType, got %v", err.Type)
	}
	if err.HTTPStatusCode() != http.StatusUnauthorized {
		t.Errorf("expected %d, got %d", http.StatusUnauthorized, err.HTTPStatusCode())
	}
}

func TestNewForbiddenError(t *testing.T) {
	err := NewForbiddenError("forbidden")
	if err.Type != ForbiddenErrorType {
		t.Errorf("expected ForbiddenErrorType, got %v", err.Type)
	}
	if err.HTTPStatusCode() != http.StatusForbidden {
		t.Errorf("expected %d, got %d", http.StatusForbidden, err.HTTPStatusCode())
	}
}

func TestNewConflictError(t *testing.T) {
	err := NewConflictError("conflict")
	if err.Type != ConflictErrorType {
		t.Errorf("expected ConflictErrorType, got %v", err.Type)
	}
	if err.HTTPStatusCode() != http.StatusConflict {
		t.Errorf("expected %d, got %d", http.StatusConflict, err.HTTPStatusCode())
	}
}

func TestNewInternalError(t *testing.T) {
	err := NewInternalError("internal error")
	if err.Type != InternalErrorType {
		t.Errorf("expected InternalErrorType, got %v", err.Type)
	}
	if err.HTTPStatusCode() != http.StatusInternalServerError {
		t.Errorf("expected %d, got %d", http.StatusInternalServerError, err.HTTPStatusCode())
	}
}

func TestErrorWithCause(t *testing.T) {
	cause := errors.New("original error")
	err := NewInternalErrorWithCause("wrapper error", cause)

	if err.Err != cause {
		t.Errorf("expected cause to be set")
	}

	expectedMsg := "wrapper error: original error"
	if err.Error() != expectedMsg {
		t.Errorf("expected '%s', got '%s'", expectedMsg, err.Error())
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("Unwrap should return the original cause")
	}
}

func TestIsErrorFunctions(t *testing.T) {
	validationErr := NewValidationError("test")
	notFoundErr := NewNotFoundError("test")
	unauthorizedErr := NewUnauthorizedError("test")
	forbiddenErr := NewForbiddenError("test")
	conflictErr := NewConflictError("test")
	internalErr := NewInternalError("test")

	if !IsValidationError(validationErr) {
		t.Error("IsValidationError should return true for validation error")
	}
	if !IsNotFoundError(notFoundErr) {
		t.Error("IsNotFoundError should return true for not found error")
	}
	if !IsUnauthorizedError(unauthorizedErr) {
		t.Error("IsUnauthorizedError should return true for unauthorized error")
	}
	if !IsForbiddenError(forbiddenErr) {
		t.Error("IsForbiddenError should return true for forbidden error")
	}
	if !IsConflictError(conflictErr) {
		t.Error("IsConflictError should return true for conflict error")
	}
	if !IsInternalError(internalErr) {
		t.Error("IsInternalError should return true for internal error")
	}

	// Test negative cases
	if IsValidationError(notFoundErr) {
		t.Error("IsValidationError should return false for not found error")
	}
}

func TestGetHTTPStatusCode(t *testing.T) {
	tests := []struct {
		err      error
		expected int
	}{
		{NewValidationError("test"), http.StatusBadRequest},
		{NewNotFoundError("test"), http.StatusNotFound},
		{NewUnauthorizedError("test"), http.StatusUnauthorized},
		{NewForbiddenError("test"), http.StatusForbidden},
		{NewConflictError("test"), http.StatusConflict},
		{NewInternalError("test"), http.StatusInternalServerError},
		{errors.New("generic error"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		got := GetHTTPStatusCode(tt.err)
		if got != tt.expected {
			t.Errorf("GetHTTPStatusCode(%v) = %d, want %d", tt.err, got, tt.expected)
		}
	}
}

func TestGetErrorMessage(t *testing.T) {
	appErr := NewValidationError("custom message")
	genericErr := errors.New("generic error")

	if msg := GetErrorMessage(appErr); msg != "custom message" {
		t.Errorf("GetErrorMessage should return 'custom message', got '%s'", msg)
	}

	if msg := GetErrorMessage(genericErr); msg != "internal server error" {
		t.Errorf("GetErrorMessage should return 'internal server error' for generic error, got '%s'", msg)
	}
}
