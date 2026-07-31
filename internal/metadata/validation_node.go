package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

func (s *NodeService) validateAddNode(
	ctx context.Context,
	projectID string,
	nodeType string,
	name string,
) error {

	// Future validations:
	// - Ensure the project exists.
	// - Ensure the project is inactive before modifying topology.
	// - Ensure node type is valid.
	// - Ensure node name is unique within the project.
	// - Enforce project/node limits if introduced.

	return nil
}

func (s *NodeService) validateRemoveNode(
	ctx context.Context,
	node repository.Node,
) error {

	// Future validations:
	// - Ensure the node is inactive.
	// - Ensure the node is not a primary node.
	// - Ensure no replicas depend on this node.
	// - Ensure the node is not part of an active replication relationship.
	// - Ensure the project is not currently active.
	// - Ensure no agents are currently assigned.

	return nil
}

func (s *NodeService) validateUpdateNodeStatus(
	ctx context.Context,
	node repository.Node,
	status bool,
) error {

	// Future validations:
	// - Prevent deactivating a required primary.
	// - Prevent activating a node without a valid connection.
	// - Ensure replication health before changing status.
	// - Ensure topology remains valid after the status change.

	return nil
}

func (s *NodeService) validateUpdateNodeType(
	ctx context.Context,
	node repository.Node,
	nodeType string,
) error {

	// Future validations:
	// - Ensure the requested node type is valid.
	// - Prevent invalid primary/replica transitions.
	// - Ensure replication topology remains consistent.
	// - Prevent changing type while the project is active.

	return nil
}
