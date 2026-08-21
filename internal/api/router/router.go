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
	projectIDValidator := middleware.UUIDPathValidator("projectId")
	nodeIDValidator := middleware.UUIDPathValidator("id")
	nodeConnectionIDValidator := middleware.UUIDPathValidator("nodeId")
	topologyIDValidator := middleware.UUIDPathValidator("projectId", "relationId")

	// nodeTopologyCreateValidator := middleware.TopologyCreateValidator
	// nodeTopologyDeleteValidator := middleware.TopologyDeleteValidator

	// Project routes
	mux.HandleFunc("GET /projects", ProjectHandler.ListProjects)
	mux.Handle("POST /projects", projectValidator(http.HandlerFunc(ProjectHandler.CreateProject)))
	mux.HandleFunc("GET /projects/ready", ProjectHandler.GetReadyProjects)
	mux.Handle("GET /projects/{id}", nodeIDValidator(http.HandlerFunc(ProjectHandler.GetProjectByID)))
	mux.Handle("DELETE /projects/{id}", nodeIDValidator(http.HandlerFunc(ProjectHandler.RemoveProject)))

	// Node routes
	mux.Handle("POST /projects/{projectId}/nodes", projectIDValidator(nodeValidator(http.HandlerFunc(NodeHandler.AddNode))))
	mux.Handle("GET /projects/{projectId}/nodes", projectIDValidator(http.HandlerFunc(NodeHandler.ListNodes)))
	mux.Handle("DELETE /nodes/{id}", nodeIDValidator(http.HandlerFunc(NodeHandler.RemoveNode)))
	mux.Handle("PUT /nodes/{id}/name", nodeIDValidator(http.HandlerFunc(NodeHandler.UpdateNodeName)))
	mux.Handle("PATCH /nodes/{id}/status", nodeIDValidator(http.HandlerFunc(NodeHandler.UpdateNodeStatus)))
	mux.Handle("PATCH /nodes/{id}/type", nodeIDValidator(http.HandlerFunc(NodeHandler.UpdateNodeType)))

	// Node connection routes
	mux.Handle("POST /nodes/{nodeId}/connection", nodeConnectionIDValidator(nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.AddNodeConnection))))
	mux.Handle("DELETE /nodes/{nodeId}/connection", nodeConnectionIDValidator(http.HandlerFunc(NodeConnectionHandler.RemoveNodeConnection)))
	mux.Handle("PATCH /nodes/{nodeId}/connection", nodeConnectionIDValidator(nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.UpdateNodeConnection))))
	mux.Handle("GET /nodes/{nodeId}/connection", nodeConnectionIDValidator(http.HandlerFunc(NodeConnectionHandler.GetNodeConnectionByID)))

	// Replication & Topology Routes routes
	// Replication routes return 501 until replication is implemented.
	mux.HandleFunc("POST /replication", ReplicationHandler.CreateReplication)
	mux.HandleFunc("DELETE /replication", ReplicationHandler.DeleteReplication)
	mux.HandleFunc("POST /replication/promote", ReplicationHandler.PromoteReplica)
	// mux.Handle("POST /topology/", nodeTopologyCreateValidator(http.HandlerFunc(NodeTopologyHandler.CreateTopology)))
	// mux.Handle("DELETE /topology/", nodeTopologyDeleteValidator(http.HandlerFunc(NodeTopologyHandler.DeleteTopology)))
	mux.Handle("GET /projects/{projectId}/topology", projectIDValidator(http.HandlerFunc(NodeTopologyHandler.ListTopologies)))
	mux.Handle("GET /projects/{projectId}/topology/{relationId}", topologyIDValidator(http.HandlerFunc(NodeTopologyHandler.GetTopologyByID)))

	return mux
}
