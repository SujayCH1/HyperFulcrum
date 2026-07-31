package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

func (s *TopologyService) validateCreateTopology(
	ctx context.Context,
	projectID string,
	replicaID string,
	shardID string,
) error {

	// Future validations:
	// - Ensure the project exists.
	// - Ensure the project is inactive.
	// - Ensure both nodes exist.
	// - Ensure both nodes belong to the same project.
	// - Prevent self-replication.
	// - Ensure the primary node is not already a replica.
	// - Ensure the replica is not already assigned.
	// - Prevent replication cycles.
	// - Ensure node types are compatible.
	// - Ensure no duplicate topology exists.

	return nil
}

func (s *TopologyService) validateDeleteTopology(
	ctx context.Context,
	topology repository.NodeTopology,
) error {

	// Future validations:
	// - Ensure the project is inactive.
	// - Ensure the topology exists.
	// - Prevent deleting topology while replication is active.
	// - Prevent deleting topology while agents are using it.

	return nil
}
