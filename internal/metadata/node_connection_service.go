package metadata

import (
	"context"
	"database/sql"
	"errors"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/connections"
	"hyperfulcrum/internal/repository"
)

type NodeConnectionService struct {
	repo              *repository.NodeConnectionRepository
	nodeRepo          *repository.NodeRepository
	cache             *cache.CacheManager
	refresher         *cache.CacheRefresher
	connectionManager *connections.ConnectionManager
}

func NewNodeConnectionService(
	repo *repository.NodeConnectionRepository,
	nodeRepo *repository.NodeRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
	connectionManager *connections.ConnectionManager,
) *NodeConnectionService {

	return &NodeConnectionService{
		repo:              repo,
		nodeRepo:          nodeRepo,
		cache:             cache,
		refresher:         refresher,
		connectionManager: connectionManager,
	}
}

func (s *NodeConnectionService) AddConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {
	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		connection.NodeId,
	)
	if err != nil {
		return err
	}

	if err := s.validateAddConnection(
		ctx,
		connection,
	); err != nil {
		return err
	}

	if err := s.repo.ConnectionAdd(
		ctx,
		connection,
	); err != nil {
		return err
	}

	return errors.Join(
		s.refresher.RefreshConnections(
			ctx,
			node.ProjectID,
		),
		s.connectionManager.SyncNode(
			ctx,
			node.ProjectID,
			connection.NodeId,
		),
	)
}

func (s *NodeConnectionService) RemoveConnection(
	ctx context.Context,
	nodeID string,
) error {

	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	if err := s.validateRemoveConnection(
		ctx,
		node,
	); err != nil {
		return err
	}

	if err := s.repo.ConnectionRemove(
		ctx,
		nodeID,
	); err != nil {
		return err
	}

	return errors.Join(
		s.refresher.RefreshConnections(
			ctx,
			node.ProjectID,
		),
		s.connectionManager.RemoveNode(node.ProjectID, nodeID),
	)
}

func (s *NodeConnectionService) UpdateConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		connection.NodeId,
	)
	if err != nil {
		return err
	}

	if err := s.validateUpdateConnection(
		ctx,
		node,
		connection,
	); err != nil {
		return err
	}

	if err := s.repo.ConnectionUpdate(
		ctx,
		connection,
	); err != nil {
		return err
	}

	return errors.Join(
		s.refresher.RefreshConnections(
			ctx,
			node.ProjectID,
		),
		s.connectionManager.SyncNode(
			ctx,
			node.ProjectID,
			connection.NodeId,
		),
	)
}

func (s *NodeConnectionService) GetConnectionByNodeID(
	ctx context.Context,
	nodeID string,
) (repository.NodeConnection, error) {

	if conn, ok := s.cache.Connections.Get(nodeID); ok {
		return conn, nil
	}

	projectID, ok := s.cache.Nodes.GetProjectID(nodeID)
	if !ok {
		node, err := s.nodeRepo.NodeGetByID(
			ctx,
			nodeID,
		)
		if err != nil {
			return repository.NodeConnection{}, err
		}

		projectID = node.ProjectID

		if err := s.refresher.RefreshNodes(
			ctx,
			projectID,
		); err != nil {
			return repository.NodeConnection{}, err
		}

		projectID, ok = s.cache.Nodes.GetProjectID(nodeID)
		if !ok {
			return repository.NodeConnection{}, sql.ErrNoRows
		}
	}

	if s.cache.Connections.ProjectLoaded(projectID) {
		return repository.NodeConnection{}, sql.ErrNoRows
	}

	if err := s.refresher.RefreshConnections(
		ctx,
		projectID,
	); err != nil {
		return repository.NodeConnection{}, err
	}

	conn, ok := s.cache.Connections.Get(nodeID)
	if !ok {
		return repository.NodeConnection{}, sql.ErrNoRows
	}

	return conn, nil
}
