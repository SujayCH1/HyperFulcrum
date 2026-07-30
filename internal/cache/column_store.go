package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ColumnStore struct {
	mu   sync.RWMutex
	data map[string]repository.Column
}

func NewColumnStore() *ColumnStore {
	return &ColumnStore{
		data: make(map[string]repository.Column),
	}
}

func columnKey(c repository.Column) string {
	return c.ProjectID + ":" + c.TableName + ":" + c.ColumnName
}

func (s *ColumnStore) Set(col repository.Column) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[columnKey(col)] = col
}

func (s *ColumnStore) GetAll() []repository.Column {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cols := make([]repository.Column, 0, len(s.data))

	for _, c := range s.data {
		cols = append(cols, c)
	}

	return cols
}

func (s *ColumnStore) Delete(col repository.Column) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, columnKey(col))
}

func (s *ColumnStore) Exists(col repository.Column) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[columnKey(col)]
	return ok
}

func (s *ColumnStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.Column)
}
