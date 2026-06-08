package router

import (
	"hyperfulcrum/internal/api/handlers"
	"hyperfulcrum/internal/api/middleware"

	"net/http"
)

func NewRouter(ProjectHandler *handlers.ProjectHandler, NodeHandler *handlers.NodeHandler, NodeConnectionHandler *handlers.NodeConnectionHandler) *http.ServeMux {

	mux := http.NewServeMux()

	projectValidator := middleware.ProjectValidator
	nodeValidator := middleware.NodeValidator
	nodeConnectionValidator := middleware.NodeConnectionValidator

	//Project routes
	mux.HandleFunc("GET /projects/", ProjectHandler.ListProjects)
	mux.Handle("POST /projects", projectValidator(http.HandlerFunc(ProjectHandler.CreateProject)))
	mux.HandleFunc("GET /projects/ready", ProjectHandler.GetReadyProjects)
	mux.HandleFunc("GET /projects/{id}", ProjectHandler.GetProjectByID)
	mux.HandleFunc("DELETE /projects/{id}", ProjectHandler.RemoveProject)

	// Node routes
	mux.Handle("POST /projects/{projectId}/nodes", nodeValidator(http.HandlerFunc(NodeHandler.AddNode)))
	mux.HandleFunc("GET /projects/{projectId}/nodes", NodeHandler.ListNodes)
	mux.HandleFunc("DELETE /nodes/{id}", NodeHandler.RemoveNode)
	mux.HandleFunc("PUT /nodes/{id}/name", NodeHandler.UpdateNodeName)
	mux.HandleFunc("PATCH /nodes/{id}/status", NodeHandler.UpdateNodeStatus)
	mux.HandleFunc("PATCH /nodes/{id}/type", NodeHandler.UpdateNodeType)

	// node connection routes
	mux.Handle("POST /nodes/{nodeId}/connection", nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.AddNodeConnection)))
	mux.HandleFunc("DELETE /nodes/{nodeId}/connection", NodeConnectionHandler.RemoveNodeConnection)
	mux.HandleFunc("PATCH /nodes/{nodeId}/connection", NodeConnectionHandler.UpdateNodeConnection)
	mux.HandleFunc("GET /nodes/{nodeId}/connectio", NodeConnectionHandler.GetNodeConnectionByID)

	//http://localhost:8080/nodes/c5f8bcd2-a2b6-4416-9001-699ea1621538/name?name=nodee-1
	return mux
}
