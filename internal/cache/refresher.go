package cache

import (
	"context"
	"hyperfulcrum/internal/repository"
)

type CacheRefresher struct {
	projectRepo    *repository.ProjectRepository
	nodeRepo       *repository.NodeRepository
	connectionRepo *repository.NodeConnectionRepository
	topologyRepo   *repository.NodeTopologyRepository

	cache *CacheManager
}

func NewCacheRefresher(
	projectRepo *repository.ProjectRepository,
	nodeRepo *repository.NodeRepository,
	connectionRepo *repository.NodeConnectionRepository,
	topologyRepo *repository.NodeTopologyRepository,
	cache *CacheManager,
) *CacheRefresher {

	return &CacheRefresher{
		projectRepo:    projectRepo,
		nodeRepo:       nodeRepo,
		connectionRepo: connectionRepo,
		topologyRepo:   topologyRepo,
		cache:          cache,
	}
}

func (r *CacheRefresher) RefreshProjectMetadata(
	ctx context.Context,
	projectID string,
) error {

	if err := r.RefreshProject(ctx, projectID); err != nil {
		return err
	}

	if err := r.RefreshNodes(ctx, projectID); err != nil {
		return err
	}

	if err := r.RefreshConnections(ctx, projectID); err != nil {
		return err
	}

	if err := r.RefreshTopology(ctx, projectID); err != nil {
		return err
	}

	return nil
}

func (r *CacheRefresher) RefreshProject(
	ctx context.Context,
	projectID string,
) error {

	project, err := r.projectRepo.ProjectGetByID(ctx, projectID)
	if err != nil {
		return err
	}

	r.cache.Projects.Set(project)

	return nil
}

func (r *CacheRefresher) RefreshNodes(
	ctx context.Context,
	projectID string,
) error {

	nodes, err := r.nodeRepo.NodesGetByPorjectID(ctx, projectID)
	if err != nil {
		return err
	}

	r.cache.ProjectNodes.Delete(projectID)

	for _, node := range nodes {

		r.cache.Nodes.Set(node)

		r.cache.ProjectNodes.Add(projectID, node.ID)
	}

	return nil
}

func (r *CacheRefresher) RefreshConnections(
	ctx context.Context,
	projectID string,
) error {

	nodes, err := r.nodeRepo.NodesGetByPorjectID(ctx, projectID)
	if err != nil {
		return err
	}

	for _, node := range nodes {

		conn, err := r.connectionRepo.GetConnectionByNodeId(
			ctx,
			node.ID,
		)

		if err != nil {
			continue
		}

		r.cache.Connections.Set(conn)
	}

	return nil
}

func (r *CacheRefresher) RefreshTopology(
	ctx context.Context,
	projectID string,
) error {

	topology, err := r.topologyRepo.TopologyGetAll(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Topology.Set(projectID, topology)

	return nil
}

func (r *CacheRefresher) RefreshAllProjects(
	ctx context.Context,
) error {

	projects, err := r.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	for _, project := range projects {

		err := r.RefreshProjectMetadata(
			ctx,
			project.ID,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
