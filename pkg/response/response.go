package response

import (
	"net/http"

	apperrors "github.com/IndigoCloud6/go-web-template/pkg/errors"
	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success sends a success response
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error sends an error response with appropriate HTTP status code
func Error(c *gin.Context, code int, message string) {
	httpStatus := http.StatusInternalServerError
	if code >= 400 && code < 600 {
		httpStatus = code
	}
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorFromAppError sends an error response based on AppError type
// This function maps custom error types to appropriate HTTP status codes
func ErrorFromAppError(c *gin.Context, err error) {
	httpStatus := apperrors.GetHTTPStatusCode(err)
	message := apperrors.GetErrorMessage(err)
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
	})
}

// SuccessWithMessage sends a success response with custom message
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// BadRequest sends a bad request error response
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// InternalServerError sends an internal server error response
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}
