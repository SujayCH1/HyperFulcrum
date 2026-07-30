package cache

type CacheManager struct {
	Projects    *ProjectStore
	Nodes       *NodeStore
	Connections *ConnectionStore
	Topology    *NodeTopologyStore

	Columns *ColumnStore
	FKEdges *FKEdgesStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectStore(),
		Nodes:       NewNodeStore(),
		Connections: NewConnectionStore(),
		Topology:    NewTopologyStore(),

		Columns: NewColumnStore(),
		FKEdges: NewFKEdgesStore(),
	}
}
