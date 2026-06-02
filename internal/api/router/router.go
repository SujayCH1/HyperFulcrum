package router

import (
	"hyperfulcrum/internal/api/handlers"
	"hyperfulcrum/internal/api/middleware"
	"net/http"
)

func NewRouter(projectHandler *handlers.ProjectHandler) *http.ServeMux {
	mux := http.NewServeMux()
	projectValidator := middleware.ProjectValidator

	mux.HandleFunc("GET /projects/", projectHandler.ListProjects)
	mux.Handle("POST /projects", projectValidator(http.HandlerFunc(projectHandler.CreateProject)))
	mux.HandleFunc("GET /projects/ready", projectHandler.GetReadyProjects)
	mux.HandleFunc("GET /projects/{id}", projectHandler.GetProjectByID)
	mux.HandleFunc("DELETE /projects/{id}", projectHandler.RemoveProject)

	return mux
}
