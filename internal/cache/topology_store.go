package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type TopologyStore struct {
	mu   sync.RWMutex
	data map[string][]repository.NodeTopology
}

func NewTopologyStore() *TopologyStore {
	return &TopologyStore{
		data: make(map[string][]repository.NodeTopology),
	}
}

func (s *TopologyStore) Set(projectID string, topology []repository.NodeTopology) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[projectID] = topology
}

func (s *TopologyStore) Get(projectID string) ([]repository.NodeTopology, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topology, ok := s.data[projectID]
	return topology, ok
}

func (s *TopologyStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}
