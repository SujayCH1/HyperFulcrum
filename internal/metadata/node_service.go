package metadata

import (
	"context"
	"errors"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/connections"
	"hyperfulcrum/internal/repository"
)

type NodeService struct {
	repo              *repository.NodeRepository
	cache             *cache.CacheManager
	refresher         *cache.CacheRefresher
	connectionManager *connections.ConnectionManager
}

func NewNodeService(
	repo *repository.NodeRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
	connectionManager *connections.ConnectionManager,
) *NodeService {

	return &NodeService{
		repo:              repo,
		cache:             cache,
		refresher:         refresher,
		connectionManager: connectionManager,
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

	nodes, loaded := s.cache.Nodes.GetByProject(projectID)
	if loaded {
		return nodes, nil
	}

	if err := s.refresher.RefreshNodes(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	nodes, _ = s.cache.Nodes.GetByProject(projectID)

	return nodes, nil
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

	return errors.Join(
		s.connectionManager.RemoveNode(node.ProjectID, nodeID),
		s.refresher.RefreshNodes(
			ctx,
			node.ProjectID,
		),
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

	return errors.Join(
		s.refresher.RefreshNodes(
			ctx,
			node.ProjectID,
		),
		s.connectionManager.SyncNode(
			ctx,
			node.ProjectID,
			nodeID,
		),
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
