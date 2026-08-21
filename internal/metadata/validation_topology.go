package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/repository"
)

func (s *TopologyService) validateCreateTopology(
	ctx context.Context,
	projectID string,
	replicaID string,
	shardID string,
) error {
	project, ok := s.cache.Projects.Get(projectID)
	if !ok {
		err := s.refresher.RefreshProject(ctx, projectID)
		if err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(projectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}
	if shardID == replicaID {
		return ErrTopologySelfRelation
	}

	nodes, loaded := s.cache.Nodes.GetByProject(projectID)
	if !loaded {
		err := s.refresher.RefreshNodes(ctx, projectID)
		if err != nil {
			return err
		}
		nodes, _ = s.cache.Nodes.GetByProject(projectID)
	}

	var shard repository.Node
	var replica repository.Node
	shardFound := false
	replicaFound := false
	for _, node := range nodes {
		if node.ID == shardID {
			shard = node
			shardFound = true
		}
		if node.ID == replicaID {
			replica = node
			replicaFound = true
		}
	}
	if !shardFound || !replicaFound {
		return sql.ErrNoRows
	}
	if shard.ProjectID != projectID || replica.ProjectID != projectID {
		return sql.ErrNoRows
	}
	if shard.Type != "shard" || replica.Type != "replica" {
		return ErrTopologyRoleMismatch
	}

	topologies, loaded := s.cache.Topology.GetByProjectID(projectID)
	if !loaded {
		err := s.refresher.RefreshTopology(ctx, projectID)
		if err != nil {
			return err
		}
		topologies, _ = s.cache.Topology.GetByProjectID(projectID)
	}

	for _, topology := range topologies {
		if topology.ShardNodeID == shardID && topology.ReplicaNodeID == replicaID {
			return ErrDuplicateTopology
		}
		if topology.ReplicaNodeID == replicaID {
			return ErrReplicaAlreadyUsed
		}
		if topology.ReplicaNodeID == shardID {
			return ErrShardIsReplica
		}
	}

	return nil
}

func (s *TopologyService) validateDeleteTopology(
	ctx context.Context,
	topology repository.NodeTopology,
) error {
	project, ok := s.cache.Projects.Get(topology.ProjectID)
	if !ok {
		err := s.refresher.RefreshProject(ctx, topology.ProjectID)
		if err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(topology.ProjectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}

	// Existence is established by GetTopologyByID before this validation.
	// Deferred until the corresponding state exists:
	// - Prevent deleting topology while replication is active.
	// - Prevent deleting topology while agents are using it.

	return nil
}
