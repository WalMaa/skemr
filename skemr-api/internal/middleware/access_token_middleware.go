package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/controller"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

func AccessTokenMiddleware(service *service.AccessTokenService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := r.Context()

			projectId, ok := controller.ParseUUIDParam(w, r, "projectId")
			if !ok {
				return
			}

			tokenHeaderValue := r.Header.Get("Authorization")
			if tokenHeaderValue == "" {
				errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
					Message: "Authorization header is required",
					Errors:  nil,
					Status:  http.StatusBadRequest,
				})
				return
			}
			token := tokenHeaderValue[len("Bearer "):]

			ok, err := authenticateToken(c, service, projectId, strings.TrimSpace(token))

			if err != nil {
				errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
					Message: "Error validating token",
					Status:  http.StatusInternalServerError,
				})
				return
			}

			if !ok {
				errormsg.WriteErrorResponse(w, r, &models.ErrorResponse{
					Message: "Invalid token",
					Status:  http.StatusUnauthorized,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func authenticateToken(c context.Context, service *service.AccessTokenService, projectId uuid.UUID, token string) (bool, error) {
	// Check if the token is valid using the service
	ok, err := service.ValidateToken(c, projectId, token)

	if err != nil {
		return false, err
	}
	return ok, nil
}
