package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type NodeConnectionService struct {
	repo      *repository.NodeConnectionRepository
	nodeRepo  *repository.NodeRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewNodeConnectionService(
	repo *repository.NodeConnectionRepository,
	nodeRepo *repository.NodeRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *NodeConnectionService {

	return &NodeConnectionService{
		repo:      repo,
		nodeRepo:  nodeRepo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *NodeConnectionService) AddConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	err := s.repo.ConnectionAdd(
		ctx,
		connection,
	)
	if err != nil {
		return err
	}

	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		connection.NodeId,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshConnections(
		ctx,
		node.ProjectID,
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

	err = s.repo.ConnectionRemove(
		ctx,
		nodeID,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshConnections(
		ctx,
		node.ProjectID,
	)
}

func (s *NodeConnectionService) UpdateConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	err := s.repo.ConnectionUpdate(
		ctx,
		connection,
	)
	if err != nil {
		return err
	}

	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		connection.NodeId,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshConnections(
		ctx,
		node.ProjectID,
	)
}

func (s *NodeConnectionService) GetConnectionByNodeID(
	ctx context.Context,
	nodeID string,
) (repository.NodeConnection, error) {

	if conn, ok := s.cache.Connections.Get(nodeID); ok {
		return conn, nil
	}

	node, err := s.nodeRepo.NodeGetByID(
		ctx,
		nodeID,
	)
	if err != nil {
		return repository.NodeConnection{}, err
	}

	if err := s.refresher.RefreshConnections(
		ctx,
		node.ProjectID,
	); err != nil {
		return repository.NodeConnection{}, err
	}

	conn, ok := s.cache.Connections.Get(nodeID)
	if !ok {
		return repository.NodeConnection{}, sql.ErrNoRows
	}

	return conn, nil
}
