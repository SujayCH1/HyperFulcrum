package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type ShardKeysStore struct {
	mu sync.RWMutex

	// projectID -> tableName -> shard key
	data           map[string]map[string]repository.ShardKey
	loadedProjects map[string]bool
}

func NewShardKeysStore() *ShardKeysStore {
	return &ShardKeysStore{
		data:           make(map[string]map[string]repository.ShardKey),
		loadedProjects: make(map[string]bool),
	}
}

func (s *ShardKeysStore) Set(key repository.ShardKey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[key.ProjectID]
	if !ok {
		s.data[key.ProjectID] = make(map[string]repository.ShardKey)
	}
	s.data[key.ProjectID][key.TableName] = key
}

func (s *ShardKeysStore) ReplaceProject(projectID string, keys []repository.ShardKey) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectKeys := make(map[string]repository.ShardKey, len(keys))
	for _, key := range keys {
		key.ProjectID = projectID
		projectKeys[key.TableName] = key
	}

	s.data[projectID] = projectKeys
	s.loadedProjects[projectID] = true
}

func (s *ShardKeysStore) Get(projectID, tableName string) (repository.ShardKey, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectKeys, ok := s.data[projectID]
	if !ok {
		return repository.ShardKey{}, false
	}

	key, ok := projectKeys[tableName]
	return key, ok
}

func (s *ShardKeysStore) GetByProject(projectID string) ([]repository.ShardKey, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectKeys, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}

	keys := make([]repository.ShardKey, 0, len(projectKeys))
	for _, key := range projectKeys {
		keys = append(keys, key)
	}

	return keys, s.loadedProjects[projectID]
}

func (s *ShardKeysStore) Delete(projectID, tableName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectKeys, ok := s.data[projectID]
	if !ok {
		return
	}

	delete(projectKeys, tableName)
	if len(projectKeys) == 0 {
		delete(s.data, projectID)
	}
}

func (s *ShardKeysStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, projectID)
	delete(s.loadedProjects, projectID)
}

func (s *ShardKeysStore) GetAll() []repository.ShardKey {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []repository.ShardKey
	for _, projectKeys := range s.data {
		for _, key := range projectKeys {
			keys = append(keys, key)
		}
	}
	return keys
}

func (s *ShardKeysStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]map[string]repository.ShardKey)
	s.loadedProjects = make(map[string]bool)
}
