package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/repository"
)

func (s *NodeService) validateAddNode(
	ctx context.Context,
	projectID string,
	nodeRole string,
	name string,
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
	if nodeRole != repository.NodeRolePrimary && nodeRole != repository.NodeRoleStandby && nodeRole != repository.NodeRoleUnassigned {
		return ErrInvalidNodeRole
	}

	nodes, loaded := s.cache.Nodes.GetByProject(projectID)
	if !loaded {
		err := s.refresher.RefreshNodes(ctx, projectID)
		if err != nil {
			return err
		}
		nodes, _ = s.cache.Nodes.GetByProject(projectID)
	}

	for _, node := range nodes {
		if node.Name == name {
			return ErrDuplicateNodeName
		}
	}

	// Project/node limits can be enforced when configured.

	return nil
}

func (s *NodeService) validateRemoveNode(
	ctx context.Context,
	node repository.Node,
) error {
	if node.Status {
		return ErrNodeActive
	}

	project, ok := s.cache.Projects.Get(node.ProjectID)
	if !ok {
		err := s.refresher.RefreshProject(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(node.ProjectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}

	topologies, loaded := s.cache.Topology.GetByProjectID(node.ProjectID)
	if !loaded {
		err := s.refresher.RefreshTopology(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		topologies, _ = s.cache.Topology.GetByProjectID(node.ProjectID)
	}

	for _, topology := range topologies {
		if topology.PrimaryNodeID == node.ID || topology.StandbyNodeID == node.ID {
			return ErrNodeInTopology
		}
	}

	shards, loaded := s.cache.Shards.GetByProject(node.ProjectID)
	if !loaded {
		if err := s.refresher.RefreshShards(ctx, node.ProjectID); err != nil {
			return err
		}
		shards, _ = s.cache.Shards.GetByProject(node.ProjectID)
	}
	for _, shard := range shards {
		if shard.PrimaryNodeID == node.ID {
			return ErrNodeOwnsShard
		}
	}

	// Agent assignments can be checked when agents are implemented.

	return nil
}

func (s *NodeService) validateUpdateNodeStatus(
	ctx context.Context,
	node repository.Node,
	status bool,
) error {
	if status {
		_, ok := s.cache.Connections.Get(node.ID)
		if !ok {
			err := s.refresher.RefreshConnections(ctx, node.ProjectID)
			if err != nil {
				return err
			}
			_, ok = s.cache.Connections.Get(node.ID)
		}
		if !ok {
			return ErrConnectionNotFound
		}
	}

	if !status {
		topologies, loaded := s.cache.Topology.GetByProjectID(node.ProjectID)
		if !loaded {
			err := s.refresher.RefreshTopology(ctx, node.ProjectID)
			if err != nil {
				return err
			}
			topologies, _ = s.cache.Topology.GetByProjectID(node.ProjectID)
		}
		for _, topology := range topologies {
			if topology.PrimaryNodeID == node.ID || topology.StandbyNodeID == node.ID {
				return ErrNodeInTopology
			}
		}
	}

	// Replication health can be checked when replication is implemented.

	return nil
}

func (s *NodeService) validateUpdateNodeName(
	ctx context.Context,
	node repository.Node,
	name string,
) error {
	project, ok := s.cache.Projects.Get(node.ProjectID)
	if !ok {
		err := s.refresher.RefreshProject(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(node.ProjectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}

	nodes, loaded := s.cache.Nodes.GetByProject(node.ProjectID)
	if !loaded {
		err := s.refresher.RefreshNodes(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		nodes, _ = s.cache.Nodes.GetByProject(node.ProjectID)
	}

	for _, existingNode := range nodes {
		if existingNode.ID != node.ID && existingNode.Name == name {
			return ErrDuplicateNodeName
		}
	}

	return nil
}

func (s *NodeService) validateUpdateNodeRole(
	ctx context.Context,
	node repository.Node,
	nodeRole string,
) error {
	if nodeRole != repository.NodeRolePrimary && nodeRole != repository.NodeRoleStandby && nodeRole != repository.NodeRoleUnassigned {
		return ErrInvalidNodeRole
	}

	project, ok := s.cache.Projects.Get(node.ProjectID)
	if !ok {
		err := s.refresher.RefreshProject(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(node.ProjectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}

	topologies, loaded := s.cache.Topology.GetByProjectID(node.ProjectID)
	if !loaded {
		err := s.refresher.RefreshTopology(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		topologies, _ = s.cache.Topology.GetByProjectID(node.ProjectID)
	}
	for _, topology := range topologies {
		if topology.PrimaryNodeID == node.ID || topology.StandbyNodeID == node.ID {
			return ErrNodeInTopology
		}
	}

	shards, loaded := s.cache.Shards.GetByProject(node.ProjectID)
	if !loaded {
		if err := s.refresher.RefreshShards(ctx, node.ProjectID); err != nil {
			return err
		}
		shards, _ = s.cache.Shards.GetByProject(node.ProjectID)
	}
	for _, shard := range shards {
		if shard.PrimaryNodeID == node.ID {
			return ErrNodeOwnsShard
		}
	}

	return nil
}
