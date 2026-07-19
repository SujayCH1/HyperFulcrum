package cache

import (
	"context"
	"database/sql"

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

func (r *CacheRefresher) RefreshProjects(
	ctx context.Context,
) error {

	projects, err := r.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	r.cache.Projects.Clear()

	for _, project := range projects {
		r.cache.Projects.Set(project)
	}

	return nil
}

func (r *CacheRefresher) RefreshProject(
	ctx context.Context,
	projectID string,
) error {

	project, err := r.projectRepo.ProjectGetByID(
		ctx,
		projectID,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			r.cache.Projects.Delete(projectID)
			return nil
		}

		return err
	}

	r.cache.Projects.Set(project)

	return nil
}

func (r *CacheRefresher) RefreshNodes(
	ctx context.Context,
	projectID string,
) error {

	nodes, err := r.nodeRepo.NodeList(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		r.cache.Nodes.Delete(node.ID)
	}

	for _, node := range nodes {
		r.cache.Nodes.Set(node)
	}

	return nil
}

func (r *CacheRefresher) RefreshConnections(
	ctx context.Context,
	projectID string,
) error {

	nodes, err := r.nodeRepo.NodeList(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		r.cache.Connections.Delete(node.ID)
	}

	for _, node := range nodes {

		conn, err := r.connectionRepo.GetConnectionByNodeId(
			ctx,
			node.ID,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			return err
		}

		r.cache.Connections.Set(conn)
	}

	return nil
}

func (r *CacheRefresher) RefreshTopology(
	ctx context.Context,
	projectID string,
) error {

	topologies, err := r.topologyRepo.TopologyGetAll(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	r.cache.Topology.Set(
		projectID,
		topologies,
	)

	return nil
}

func (r *CacheRefresher) RefreshAllProjects(
	ctx context.Context,
) error {

	projects, err := r.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	r.cache.Projects.Clear()

	for _, project := range projects {

		r.cache.Projects.Set(project)

		if err := r.RefreshNodes(
			ctx,
			project.ID,
		); err != nil {
			return err
		}

		if err := r.RefreshConnections(
			ctx,
			project.ID,
		); err != nil {
			return err
		}

		if err := r.RefreshTopology(
			ctx,
			project.ID,
		); err != nil {
			return err
		}
	}

	return nil
}
