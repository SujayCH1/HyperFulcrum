package connections

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"hyperfulcrum/internal/repository"
)

var (
	ErrActiveProjectHasNoNodes = errors.New("active project has no nodes")
	ErrActiveProjectNode       = errors.New("active project contains an inactive node")
	ErrConnectionConfig        = errors.New("node connection configuration not found")
	ErrNodeProjectMismatch     = errors.New("node does not belong to project")
)

type projectRepository interface {
	ProjectGetRunning(ctx context.Context) (repository.Project, error)
	ProjectGetByID(ctx context.Context, id string) (repository.Project, error)
}

type nodeRepository interface {
	NodeList(ctx context.Context, projectID string) ([]repository.Node, error)
	NodeGetByID(ctx context.Context, nodeID string) (repository.Node, error)
}

type nodeConnectionRepository interface {
	GetConnectionByNodeId(ctx context.Context, nodeID string) (repository.NodeConnection, error)
	ConnectionsListByProjectID(ctx context.Context, projectID string) ([]repository.NodeConnection, error)
}

type ConnectionManager struct {
	store        *PoolStore
	projectRepo  projectRepository
	nodeRepo     nodeRepository
	nodeConnRepo nodeConnectionRepository
	connect      func(context.Context, string) (*sql.DB, error)
}

func NewConnectionManager(
	store *PoolStore,
	projectRepo projectRepository,
	nodeRepo nodeRepository,
	nodeConnRepo nodeConnectionRepository,
) *ConnectionManager {
	return &ConnectionManager{
		store:        store,
		projectRepo:  projectRepo,
		nodeRepo:     nodeRepo,
		nodeConnRepo: nodeConnRepo,
		connect:      NewConnection,
	}
}

func (m *ConnectionManager) InitializeActiveConnections(ctx context.Context) error {
	return m.SyncActiveProject(ctx)
}

func (m *ConnectionManager) SyncActiveProject(ctx context.Context) error {
	pools, err := m.buildActivePools(ctx)
	if err != nil {
		return err
	}

	return m.store.ReplaceAll(pools)
}

func (m *ConnectionManager) buildActivePools(
	ctx context.Context,
) (map[string]map[string]*sql.DB, error) {
	pools := make(map[string]map[string]*sql.DB)

	project, err := m.projectRepo.ProjectGetRunning(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return pools, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get active project: %w", err)
	}

	nodes, err := m.nodeRepo.NodeList(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("list nodes for active project %s: %w", project.ID, err)
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("project %s: %w", project.ID, ErrActiveProjectHasNoNodes)
	}

	connections, err := m.nodeConnRepo.ConnectionsListByProjectID(ctx, project.ID)
	if err != nil {
		return nil, fmt.Errorf("list connections for active project %s: %w", project.ID, err)
	}

	connectionByNode := make(map[string]repository.NodeConnection, len(connections))
	for _, connection := range connections {
		connectionByNode[connection.NodeId] = connection
	}

	projectPools := make(map[string]*sql.DB, len(nodes))
	pools[project.ID] = projectPools

	for _, node := range nodes {
		if !node.Status {
			return nil, joinPoolError(
				fmt.Errorf("project %s node %s: %w", project.ID, node.ID, ErrActiveProjectNode),
				projectPools,
			)
		}

		connection, ok := connectionByNode[node.ID]
		if !ok {
			return nil, joinPoolError(
				fmt.Errorf("project %s node %s: %w", project.ID, node.ID, ErrConnectionConfig),
				projectPools,
			)
		}

		db, err := m.connect(ctx, buildDSN(connection))
		if err != nil {
			return nil, joinPoolError(
				fmt.Errorf("connect project %s node %s: %w", project.ID, node.ID, err),
				projectPools,
			)
		}

		projectPools[node.ID] = db
	}

	return pools, nil
}

func (m *ConnectionManager) SyncNode(
	ctx context.Context,
	projectID string,
	nodeID string,
) error {
	project, err := m.projectRepo.ProjectGetByID(ctx, projectID)
	if err != nil {
		return err
	}

	node, err := m.nodeRepo.NodeGetByID(ctx, nodeID)
	if err != nil {
		return err
	}

	if node.ProjectID != projectID {
		return ErrNodeProjectMismatch
	}

	if !project.Running || !node.Status {
		return m.store.Remove(projectID, nodeID)
	}

	connection, err := m.nodeConnRepo.GetConnectionByNodeId(ctx, nodeID)
	if err != nil {
		removeErr := m.store.Remove(projectID, nodeID)
		return errors.Join(
			fmt.Errorf(
				"node %s: %w",
				nodeID,
				errors.Join(ErrConnectionConfig, err),
			),
			removeErr,
		)
	}

	db, err := m.connect(ctx, buildDSN(connection))
	if err != nil {
		removeErr := m.store.Remove(projectID, nodeID)
		return errors.Join(
			fmt.Errorf("connect project %s node %s: %w", projectID, nodeID, err),
			removeErr,
		)
	}

	return m.store.Set(projectID, nodeID, db)
}

func (m *ConnectionManager) RemoveNode(projectID string, nodeID string) error {
	return m.store.Remove(projectID, nodeID)
}

func (m *ConnectionManager) RemoveProject(projectID string) error {
	return m.store.RemoveProject(projectID)
}

func (m *ConnectionManager) CheckConnectionHealth(
	ctx context.Context,
	projectID string,
	nodeID string,
) (bool, error) {
	db, err := m.store.Get(projectID, nodeID)
	if err != nil {
		return false, err
	}

	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(healthCtx); err != nil {
		removeErr := m.store.Remove(projectID, nodeID)
		return false, errors.Join(err, removeErr)
	}

	return true, nil
}

func (m *ConnectionManager) Close() error {
	return m.store.Close()
}

func joinPoolError(err error, pools map[string]*sql.DB) error {
	return errors.Join(err, closePools(pools))
}
