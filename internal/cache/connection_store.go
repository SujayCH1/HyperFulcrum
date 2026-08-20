package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ConnectionStore struct {
	mu sync.RWMutex

	// projectID -> nodeID -> Connection
	data map[string]map[string]repository.NodeConnection

	// nodeID -> projectID
	projectIndex map[string]string

	loadedProjects map[string]bool
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		data:           make(map[string]map[string]repository.NodeConnection),
		projectIndex:   make(map[string]string),
		loadedProjects: make(map[string]bool),
	}
}

func (s *ConnectionStore) Set(
	projectID string,
	connection repository.NodeConnection,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[projectID]; !ok {
		s.data[projectID] = make(map[string]repository.NodeConnection)
	}

	s.data[projectID][connection.NodeId] = connection
	s.projectIndex[connection.NodeId] = projectID
}

func (s *ConnectionStore) ReplaceProject(
	projectID string,
	connections []repository.NodeConnection,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if projectConnections, ok := s.data[projectID]; ok {
		for nodeID := range projectConnections {
			delete(s.projectIndex, nodeID)
		}
	}

	projectConnections := make(
		map[string]repository.NodeConnection,
		len(connections),
	)

	for _, connection := range connections {
		projectConnections[connection.NodeId] = connection
		s.projectIndex[connection.NodeId] = projectID
	}

	s.data[projectID] = projectConnections
	s.loadedProjects[projectID] = true
}

func (s *ConnectionStore) Get(nodeID string) (repository.NodeConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectID, ok := s.projectIndex[nodeID]
	if !ok {
		return repository.NodeConnection{}, false
	}

	conn, ok := s.data[projectID][nodeID]
	return conn, ok
}

func (s *ConnectionStore) GetByProject(projectID string) ([]repository.NodeConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectConnections, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}

	connections := make([]repository.NodeConnection, 0, len(projectConnections))

	for _, conn := range projectConnections {
		connections = append(connections, conn)
	}

	return connections, s.loadedProjects[projectID]
}

func (s *ConnectionStore) GetAll() []repository.NodeConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := make([]repository.NodeConnection, 0)

	for _, projectConnections := range s.data {
		for _, conn := range projectConnections {
			connections = append(connections, conn)
		}
	}

	return connections
}

func (s *ConnectionStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectID, ok := s.projectIndex[nodeID]
	if !ok {
		return
	}

	delete(s.data[projectID], nodeID)

	if len(s.data[projectID]) == 0 {
		delete(s.data, projectID)
	}

	delete(s.projectIndex, nodeID)
}

func (s *ConnectionStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectConnections, ok := s.data[projectID]
	if ok {
		for nodeID := range projectConnections {
			delete(s.projectIndex, nodeID)
		}
	}

	delete(s.data, projectID)
	delete(s.loadedProjects, projectID)
}

func (s *ConnectionStore) Exists(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.projectIndex[nodeID]
	return ok
}

func (s *ConnectionStore) ProjectLoaded(projectID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loadedProjects[projectID]
}

func (s *ConnectionStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]map[string]repository.NodeConnection)
	s.projectIndex = make(map[string]string)
	s.loadedProjects = make(map[string]bool)
}
