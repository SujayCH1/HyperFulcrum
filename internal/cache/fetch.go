package cache

import (
	"context"
	"fmt"

	"hyperfulcrum/internal/repository"
)

type Fetcher struct {
	projectRepo   repository.ProjectRepository
	nodesRepo     repository.NodeRepository
	nodesConnRepo repository.NodeConnectionRepository
	cache         *CacheManager
}

func NewFetcher(
	projects repository.ProjectRepository,
	nodes repository.NodeRepository,
	connections repository.NodeConnectionRepository,
	cache *CacheManager,
) *Fetcher {
	return &Fetcher{
		projectRepo:   projects,
		nodesRepo:     nodes,
		nodesConnRepo: connections,
		cache:         cache,
	}
}

func (f *Fetcher) RefreshProject(ctx context.Context, projectID string) error {
	project, err := f.projectRepo.ProjectGetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("fetch project %s: %w", projectID, err)
	}
	f.cache.Projects.Set(projectID, *project)

	nodes, err := f.nodesRepo.NodeGetByPorjectID(ctx, projectID) //GetNodesByProjectID
	if err != nil {
		return fmt.Errorf("fetch nodes for project %s: %w", projectID, err)
	}
	f.cache.Nodes.Set(projectID, nodes)

	for _, node := range nodes {
		connections, err := f.nodesConnRepo.GetConnectionByNodeId(ctx, node.ID)
		if err != nil {
			return fmt.Errorf("fetch connections for node %s: %w", node.ID, err)
		}
		f.cache.Connections.Set(node.ID, connections)
	}

	return nil
}

func (f *Fetcher) GetProject(ctx context.Context, projectID string) (repository.Project, error) {
	if project, ok := f.cache.Projects.Get(projectID); ok {
		return project, nil
	}

	if err := f.RefreshProject(ctx, projectID); err != nil {
		return repository.Project{}, err
	}

	project, ok := f.cache.Projects.Get(projectID)
	if !ok {
		return repository.Project{}, fmt.Errorf("project %s not found after refresh", projectID)
	}

	return project, nil
}

func (f *Fetcher) GetNodes(ctx context.Context, projectID string) ([]repository.Node, error) {
	if nodes, ok := f.cache.Nodes.Get(projectID); ok {
		return nodes, nil
	}

	if err := f.RefreshProject(ctx, projectID); err != nil {
		return nil, err
	}

	nodes, ok := f.cache.Nodes.Get(projectID)
	if !ok {
		return nil, fmt.Errorf("nodes for project %s not found after refresh", projectID)
	}

	return nodes, nil
}

func (f *Fetcher) GetConnections(ctx context.Context, nodeID string) ([]repository.NodeConnection, error) {
	if connections, ok := f.cache.Connections.Get(nodeID); ok {
		return connections, nil
	}

	// load directly for this node if missing
	connections, err := f.nodesConnRepo.GetConnectionByNodeId(ctx, nodeID)
	if err != nil {
		return nil, fmt.Errorf("fetch connections for node %s: %w", nodeID, err)
	}

	f.cache.Connections.Set(nodeID, connections)
	return connections, nil
}
