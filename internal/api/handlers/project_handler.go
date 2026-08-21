package handlers

import (
	"net/http"

	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/pkg/utils"
)

type ProjectHandler struct {
	service *metadata.ProjectService
}

func NewProjectHandler(
	service *metadata.ProjectService,
) *ProjectHandler {

	return &ProjectHandler{
		service: service,
	}
}

func (h *ProjectHandler) CreateProject(
	w http.ResponseWriter,
	r *http.Request,
) {

	payload, ok := middleware.GetProjectPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(
			w,
			http.StatusInternalServerError,
			"Failed to retrieve payload",
			nil,
		)
		return
	}

	project, err := h.service.CreateProject(
		r.Context(),
		payload.Name,
		payload.Description,
	)
	if err != nil {
		writeHandlerError(
			w,
			"Project not found",
			"Failed to create project",
			err,
		)
		return
	}

	utils.WriteJSONSuccessResponse(
		w,
		http.StatusCreated,
		"Project created successfully",
		project,
	)
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to list projects", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Projects retrieved successfully", projects)
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	project, err := h.service.GetProjectByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to get project", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Project retrieved successfully", project)
}

func (h *ProjectHandler) RemoveProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.DeleteProject(r.Context(), id)
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to remove project", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) GetReadyProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetReadyProjects(r.Context())
	if err != nil {
		writeHandlerError(w, "Project not found", "Failed to get ready projects", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Ready projects retrieved successfully", projects)
}
