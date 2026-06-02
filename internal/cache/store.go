package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ProjectCacheStore struct {
	mu sync.RWMutex
	// key is projectID
	projects map[string]repository.Project
}

type NodeCacheStore struct {
	mu sync.RWMutex
	// key is projectID
	nodes map[string][]repository.Node
}

type ConnectionsCacheStore struct {
	mu sync.RWMutex
	// key is nodeID
	connections map[string][]repository.NodeConnection
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
		connections: make(map[string][]repository.NodeConnection),
	}
}

// functions for Project cache
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

// functions for Node cache
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

// functions for Connections cache
func (s *ConnectionsCacheStore) Set(projectID string, connections []repository.NodeConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.connections[projectID] = connections
}

func (s *ConnectionsCacheStore) Get(projectID string) ([]repository.NodeConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections, exists := s.connections[projectID]
	return connections, exists
}

func (s *ConnectionsCacheStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.connections, projectID)
}
