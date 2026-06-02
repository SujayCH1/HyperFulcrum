package handlers

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/internal/utils"
	"net/http"
)

type ProjectHandler struct {
	repo *repository.ProjectRepository
}

func NewProjectHandler(repo *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo: repo}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	payload, ok := middleware.GetProjectPayload(r)
	if !ok {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve payload from context", nil)
		return
	}

	project, err := h.repo.ProjectAdd(r.Context(), payload.Name, payload.Description)

	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to create project", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusCreated, "Project created successfully", project)
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.ProjectList(r.Context())
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to list projects", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Projects retrieved successfully", projects)
}

func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Project ID is required", nil)
		return
	}

	project, err := h.repo.ProjectGetByID(r.Context(), id)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusNotFound, "Project not found", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Project retrieved successfully", project)
}

func (h *ProjectHandler) RemoveProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.WriteJSONErrorResponse(w, http.StatusBadRequest, "Project ID is required", nil)
		return
	}

	err := h.repo.ProjectRemove(r.Context(), id)
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to remove project", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Project removed successfully", nil)
}

func (h *ProjectHandler) GetReadyProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.repo.ProjectGetReady(r.Context())
	if err != nil {
		utils.WriteJSONErrorResponse(w, http.StatusInternalServerError, "Failed to get ready projects", err)
		return
	}

	utils.WriteJSONSuccessResponse(w, http.StatusOK, "Ready projects retrieved successfully", projects)
}

