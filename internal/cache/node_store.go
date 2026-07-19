package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type NodeStore struct {
	mu   sync.RWMutex
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
	s.mu.RLock()
	defer s.mu.RUnlock()

	node, ok := s.data[nodeID]
	return node, ok
}

func (s *NodeStore) GetAll() []repository.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]repository.Node, 0, len(s.data))

	for _, node := range s.data {
		nodes = append(nodes, node)
	}

	return nodes
}

func (s *NodeStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, nodeID)
}

func (s *NodeStore) Exists(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[nodeID]
	return ok
}

func (s *NodeStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.Node)
}
