package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ctxProjectID = "projectID"

func ProjectIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("projectId")
		// if not a project path, skip
		if param == "" {
			return
		}
		// If param exists it must pass uuid parsing
		id, err := uuid.Parse(param)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid projectId"})
			c.Abort()
			return
		}
		c.Set(ctxProjectID, id)
		c.Next()
	}
}
