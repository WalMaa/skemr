package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-common/models"
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
			errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid project id",
			})
			return
		}
		ctx := context.WithValue(r.Context(), CtxProjectID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
