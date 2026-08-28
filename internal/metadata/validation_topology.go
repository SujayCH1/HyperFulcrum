package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/repository"
)

func (s *TopologyService) validateCreateTopology(
	ctx context.Context,
	projectID string,
	shardID string,
	standbyID string,
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
	shards, loaded := s.cache.Shards.GetByProject(projectID)
	if !loaded {
		if err := s.refresher.RefreshShards(ctx, projectID); err != nil {
			return err
		}
		shards, _ = s.cache.Shards.GetByProject(projectID)
	}
	var shard repository.Shard
	shardFound := false
	for _, value := range shards {
		if value.ID == shardID {
			shard = value
			shardFound = true
			break
		}
	}
	if !shardFound {
		return sql.ErrNoRows
	}
	if shard.PrimaryNodeID == standbyID {
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

	var primary repository.Node
	var standby repository.Node
	primaryFound := false
	standbyFound := false
	for _, node := range nodes {
		if node.ID == shard.PrimaryNodeID {
			primary = node
			primaryFound = true
		}
		if node.ID == standbyID {
			standby = node
			standbyFound = true
		}
	}
	if !primaryFound || !standbyFound {
		return sql.ErrNoRows
	}
	if primary.ProjectID != projectID || standby.ProjectID != projectID {
		return sql.ErrNoRows
	}
	if primary.Role != repository.NodeRolePrimary || standby.Role != repository.NodeRoleStandby {
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
		if topology.ShardID == shardID && topology.StandbyNodeID == standbyID {
			return ErrDuplicateTopology
		}
		if topology.StandbyNodeID == standbyID {
			return ErrReplicaAlreadyUsed
		}
		if topology.StandbyNodeID == shard.PrimaryNodeID {
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
