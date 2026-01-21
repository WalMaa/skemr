package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"

	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := authenticate(c.Request.Header.Get("Authorization"))
		if user == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func authenticate(token string) *models.User {
	return &models.User{
		ID:    uuid.New(),
		Email: "example@gmail.com",
	}
}
