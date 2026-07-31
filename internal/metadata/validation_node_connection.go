package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

func (s *NodeConnectionService) validateAddConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {

	// Future validations:
	// - Ensure the node exists.
	// - Ensure the node does not already have a connection.
	// - Ensure the project is inactive before modifying connections.
	// - Validate connection parameters if business rules require it.
	// - Ensure the node is not currently participating in replication.

	return nil
}

func (s *NodeConnectionService) validateRemoveConnection(
	ctx context.Context,
	node repository.Node,
) error {

	// Future validations:
	// - Ensure the node has an existing connection.
	// - Prevent removing the connection from an active node.
	// - Prevent removing the connection while replication exists.
	// - Prevent removing the connection while the project is active.

	return nil
}

func (s *NodeConnectionService) validateUpdateConnection(
	ctx context.Context,
	node repository.Node,
	connection *repository.NodeConnection,
) error {

	// Future validations:
	// - Ensure the node has an existing connection.
	// - Prevent updating while the node is active.
	// - Prevent updating while replication depends on this node.
	// - Validate any immutable connection fields if required.

	return nil
}
