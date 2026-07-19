package cache

import "sync"

type ProjectNodeStore struct {
	mu   sync.RWMutex
	data map[string][]string
}

func NewProjectNodeStore() *ProjectNodeStore {
	return &ProjectNodeStore{
		data: make(map[string][]string),
	}
}

func (s *ProjectNodeStore) Set(projectID string, nodeIDs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[projectID] = nodeIDs
}

func (s *ProjectNodeStore) Get(projectID string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes, ok := s.data[projectID]
	return nodes, ok
}

func (s *ProjectNodeStore) Delete(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
}

func (s *ProjectNodeStore) Add(projectID, nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[projectID] = append(s.data[projectID], nodeID)
}
