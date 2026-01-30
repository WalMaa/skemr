package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-common/models"
)

// Context key for user
const CtxUser = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := authenticate(r.Header.Get("Authorization"))
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), CtxUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authenticate(token string) *models.User {
	return &models.User{
		ID:    uuid.New(),
		Email: "example@gmail.com",
	}
}
