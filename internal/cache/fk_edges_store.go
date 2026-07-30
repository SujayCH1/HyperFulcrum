package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type FKEdgesStore struct {
	mu   sync.RWMutex
	data map[string]repository.FkEdges
}

func NewFKEdgesStore() *FKEdgesStore {
	return &FKEdgesStore{
		data: make(map[string]repository.FkEdges),
	}
}

func edgeKey(e repository.FkEdges) string {
	return e.ProjectId +
		":" + e.ParentTable +
		":" + e.ParentColumn +
		":" + e.ChildTable +
		":" + e.ChildColumn
}

func (s *FKEdgesStore) Set(edge repository.FkEdges) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[edgeKey(edge)] = edge
}

func (s *FKEdgesStore) GetAll() []repository.FkEdges {
	s.mu.RLock()
	defer s.mu.RUnlock()

	edges := make([]repository.FkEdges, 0, len(s.data))

	for _, e := range s.data {
		edges = append(edges, e)
	}

	return edges
}

func (s *FKEdgesStore) Delete(edge repository.FkEdges) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, edgeKey(edge))
}

func (s *FKEdgesStore) Exists(edge repository.FkEdges) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[edgeKey(edge)]
	return ok
}

func (s *FKEdgesStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.FkEdges)
}
