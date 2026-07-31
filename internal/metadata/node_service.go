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

	if err := s.validateAddNode(
		ctx,
		projectID,
		nodeType,
		name,
	); err != nil {
		return err
	}

	if err := s.repo.NodeAdd(
		ctx,
		projectID,
		nodeType,
		name,
	); err != nil {
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

	nodes := s.cache.Nodes.GetByProject(projectID)
	if len(nodes) > 0 {
		return nodes, nil
	}

	if err := s.refresher.RefreshNodes(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	return s.cache.Nodes.GetByProject(projectID), nil
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

	if err := s.validateRemoveNode(
		ctx,
		node,
	); err != nil {
		return err
	}

	if err := s.repo.NodeRemove(
		ctx,
		nodeID,
	); err != nil {
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

	if err := s.repo.NodeUpdateName(
		ctx,
		nodeID,
		name,
	); err != nil {
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

	if err := s.validateUpdateNodeStatus(
		ctx,
		node,
		status,
	); err != nil {
		return err
	}

	if err := s.repo.NodeUpdateStatus(
		ctx,
		nodeID,
		status,
	); err != nil {
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

	if err := s.validateUpdateNodeType(
		ctx,
		node,
		nodeType,
	); err != nil {
		return err
	}

	if err := s.repo.NodeUpdateType(
		ctx,
		nodeID,
		nodeType,
	); err != nil {
		return err
	}

	return s.refresher.RefreshNodes(
		ctx,
		node.ProjectID,
	)
}
