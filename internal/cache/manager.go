package cache

type CacheManager struct {
	Projects    *ProjectCacheStore
	Nodes       *NodeCacheStore
	Connections *ConnectionsCacheStore
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Projects:    NewProjectCacheStore(),
		Nodes:       NewNodeCacheStore(),
		Connections: NewConnectionCacheStore(),
	}
}
