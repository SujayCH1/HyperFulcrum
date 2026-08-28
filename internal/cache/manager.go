package cache

type CacheManager struct {
	Projects    *ProjectStore
	Nodes       *NodeStore
	Connections *ConnectionStore
	Topology    *NodeTopologyStore
	Shards      *ShardStore
	Runtime     *NodeRuntimeStateStore

	Columns       *ColumnStore
	FKEdges       *FKEdgesStore
	SchemaVersion *SchemaVersionStore
	ShardKeys     *ShardKeysStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectStore(),
		Nodes:       NewNodeStore(),
		Connections: NewConnectionStore(),
		Topology:    NewTopologyStore(),
		Shards:      NewShardStore(),
		Runtime:     NewNodeRuntimeStateStore(),

		Columns:       NewColumnStore(),
		FKEdges:       NewFKEdgesStore(),
		SchemaVersion: NewSchemaVersionStore(),
		ShardKeys:     NewShardKeysStore(),
	}
}

func (m *CacheManager) DeleteProject(projectID string) {
	m.Projects.Delete(projectID)
	m.Nodes.DeleteProject(projectID)
	m.Connections.DeleteProject(projectID)
	m.Topology.DeleteByProjectID(projectID)
	m.Shards.DeleteProject(projectID)
	m.Runtime.DeleteProject(projectID)
	m.Columns.DeleteProject(projectID)
	m.FKEdges.DeleteProject(projectID)
	m.SchemaVersion.DeleteProject(projectID)
	m.ShardKeys.DeleteProject(projectID)
}

func (m *CacheManager) Clear() {
	m.Projects.Clear()
	m.Nodes.Clear()
	m.Connections.Clear()
	m.Topology.Clear()
	m.Shards.Clear()
	m.Runtime.Clear()
	m.Columns.Clear()
	m.FKEdges.Clear()
	m.SchemaVersion.Clear()
	m.ShardKeys.Clear()
}
