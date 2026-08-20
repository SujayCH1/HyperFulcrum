package connections

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"hyperfulcrum/internal/repository"
)

type projectRepoStub struct {
	running    repository.Project
	runningErr error
	project    repository.Project
	projectErr error
}

func (s *projectRepoStub) ProjectGetRunning(ctx context.Context) (repository.Project, error) {
	return s.running, s.runningErr
}

func (s *projectRepoStub) ProjectGetByID(ctx context.Context, id string) (repository.Project, error) {
	return s.project, s.projectErr
}

type nodeRepoStub struct {
	nodes   []repository.Node
	listErr error
	node    repository.Node
	getErr  error
}

func (s *nodeRepoStub) NodeList(ctx context.Context, projectID string) ([]repository.Node, error) {
	return s.nodes, s.listErr
}

func (s *nodeRepoStub) NodeGetByID(ctx context.Context, nodeID string) (repository.Node, error) {
	return s.node, s.getErr
}

type connectionRepoStub struct {
	connections []repository.NodeConnection
	listErr     error
	connection  repository.NodeConnection
	getErr      error
}

func (s *connectionRepoStub) ConnectionsListByProjectID(
	ctx context.Context,
	projectID string,
) ([]repository.NodeConnection, error) {
	return s.connections, s.listErr
}

func (s *connectionRepoStub) GetConnectionByNodeId(
	ctx context.Context,
	nodeID string,
) (repository.NodeConnection, error) {
	return s.connection, s.getErr
}

func TestInitializeActiveConnectionsReplacesPoolsAtomically(t *testing.T) {
	store := NewPoolStore()
	oldDB := newTestPool(t)
	if err := store.Set("old-project", "old-node", oldDB); err != nil {
		t.Fatal(err)
	}

	projectRepo := &projectRepoStub{
		running: repository.Project{ID: "project-1", Running: true},
	}
	nodeRepo := &nodeRepoStub{
		nodes: []repository.Node{
			{ID: "node-1", ProjectID: "project-1", Status: true},
			{ID: "node-2", ProjectID: "project-1", Status: true},
		},
	}
	connectionRepo := &connectionRepoStub{
		connections: []repository.NodeConnection{
			{NodeId: "node-1"},
			{NodeId: "node-2"},
		},
	}

	manager := NewConnectionManager(store, projectRepo, nodeRepo, connectionRepo)
	manager.connect = func(ctx context.Context, dsn string) (*sql.DB, error) {
		return newTestPool(t), nil
	}
	defer manager.Close()

	if err := manager.InitializeActiveConnections(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !poolClosed(oldDB) {
		t.Fatal("old pool store was not closed")
	}

	if _, err := store.Get("project-1", "node-1"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.Get("project-1", "node-2"); err != nil {
		t.Fatal(err)
	}
}

func TestInitializeActiveConnectionsKeepsOldPoolsOnFailure(t *testing.T) {
	store := NewPoolStore()
	oldDB := newTestPool(t)
	if err := store.Set("old-project", "old-node", oldDB); err != nil {
		t.Fatal(err)
	}

	manager := NewConnectionManager(
		store,
		&projectRepoStub{running: repository.Project{ID: "project-1", Running: true}},
		&nodeRepoStub{nodes: []repository.Node{
			{ID: "node-1", ProjectID: "project-1", Status: true},
			{ID: "node-2", ProjectID: "project-1", Status: true},
		}},
		&connectionRepoStub{connections: []repository.NodeConnection{
			{NodeId: "node-1"},
			{NodeId: "node-2"},
		}},
	)

	var stagedDB *sql.DB
	connectionCount := 0
	manager.connect = func(ctx context.Context, dsn string) (*sql.DB, error) {
		connectionCount++
		if connectionCount == 2 {
			return nil, errors.New("connection failed")
		}

		stagedDB = newTestPool(t)
		return stagedDB, nil
	}
	defer manager.Close()

	if err := manager.InitializeActiveConnections(context.Background()); err == nil {
		t.Fatal("expected initialization error")
	}

	db, err := store.Get("old-project", "old-node")
	if err != nil || db != oldDB {
		t.Fatal("old pool store was changed after failed initialization")
	}

	if stagedDB == nil || !poolClosed(stagedDB) {
		t.Fatal("staged pool was not closed after failed initialization")
	}
}

func TestInitializeActiveConnectionsClearsPoolsWithoutActiveProject(t *testing.T) {
	store := NewPoolStore()
	db := newTestPool(t)
	if err := store.Set("project-1", "node-1", db); err != nil {
		t.Fatal(err)
	}

	manager := NewConnectionManager(
		store,
		&projectRepoStub{runningErr: sql.ErrNoRows},
		&nodeRepoStub{},
		&connectionRepoStub{},
	)

	if err := manager.InitializeActiveConnections(context.Background()); err != nil {
		t.Fatal(err)
	}

	if !poolClosed(db) {
		t.Fatal("inactive project pool was not closed")
	}

	if _, err := store.Get("project-1", "node-1"); !errors.Is(err, ErrProjectPoolsNotFound) {
		t.Fatalf("expected project pools not found, got %v", err)
	}
}

func TestSyncNodeRemovesStalePoolAfterConnectionFailure(t *testing.T) {
	store := NewPoolStore()
	oldDB := newTestPool(t)
	if err := store.Set("project-1", "node-1", oldDB); err != nil {
		t.Fatal(err)
	}

	manager := NewConnectionManager(
		store,
		&projectRepoStub{project: repository.Project{ID: "project-1", Running: true}},
		&nodeRepoStub{node: repository.Node{ID: "node-1", ProjectID: "project-1", Status: true}},
		&connectionRepoStub{connection: repository.NodeConnection{NodeId: "node-1"}},
	)
	manager.connect = func(ctx context.Context, dsn string) (*sql.DB, error) {
		return nil, errors.New("connection failed")
	}

	if err := manager.SyncNode(context.Background(), "project-1", "node-1"); err == nil {
		t.Fatal("expected synchronization error")
	}

	if !poolClosed(oldDB) {
		t.Fatal("stale pool was not closed")
	}

	if _, err := store.Get("project-1", "node-1"); !errors.Is(err, ErrProjectPoolsNotFound) {
		t.Fatalf("expected project pools not found, got %v", err)
	}
}

func TestCheckConnectionHealthRemovesFailedPool(t *testing.T) {
	store := NewPoolStore()
	db := newTestPool(t)
	if err := store.Set("project-1", "node-1", db); err != nil {
		t.Fatal(err)
	}

	manager := NewConnectionManager(
		store,
		&projectRepoStub{},
		&nodeRepoStub{},
		&connectionRepoStub{},
	)

	healthy, err := manager.CheckConnectionHealth(
		context.Background(),
		"project-1",
		"node-1",
	)
	if err == nil || healthy {
		t.Fatal("expected failed health check")
	}

	if _, err := store.Get("project-1", "node-1"); !errors.Is(err, ErrProjectPoolsNotFound) {
		t.Fatalf("expected project pools not found, got %v", err)
	}
}
