package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request first

		// Check if errors were added to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			c.JSON(http.StatusInternalServerError, APIError{Path: c.FullPath(), Status: http.StatusInternalServerError, Message: err.Error()})

		}

	}
}
