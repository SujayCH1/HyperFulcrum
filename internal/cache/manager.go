package cache

type CacheManager struct {
	Projects    *ProjectStore
	Nodes       *NodeStore
	Connections *ConnectionStore
	Topology    *NodeTopologyStore

	Columns       *ColumnStore
	FKEdges       *FKEdgesStore
	SchemaVersion *SchemaVersionStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectStore(),
		Nodes:       NewNodeStore(),
		Connections: NewConnectionStore(),
		Topology:    NewTopologyStore(),

		Columns:       NewColumnStore(),
		FKEdges:       NewFKEdgesStore(),
		SchemaVersion: NewSchemaVersionStore(),
	}
}
