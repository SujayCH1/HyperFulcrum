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

func (s *ConnectionStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, nodeID)
}
