package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ColumnStore struct {
	mu sync.RWMutex

	// projectID -> columnKey -> Column
	data map[string]map[columnCacheKey]repository.Column

	loadedProjects map[string]bool
}

type columnCacheKey struct {
	tableName  string
	columnName string
}

func NewColumnStore() *ColumnStore {
	return &ColumnStore{
		data:           make(map[string]map[columnCacheKey]repository.Column),
		loadedProjects: make(map[string]bool),
	}
}

func columnKey(tableName, columnName string) columnCacheKey {
	return columnCacheKey{
		tableName:  tableName,
		columnName: columnName,
	}
}

func (s *ColumnStore) Set(col repository.Column) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[col.ProjectID]; !ok {
		s.data[col.ProjectID] = make(map[columnCacheKey]repository.Column)
	}

	s.data[col.ProjectID][columnKey(col.TableName, col.ColumnName)] = col
}

func (s *ColumnStore) ReplaceProject(
	projectID string,
	columns []repository.Column,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectColumns := make(
		map[columnCacheKey]repository.Column,
		len(columns),
	)

	for _, col := range columns {
		col.ProjectID = projectID
		projectColumns[columnKey(col.TableName, col.ColumnName)] = col
	}

	s.data[projectID] = projectColumns
	s.loadedProjects[projectID] = true
}

func (s *ColumnStore) Get(
	projectID,
	tableName,
	columnName string,
) (repository.Column, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectColumns, ok := s.data[projectID]
	if !ok {
		return repository.Column{}, false
	}

	col, ok := projectColumns[columnKey(tableName, columnName)]
	return col, ok
}

func (s *ColumnStore) GetByProject(
	projectID string,
) ([]repository.Column, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectColumns, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}

	columns := make([]repository.Column, 0, len(projectColumns))

	for _, col := range projectColumns {
		columns = append(columns, col)
	}

	return columns, s.loadedProjects[projectID]
}

func (s *ColumnStore) GetByTable(
	projectID,
	tableName string,
) ([]repository.Column, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	projectColumns, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}

	columns := make([]repository.Column, 0)

	for _, col := range projectColumns {
		if col.TableName == tableName {
			columns = append(columns, col)
		}
	}

	return columns, s.loadedProjects[projectID]
}

func (s *ColumnStore) Delete(col repository.Column) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectColumns, ok := s.data[col.ProjectID]
	if !ok {
		return
	}

	delete(projectColumns, columnKey(col.TableName, col.ColumnName))

	if len(projectColumns) == 0 {
		delete(s.data, col.ProjectID)
	}
}

func (s *ColumnStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
	delete(s.loadedProjects, projectID)
}

func (s *ColumnStore) Exists(col repository.Column) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectColumns, ok := s.data[col.ProjectID]
	if !ok {
		return false
	}

	_, ok = projectColumns[columnKey(col.TableName, col.ColumnName)]
	return ok
}

func (s *ColumnStore) GetAll() []repository.Column {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var columns []repository.Column

	for _, projectColumns := range s.data {
		for _, col := range projectColumns {
			columns = append(columns, col)
		}
	}

	return columns
}

func (s *ColumnStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]map[columnCacheKey]repository.Column)
	s.loadedProjects = make(map[string]bool)
}
