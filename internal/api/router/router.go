package router

import (
	"hyperfulcrum/internal/api/middleware"
	"hyperfulcrum/internal/api/services"
	"net/http"
)

func NewRouter(projectService *services.ProjectService,nodeService *services.NodeService) *http.ServeMux {

	mux := http.NewServeMux()

	projectValidator := middleware.ProjectValidator
	nodeValidator := middleware.NodeValidator

	//Project routes
	mux.HandleFunc("GET /projects/", projectService.ListProjects)
	mux.Handle("POST /projects", projectValidator(http.HandlerFunc(projectService.CreateProject)))
	mux.HandleFunc("GET /projects/ready", projectService.GetReadyProjects)
	mux.HandleFunc("GET /projects/{id}", projectService.GetProjectByID)
	mux.HandleFunc("DELETE /projects/{id}", projectService.RemoveProject)

	// Node routes
	mux.Handle("POST /projects/{projectId}/nodes", nodeValidator(http.HandlerFunc(nodeService.AddNode)))
	mux.HandleFunc("GET /projects/{projectId}/nodes", nodeService.ListNodes)
	mux.HandleFunc("DELETE /nodes/{id}", nodeService.RemoveNode)
	mux.HandleFunc("PUT /nodes/{id}/name", nodeService.UpdateNodeName)
	mux.HandleFunc("PATCH /nodes/{id}/status", nodeService.UpdateNodeStatus)
	mux.HandleFunc("PATCH /nodes/{id}/type", nodeService.UpdateNodeType)

	//http://localhost:8080/nodes/c5f8bcd2-a2b6-4416-9001-699ea1621538/name?name=nodee-1
	return mux
}
