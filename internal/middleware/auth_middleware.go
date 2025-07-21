package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/walmaa/skemr/db/sqlc"
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

func authenticate(token string) *sqlc.User {
	return &sqlc.User{
		ID:       uuid.New(),
		Email:    "user@mail.com",
		Password: "password",
	}
}
