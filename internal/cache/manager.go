package cache

type CacheManager struct {
	Projects     *ProjectStore
	Nodes        *NodeStore
	ProjectNodes *ProjectNodeStore
	Connections  *ConnectionStore
	Topology     *TopologyStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:     NewProjectStore(),
		Nodes:        NewNodeStore(),
		ProjectNodes: NewProjectNodeStore(),
		Connections:  NewConnectionStore(),
		Topology:     NewTopologyStore(),
	}
}
