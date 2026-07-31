package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type SchemaVersionStore struct {
	mu sync.RWMutex

	// projectID -> SchemaVersion
	data map[string]repository.SchemaVersion
}

func NewSchemaVersionStore() *SchemaVersionStore {
	return &SchemaVersionStore{
		data: make(map[string]repository.SchemaVersion),
	}
}

func (s *SchemaVersionStore) Set(
	schema repository.SchemaVersion,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[schema.ProjectID] = schema
}

func (s *SchemaVersionStore) Get(
	projectID string,
) (repository.SchemaVersion, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	schema, ok := s.data[projectID]

	return schema, ok
}

func (s *SchemaVersionStore) GetAll() []repository.SchemaVersion {

	s.mu.RLock()
	defer s.mu.RUnlock()

	schemas := make(
		[]repository.SchemaVersion,
		0,
		len(s.data),
	)

	for _, schema := range s.data {
		schemas = append(schemas, schema)
	}

	return schemas
}

func (s *SchemaVersionStore) DeleteProject(
	projectID string,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *SchemaVersionStore) Exists(
	projectID string,
) bool {

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[projectID]

	return ok
}

func (s *SchemaVersionStore) Clear() {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.SchemaVersion)
}
