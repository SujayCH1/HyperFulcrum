package replication

import (
	"context"

	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/repository"
)

type ReplicationService struct {
	topology *metadata.TopologyService
	nodes    *metadata.NodeService
}

func NewReplicationService(
	topology *metadata.TopologyService,
	nodes *metadata.NodeService,
) *ReplicationService {

	return &ReplicationService{
		topology: topology,
		nodes:    nodes,
	}
}

func (s *ReplicationService) CreateReplication(
	ctx context.Context,
	projectID string,
	shardID string,
	replicaID string,
) (repository.NodeTopology, error) {

	// TODO:
	// 1. Validate request
	// 2. Configure PostgreSQL replication
	// 3. Create topology metadata

	return repository.NodeTopology{}, nil
}

func (s *ReplicationService) DeleteReplication(
	ctx context.Context,
	relationID string,
	projectID string,
) error {

	// TODO:
	// 1. Validate request
	// 2. Stop PostgreSQL replication
	// 3. Remove topology metadata

	return nil
}

func (s *ReplicationService) PromoteReplica(
	ctx context.Context,
	relationID string,
	shardID string,
	replicaID string,
) error {

	// TODO:
	// 1. Validate request
	// 2. Promote replica inside PostgreSQL
	// 3. Update topology
	// 4. Update node types

	return nil
}
