package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

type NodeService struct {
	repo *repository.NodeRepository
}

func NewNodeService(
	repo *repository.NodeRepository,
) *NodeService {

	return &NodeService{
		repo: repo,
	}
}

func (s *NodeService) AddNode(
	ctx context.Context,
	projectID string,
	nodeType string,
	name string,
) error {

	return s.repo.NodeAdd(
		ctx,
		projectID,
		nodeType,
		name,
	)
}

func (s *NodeService) ListNodes(
	ctx context.Context,
	projectID string,
) ([]repository.Node, error) {

	return s.repo.NodeList(
		ctx,
		projectID,
	)
}

func (s *NodeService) RemoveNode(
	ctx context.Context,
	nodeID string,
) error {

	return s.repo.NodeRemove(
		ctx,
		nodeID,
	)
}

func (s *NodeService) UpdateNodeName(
	ctx context.Context,
	nodeID string,
	name string,
) error {

	return s.repo.NodeUpdateName(
		ctx,
		nodeID,
		name,
	)
}

func (s *NodeService) UpdateNodeStatus(
	ctx context.Context,
	nodeID string,
	status bool,
) error {

	return s.repo.NodeUpdateStatus(
		ctx,
		nodeID,
		status,
	)
}

func (s *NodeService) UpdateNodeType(
	ctx context.Context,
	nodeID string,
	nodeType string,
) error {

	return s.repo.NodeUpdateType(
		ctx,
		nodeID,
		nodeType,
	)
}
