package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type ProjectStore struct {
	mu   sync.Mutex
	data map[string]repository.Project
}

func NewProjectStore() *ProjectStore {
	return &ProjectStore{
		data: make(map[string]repository.Project),
	}
}

func (s *ProjectStore) Set(project repository.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[project.ID] = project
}

func (s *ProjectStore) Get(projectID string) (repository.Project, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	project, ok := s.data[projectID]
	return project, ok
}

func (s *ProjectStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *ProjectStore) Exists(projectID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[projectID]
	return ok
}
