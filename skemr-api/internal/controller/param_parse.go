package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func paramUUID(c *gin.Context, name string) (uuid.UUID, bool) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid " + name})
		return uuid.Nil, false
	}
	return id, true
}
