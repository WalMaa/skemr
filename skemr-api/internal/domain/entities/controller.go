package entities

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/requestparams"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityController struct {
	Service *DatabaseEntityService
}

func NewDatabaseEntityController(s *DatabaseEntityService) *DatabaseEntityController {
	return &DatabaseEntityController{Service: s}
}

func (h *DatabaseEntityController) RegisterRoutes(r chi.Router) {
	r.Route("/databases/{databaseId}/entities", func(r chi.Router) {
		r.Get("/", h.GetDatabaseEntities)
		r.Get("/{entityId}", h.GetDatabaseEntity)
	})
}

type Query struct {
	EntityType models.DatabaseEntityType
}

func (h *DatabaseEntityController) GetDatabaseEntity(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	entityId, ok := requestparams.ParseUUIDParam(w, r, "entityId")
	if !ok {
		return
	}

	entity, err := h.Service.GetDatabaseEntityByID(r.Context(), projectId, databaseId, entityId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, entity)
}

func (h *DatabaseEntityController) GetDatabaseEntities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	entityTypeQuery := r.URL.Query().Get("type")
	var entityType *models.DatabaseEntityType
	if entityTypeQuery != "" {
		et := models.DatabaseEntityType(entityTypeQuery)
		entityType = &et
	}
	parentIdQuery := r.URL.Query().Get("parentId")
	var parentId *uuid.UUID
	if parentIdQuery != "" {
		pId, err := uuid.Parse(parentIdQuery)
		if err != nil {
			http.Error(w, "Invalid parent ID format", http.StatusBadRequest)
			return
		}
		parentId = &pId
	}
	entities, err := h.Service.ListDatabaseEntitiesByDatabase(ctx, projectId, databaseId, entityType, parentId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, entities)
}
