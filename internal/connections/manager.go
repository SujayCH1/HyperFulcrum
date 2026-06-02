package connections

import (
	"context"
	"hyperfulcrum/internal/repository"
)

type ConnectionManager struct {
	store        *ConnectionStore
	projectRepo  *repository.ProjectRepository
	nodeRepo     *repository.NodeRepository
	nodeConnRepo *repository.NodeConnectionRepository
}

func NewConnectionManager(
	store *ConnectionStore, projectRepo *repository.ProjectRepository,
	nodeRepo *repository.NodeRepository, nodeConnRepo *repository.NodeConnectionRepository) *ConnectionManager {
	return &ConnectionManager{
		store:        store,
		projectRepo:  projectRepo,
		nodeRepo:     nodeRepo,
		nodeConnRepo: nodeConnRepo,
	}
}

func (m *ConnectionManager) InitiateActiveConnections(ctx context.Context) error {

	projects, err := m.projectRepo.ProjectGetReady(ctx)
	if err != nil {
		return err
	}

	for _, project := range projects {
		nodes, err := m.nodeRepo.NodeList(ctx, project.ID)
		if err != nil {
			return err
		}

		for _, node := range nodes {

			connInfo, err := m.nodeConnRepo.GetConnectionByNodeId(ctx, node.ID)
			if err != nil {
				continue
			}

			dsn := buildDSN(*connInfo)

			db, err := NewConnection(ctx, dsn)
			if err != nil {
				continue
			}
			m.store.Set(project.ID, node.ID, db)
		}
	}
	return nil
}

func (m *ConnectionManager) InititateConnectionsAll(ctx context.Context) error {

	projects, err := m.projectRepo.ProjectList(ctx)
	if err != nil {
		return err
	}

	for _, project := range projects {

		nodes, err := m.nodeRepo.NodeList(ctx, project.ID)
		if err != nil {
			return err
		}

		for _, node := range nodes {

			connInfo, err := m.nodeConnRepo.GetConnectionByNodeId(ctx, node.ID)
			if err != nil {
				continue
			}

			dsn := buildDSN(*connInfo)

			db, err := NewConnection(ctx, dsn)
			if err != nil {
				continue
			}

			m.store.Set(project.ID, node.ID, db)
		}
	}

	return nil
}

func (m *ConnectionManager) CheckConnectionHealth(
	ctx context.Context,
	projectID string,
	nodeID string,
) (bool, error) {

	conn, err := m.store.Get(projectID, nodeID)
	if err != nil {
		// _ = m.nodeRepo.NodeUpdateStatus(ctx, nodeID, false)
		return false, err
	}

	if err := conn.PingContext(ctx); err != nil {
		// _ = m.nodeRepo.NodeUpdateStatus(ctx, nodeID, false)
		return false, err
	}
	//_ = m.nodeRepo.NodeUpdateStatus(ctx, nodeID, true)

	return true, nil
}
