package databases

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/walmaa/skemr-api/internal/errormsg"
	"github.com/walmaa/skemr-api/internal/requestparams"
	"github.com/walmaa/skemr-api/internal/validation"
)

type DatabaseController struct {
	Service *DatabaseService
}

func NewDatabaseController(s *DatabaseService) *DatabaseController {
	return &DatabaseController{Service: s}
}

func (h *DatabaseController) RegisterRoutes(r chi.Router) {
	r.Route("/databases", func(r chi.Router) {
		r.Post("/", h.createDatabase)
		r.Get("/", h.listDatabasesByProject)
		r.Get("/{databaseId}", h.getDatabase)
		r.Delete("/{databaseId}", h.deleteDatabase)
		r.Patch("/{databaseId}", h.updateDatabase)
		r.Post("/{databaseId}/sync", h.syncDatabase)
	})
}

func (h *DatabaseController) deleteDatabase(w http.ResponseWriter, r *http.Request) {

	// TODO: scope this to the project
	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}
	err := h.Service.DeleteDatabase(r.Context(), databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *DatabaseController) listDatabasesByProject(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databases, err := h.Service.ListDatabasesByProject(r.Context(), projectId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, databases)
}

func (h *DatabaseController) createDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	var body DatabaseCreationDto

	err := render.Decode(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validation.Validate.Struct(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database, err := h.Service.CreateDatabase(r.Context(), projectId, body)
	if err != nil {
		errormsg.WriteErrorResponse(w, r, err)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, database)
}

func (h *DatabaseController) updateDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	var body DatabaseUpdateDto
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	database, err := h.Service.UpdateDatabase(r.Context(), projectId, databaseId, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, database)
}

func (h *DatabaseController) syncDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	err := h.Service.EnqueueManualDatabaseSync(r.Context(), projectId, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusAccepted)
}

func (h *DatabaseController) getDatabase(w http.ResponseWriter, r *http.Request) {
	projectId, ok := requestparams.ParseUUIDParam(w, r, "projectId")
	if !ok {
		return
	}

	databaseId, ok := requestparams.ParseUUIDParam(w, r, "databaseId")
	if !ok {
		return
	}

	database, err := h.Service.GetDatabase(r.Context(), projectId, databaseId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, database)
}
