package metadata

import (
	"context"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type NodeService struct {
	repo      *repository.NodeRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewNodeService(
	repo *repository.NodeRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *NodeService {

	return &NodeService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *NodeService) AddNode(
	ctx context.Context,
	projectID string,
	nodeType string,
	name string,
) error {

	err := s.repo.NodeAdd(
		ctx,
		projectID,
		nodeType,
		name,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		projectID,
	)
}

func (s *NodeService) ListNodes(
	ctx context.Context,
	projectID string,
) ([]repository.Node, error) {

	cachedNodes := s.cache.Nodes.GetAll()

	var projectNodes []repository.Node

	for _, node := range cachedNodes {
		if node.ProjectID == projectID {
			projectNodes = append(projectNodes, node)
		}
	}

	if len(projectNodes) > 0 {
		return projectNodes, nil
	}

	if err := s.refresher.RefreshNodes(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	cachedNodes = s.cache.Nodes.GetAll()

	projectNodes = nil

	for _, node := range cachedNodes {
		if node.ProjectID == projectID {
			projectNodes = append(projectNodes, node)
		}
	}

	return projectNodes, nil
}

func (s *NodeService) RemoveNode(
	ctx context.Context,
	nodeID string,
) error {

	node, err := s.repo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	err = s.repo.NodeRemove(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		node.ProjectID,
	)
}

func (s *NodeService) UpdateNodeName(
	ctx context.Context,
	nodeID string,
	name string,
) error {

	node, err := s.repo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	err = s.repo.NodeUpdateName(
		ctx,
		nodeID,
		name,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		node.ProjectID,
	)
}

func (s *NodeService) UpdateNodeStatus(
	ctx context.Context,
	nodeID string,
	status bool,
) error {

	node, err := s.repo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	err = s.repo.NodeUpdateStatus(
		ctx,
		nodeID,
		status,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		node.ProjectID,
	)
}

func (s *NodeService) UpdateNodeType(
	ctx context.Context,
	nodeID string,
	nodeType string,
) error {

	node, err := s.repo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	err = s.repo.NodeUpdateType(
		ctx,
		nodeID,
		nodeType,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		node.ProjectID,
	)
}
