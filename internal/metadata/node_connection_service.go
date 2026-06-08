package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

type NodeConnectionService struct {
	repo *repository.NodeConnectionRepository
}

func NewNodeConnectionService(
	repo *repository.NodeConnectionRepository,
) *NodeConnectionService {

	return &NodeConnectionService{
		repo: repo,
	}
}

func (s *NodeConnectionService) AddConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	return s.repo.ConnectionAdd(
		ctx,
		connection,
	)
}

func (s *NodeConnectionService) RemoveConnection(
	ctx context.Context,
	nodeID string,
) error {

	return s.repo.ConnectionRemove(
		ctx,
		nodeID,
	)
}

func (s *NodeConnectionService) UpdateConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	return s.repo.ConnectionUpdate(
		ctx,
		connection,
	)
}

func (s *NodeConnectionService) GetConnectionByNodeID(
	ctx context.Context,
	nodeID string,
) (repository.NodeConnection, error) {

	return s.repo.GetConnectionByNodeId(
		ctx,
		nodeID,
	)
}
