package router

import (
	"hyperfulcrum/internal/api/handlers"
	"hyperfulcrum/internal/api/middleware"

	"net/http"
)

func NewRouter(
	ProjectHandler *handlers.ProjectHandler,
	NodeHandler *handlers.NodeHandler,
	NodeConnectionHandler *handlers.NodeConnectionHandler,
	NodeTopologyHandler *handlers.TopologyHandler,
	ReplicationHandler *handlers.ReplicationHandler,
) *http.ServeMux {

	mux := http.NewServeMux()

	projectValidator := middleware.ProjectValidator
	nodeValidator := middleware.NodeValidator
	nodeConnectionValidator := middleware.NodeConnectionValidator

	// nodeTopologyCreateValidator := middleware.TopologyCreateValidator
	// nodeTopologyDeleteValidator := middleware.TopologyDeleteValidator

	replicationCreateValidator := middleware.ReplicationCreateValidator
	replicationDeleteValidator := middleware.ReplicationDeleteValidator
	replicationPromoteValidator := middleware.ReplicationPromoteValidator

	// Project routes
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

	// Node connection routes
	mux.Handle("POST /nodes/{nodeId}/connection", nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.AddNodeConnection)))
	mux.HandleFunc("DELETE /nodes/{nodeId}/connection", NodeConnectionHandler.RemoveNodeConnection)
	mux.HandleFunc("PATCH /nodes/{nodeId}/connection", NodeConnectionHandler.UpdateNodeConnection)
	mux.HandleFunc("GET /nodes/{nodeId}/connection", NodeConnectionHandler.GetNodeConnectionByID)

	// Replication & Topology Routes routes
	mux.Handle("POST /replication/", replicationCreateValidator(http.HandlerFunc(ReplicationHandler.CreateReplication)))
	mux.Handle("DELETE /replication/", replicationDeleteValidator(http.HandlerFunc(ReplicationHandler.DeleteReplication)))
	mux.Handle("POST /replication/promote", replicationPromoteValidator(http.HandlerFunc(ReplicationHandler.PromoteReplica)))
	// mux.Handle("POST /topology/", nodeTopologyCreateValidator(http.HandlerFunc(NodeTopologyHandler.CreateTopology)))
	// mux.Handle("DELETE /topology/", nodeTopologyDeleteValidator(http.HandlerFunc(NodeTopologyHandler.DeleteTopology)))
	mux.HandleFunc("GET /projects/{projectId}/topology", NodeTopologyHandler.ListTopologies)
	mux.HandleFunc("GET /projects/{projectId}/topology/{relationId}", NodeTopologyHandler.GetTopologyByID)

	return mux
}
