package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const CtxProjectID = "projectId"

func ProjectIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "projectId")
		if param == "" {
			next.ServeHTTP(w, r)
			return
		}
		id, err := uuid.Parse(param)
		if err != nil {
			http.Error(w, "invalid projectId", http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), CtxProjectID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
