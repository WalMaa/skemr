package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/walmaa/skemr-api/internal/service"
	"github.com/walmaa/skemr-common/models"
)

type DatabaseEntityController struct {
	Service *service.DatabaseEntityService
}

func NewDatabaseEntityController(s *service.DatabaseEntityService) *DatabaseEntityController {
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
	ctx := r.Context()
	projectId, ok := ctx.Value("projectId").(uuid.UUID)
	if !ok {
		http.Error(w, "projectId not found in context", http.StatusBadRequest)
		return
	}
	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid database ID format", http.StatusBadRequest)
		return
	}
	entityId, err := uuid.Parse(chi.URLParam(r, "entityId"))
	if err != nil {
		http.Error(w, "Invalid entity ID format", http.StatusBadRequest)
		return
	}
	entity, err := h.Service.GetDatabaseEntityByID(ctx, projectId, databaseId, entityId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, entity)
}

func (h *DatabaseEntityController) GetDatabaseEntities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectId := ctx.Value("projectId").(uuid.UUID)

	databaseId, err := uuid.Parse(chi.URLParam(r, "databaseId"))
	if err != nil {
		http.Error(w, "Invalid database ID format", http.StatusBadRequest)
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
