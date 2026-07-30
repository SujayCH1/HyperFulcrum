package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type NodeTopologyStore struct {
	mu   sync.RWMutex
	data map[string][]repository.NodeTopology
}

func NewTopologyStore() *NodeTopologyStore {
	return &NodeTopologyStore{
		data: make(map[string][]repository.NodeTopology),
	}
}

func (s *NodeTopologyStore) Set(
	projectID string,
	topologies []repository.NodeTopology,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copied := make([]repository.NodeTopology, len(topologies))
	copy(copied, topologies)

	s.data[projectID] = copied
}

func (s *NodeTopologyStore) GetByProjectID(
	projectID string,
) ([]repository.NodeTopology, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	topologies, ok := s.data[projectID]
	if !ok {
		return nil, false
	}

	copied := make([]repository.NodeTopology, len(topologies))
	copy(copied, topologies)

	return copied, true
}

func (s *NodeTopologyStore) GetAll() []repository.NodeTopology {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]repository.NodeTopology, 0)

	for _, topologies := range s.data {
		result = append(result, topologies...)
	}

	return result
}

func (s *NodeTopologyStore) DeleteByProjectID(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *NodeTopologyStore) ExistsByProjectID(projectID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[projectID]
	return ok
}

func (s *NodeTopologyStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string][]repository.NodeTopology)
}
