package middleware

import (
	"net/http"

	"github.com/IndigoCloud6/go-web-template/pkg/logger"
	"github.com/IndigoCloud6/go-web-template/pkg/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery middleware recovers from panics and returns an error response
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				c.JSON(http.StatusInternalServerError, response.Response{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
