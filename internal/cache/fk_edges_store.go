package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type FKEdgesStore struct {
	mu sync.RWMutex

	// projectID -> edgeKey -> FkEdges
	data map[string]map[string]repository.FkEdges
}

func NewFKEdgesStore() *FKEdgesStore {
	return &FKEdgesStore{
		data: make(map[string]map[string]repository.FkEdges),
	}
}

func edgeKey(
	parentTable,
	parentColumn,
	childTable,
	childColumn string,
) string {
	return parentTable +
		":" + parentColumn +
		":" + childTable +
		":" + childColumn
}

func (s *FKEdgesStore) Set(edge repository.FkEdges) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[edge.ProjectId]; !ok {
		s.data[edge.ProjectId] = make(map[string]repository.FkEdges)
	}

	s.data[edge.ProjectId][edgeKey(
		edge.ParentTable,
		edge.ParentColumn,
		edge.ChildTable,
		edge.ChildColumn,
	)] = edge
}

func (s *FKEdgesStore) Get(
	projectID,
	parentTable,
	parentColumn,
	childTable,
	childColumn string,
) (repository.FkEdges, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectEdges, ok := s.data[projectID]
	if !ok {
		return repository.FkEdges{}, false
	}

	edge, ok := projectEdges[edgeKey(
		parentTable,
		parentColumn,
		childTable,
		childColumn,
	)]

	return edge, ok
}

func (s *FKEdgesStore) GetByProject(
	projectID string,
) []repository.FkEdges {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectEdges, ok := s.data[projectID]
	if !ok {
		return nil
	}

	edges := make([]repository.FkEdges, 0, len(projectEdges))

	for _, edge := range projectEdges {
		edges = append(edges, edge)
	}

	return edges
}

func (s *FKEdgesStore) GetByTable(
	projectID,
	tableName string,
) []repository.FkEdges {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectEdges, ok := s.data[projectID]
	if !ok {
		return nil
	}

	edges := make([]repository.FkEdges, 0)

	for _, edge := range projectEdges {
		if edge.ParentTable == tableName ||
			edge.ChildTable == tableName {

			edges = append(edges, edge)
		}
	}

	return edges
}

func (s *FKEdgesStore) Delete(edge repository.FkEdges) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectEdges, ok := s.data[edge.ProjectId]
	if !ok {
		return
	}

	delete(projectEdges, edgeKey(
		edge.ParentTable,
		edge.ParentColumn,
		edge.ChildTable,
		edge.ChildColumn,
	))

	if len(projectEdges) == 0 {
		delete(s.data, edge.ProjectId)
	}
}

func (s *FKEdgesStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *FKEdgesStore) Exists(edge repository.FkEdges) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectEdges, ok := s.data[edge.ProjectId]
	if !ok {
		return false
	}

	_, ok = projectEdges[edgeKey(
		edge.ParentTable,
		edge.ParentColumn,
		edge.ChildTable,
		edge.ChildColumn,
	)]

	return ok
}

func (s *FKEdgesStore) GetAll() []repository.FkEdges {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var edges []repository.FkEdges

	for _, projectEdges := range s.data {
		for _, edge := range projectEdges {
			edges = append(edges, edge)
		}
	}

	return edges
}

func (s *FKEdgesStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]map[string]repository.FkEdges)
}
