package cache

type CacheManager struct {
	Projects    *ProjectStore
	Nodes       *NodeStore
	Connections *ConnectionStore
	Topology    *NodeTopologyStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectStore(),
		Nodes:       NewNodeStore(),
		Connections: NewConnectionStore(),
		Topology:    NewTopologyStore(),
	}
}
