package cache

import (
	"hyperfulcrum/internal/repository"
)

type Fetcher struct {
	projectRepo   repository.ProjectRepository
	nodesRepo     repository.NodeRepository
	nodesConnRepo repository.NodeRepository
	cache         *CacheManager
}

func NewFetcher(
	projects repository.ProjectRepository,
	nodes repository.NodeRepository,
	connections repository.NodeRepository,
	cache *CacheManager,
) *Fetcher {
	return &Fetcher{
		projectRepo:   projects,
		nodesRepo:     nodes,
		nodesConnRepo: connections,
		cache:         cache,
	}
}
