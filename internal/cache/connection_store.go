package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ConnectionStore struct {
	mu   sync.RWMutex
	data map[string]repository.NodeConnection
}

func NewConnectionStore() *ConnectionStore {
	return &ConnectionStore{
		data: make(map[string]repository.NodeConnection),
	}
}

func (s *ConnectionStore) Set(connection repository.NodeConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[connection.NodeId] = connection
}

func (s *ConnectionStore) Get(nodeID string) (repository.NodeConnection, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	conn, ok := s.data[nodeID]
	return conn, ok
}

func (s *ConnectionStore) GetAll() []repository.NodeConnection {
	s.mu.RLock()
	defer s.mu.RUnlock()

	connections := make([]repository.NodeConnection, 0, len(s.data))

	for _, connection := range s.data {
		connections = append(connections, connection)
	}

	return connections
}

func (s *ConnectionStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, nodeID)
}

func (s *ConnectionStore) Exists(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[nodeID]
	return ok
}

func (s *ConnectionStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.NodeConnection)
}
