package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ProjectCacheStore struct {
	mu       sync.RWMutex
	projects map[string]repository.Project
}

type NodeCacheStore struct {
	mu    sync.RWMutex
	nodes map[string][]repository.Node
}

type ConnectionsCacheStore struct {
	mu          sync.RWMutex
	connections map[string]repository.NodeConnection
}

func NewProjectCacheStore() *ProjectCacheStore {
	return &ProjectCacheStore{
		projects: make(map[string]repository.Project),
	}
}

func NewNodeCacheStore() *NodeCacheStore {
	return &NodeCacheStore{
		nodes: make(map[string][]repository.Node),
	}
}

func NewConnectionCacheStore() *ConnectionsCacheStore {
	return &ConnectionsCacheStore{
		connections: make(map[string]repository.NodeConnection),
	}
}

// Project cache
func (s *ProjectCacheStore) Set(projectID string, project repository.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.projects[projectID] = project
}

func (s *ProjectCacheStore) Get(projectID string) (repository.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, exists := s.projects[projectID]
	return project, exists
}

func (s *ProjectCacheStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.projects, projectID)
}

// Node cache
func (s *NodeCacheStore) Set(projectID string, nodes []repository.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodes[projectID] = nodes
}

func (s *NodeCacheStore) Get(projectID string) ([]repository.Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes, exists := s.nodes[projectID]
	return nodes, exists
}

func (s *NodeCacheStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.nodes, projectID)
}

// Connection cache
func (s *ConnectionsCacheStore) Set(nodeID string, connection repository.NodeConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[nodeID] = connection
}

func (s *ConnectionsCacheStore) Get(nodeID string) (repository.NodeConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connection, exists := s.connections[nodeID]
	return connection, exists
}

func (s *ConnectionsCacheStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.connections, nodeID)
}
