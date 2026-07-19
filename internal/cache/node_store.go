package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type NodeStore struct {
	mu   sync.Mutex
	data map[string]repository.Node
}

func NewNodeStore() *NodeStore {
	return &NodeStore{
		data: make(map[string]repository.Node),
	}
}

func (s *NodeStore) Set(node repository.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[node.ID] = node
}

func (s *NodeStore) Get(nodeID string) (repository.Node, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	node, ok := s.data[nodeID]
	return node, ok
}

func (s *NodeStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, nodeID)

}

func (s *NodeStore) Exists(nodeID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[nodeID]
	return ok
}
