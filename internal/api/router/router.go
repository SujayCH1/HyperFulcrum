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
	ShardHandler *handlers.ShardHandler,
	RuntimeStateHandler *handlers.NodeRuntimeStateHandler,
	ShardKeyHandler *handlers.ShardKeyHandler,
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
	shardIDValidator := middleware.UUIDPathValidator("id")
	shardKeyIDValidator := middleware.UUIDPathValidator("projectId", "id")

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
	mux.Handle("PATCH /nodes/{id}/role", nodeIDValidator(http.HandlerFunc(NodeHandler.UpdateNodeRole)))
	mux.Handle("GET /nodes/{id}/runtime-state", nodeIDValidator(http.HandlerFunc(RuntimeStateHandler.GetByNodeID)))
	mux.Handle("GET /projects/{projectId}/runtime-states", projectIDValidator(http.HandlerFunc(RuntimeStateHandler.ListByProject)))

	// Node connection routes
	mux.Handle("POST /nodes/{nodeId}/connection", nodeConnectionIDValidator(nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.AddNodeConnection))))
	mux.Handle("DELETE /nodes/{nodeId}/connection", nodeConnectionIDValidator(http.HandlerFunc(NodeConnectionHandler.RemoveNodeConnection)))
	mux.Handle("PATCH /nodes/{nodeId}/connection", nodeConnectionIDValidator(nodeConnectionValidator(http.HandlerFunc(NodeConnectionHandler.UpdateNodeConnection))))
	mux.Handle("GET /nodes/{nodeId}/connection", nodeConnectionIDValidator(http.HandlerFunc(NodeConnectionHandler.GetNodeConnectionByID)))

	// Shard key routes
	mux.Handle("POST /projects/{projectId}/shard-keys", projectIDValidator(middleware.ShardKeyValidator(http.HandlerFunc(ShardKeyHandler.AddShardKey))))
	mux.Handle("GET /projects/{projectId}/shard-keys", projectIDValidator(http.HandlerFunc(ShardKeyHandler.ListShardKeys)))
	mux.Handle("GET /projects/{projectId}/shard-keys/{tableName}", projectIDValidator(http.HandlerFunc(ShardKeyHandler.GetShardKey)))
	mux.Handle("DELETE /projects/{projectId}/shard-keys/{id}", shardKeyIDValidator(http.HandlerFunc(ShardKeyHandler.DeleteShardKey)))

	// Logical shard routes
	mux.Handle("POST /projects/{projectId}/shards", projectIDValidator(middleware.ShardValidator(http.HandlerFunc(ShardHandler.AddShard))))
	mux.Handle("GET /projects/{projectId}/shards", projectIDValidator(http.HandlerFunc(ShardHandler.ListShards)))
	mux.Handle("GET /shards/{id}", shardIDValidator(http.HandlerFunc(ShardHandler.GetShard)))
	mux.Handle("PATCH /shards/{id}/name", shardIDValidator(http.HandlerFunc(ShardHandler.RenameShard)))
	mux.Handle("DELETE /shards/{id}", shardIDValidator(http.HandlerFunc(ShardHandler.RemoveShard)))

	// Replication & Topology Routes routes
	// Replication routes return 501 until replication is implemented.
	mux.HandleFunc("POST /replication", ReplicationHandler.CreateReplication)
	mux.HandleFunc("DELETE /replication", ReplicationHandler.DeleteReplication)
	mux.HandleFunc("POST /replication/promote", ReplicationHandler.PromoteReplica)
	mux.Handle("POST /projects/{projectId}/topology", projectIDValidator(middleware.TopologyCreateValidator(http.HandlerFunc(NodeTopologyHandler.CreateTopology))))
	mux.Handle("GET /projects/{projectId}/topology", projectIDValidator(http.HandlerFunc(NodeTopologyHandler.ListTopologies)))
	mux.Handle("GET /projects/{projectId}/topology/{relationId}", topologyIDValidator(http.HandlerFunc(NodeTopologyHandler.GetTopologyByID)))
	mux.Handle("DELETE /projects/{projectId}/topology/{relationId}", topologyIDValidator(http.HandlerFunc(NodeTopologyHandler.DeleteTopologyByPath)))

	return mux
}
