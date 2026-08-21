package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type ProjectStore struct {
	mu     sync.RWMutex
	data   map[string]repository.Project
	loaded bool
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

func (s *ProjectStore) Replace(projects []repository.Project) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := make(map[string]repository.Project, len(projects))

	for _, project := range projects {
		data[project.ID] = project
	}

	s.data = data
	s.loaded = true
}

func (s *ProjectStore) Get(projectID string) (repository.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	project, ok := s.data[projectID]
	return project, ok
}

func (s *ProjectStore) GetAll() ([]repository.Project, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projects := make([]repository.Project, 0, len(s.data))

	for _, project := range s.data {
		projects = append(projects, project)
	}

	return projects, s.loaded
}

func (s *ProjectStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *ProjectStore) Exists(projectID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[projectID]
	return ok
}

func (s *ProjectStore) Loaded() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.loaded
}

func (s *ProjectStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]repository.Project)
	s.loaded = false
}
