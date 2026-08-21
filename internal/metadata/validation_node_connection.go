package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/repository"
)

func (s *NodeConnectionService) validateAddConnection(
	ctx context.Context,
	connection *repository.NodeConnection,
) error {
	node, err := s.nodeRepo.NodeGetByID(ctx, connection.NodeId)
	if err != nil {
		return err
	}

	project, ok := s.cache.Projects.Get(node.ProjectID)
	if !ok {
		err = s.refresher.RefreshProject(ctx, node.ProjectID)
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

	_, ok = s.cache.Connections.Get(connection.NodeId)
	if !ok && !s.cache.Connections.ProjectLoaded(node.ProjectID) {
		err = s.refresher.RefreshConnections(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		_, ok = s.cache.Connections.Get(connection.NodeId)
	}
	if ok {
		return ErrConnectionExists
	}

	err = validateConnectionFields(connection)
	if err != nil {
		return err
	}

	return nil
}

func (s *NodeConnectionService) validateRemoveConnection(
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

	_, ok = s.cache.Connections.Get(node.ID)
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

	err := ensureNodeOutsideTopology(ctx, s.refresher, s.cache, node)
	if err != nil {
		return err
	}

	return nil
}

func (s *NodeConnectionService) validateUpdateConnection(
	ctx context.Context,
	node repository.Node,
	connection *repository.NodeConnection,
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

	_, ok = s.cache.Connections.Get(node.ID)
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

	err := validateConnectionFields(connection)
	if err != nil {
		return err
	}

	err = ensureNodeOutsideTopology(ctx, s.refresher, s.cache, node)
	if err != nil {
		return err
	}

	return nil
}

func validateConnectionFields(connection *repository.NodeConnection) error {
	if connection.Host == "" ||
		connection.Port < 1 ||
		connection.Port > 65535 ||
		connection.DatabaseName == "" ||
		connection.Username == "" ||
		connection.Password == "" {
		return ErrInvalidConnection
	}

	return nil
}
