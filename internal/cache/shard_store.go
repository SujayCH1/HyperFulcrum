package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type ShardStore struct {
	mu             sync.RWMutex
	data           map[string]map[string]repository.Shard
	projectIndex   map[string]string
	loadedProjects map[string]bool
}

func NewShardStore() *ShardStore {
	return &ShardStore{data: make(map[string]map[string]repository.Shard),
		projectIndex: make(map[string]string), loadedProjects: make(map[string]bool)}
}

func (s *ShardStore) ReplaceProject(projectID string, shards []repository.Shard) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if current := s.data[projectID]; current != nil {
		for id := range current {
			delete(s.projectIndex, id)
		}
	}
	values := make(map[string]repository.Shard, len(shards))
	for _, shard := range shards {
		shard.ProjectID = projectID
		values[shard.ID] = shard
		s.projectIndex[shard.ID] = projectID
	}
	s.data[projectID] = values
	s.loadedProjects[projectID] = true
}

func (s *ShardStore) GetByProject(projectID string) ([]repository.Shard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}
	shards := make([]repository.Shard, 0, len(values))
	for _, shard := range values {
		shards = append(shards, shard)
	}
	return shards, s.loadedProjects[projectID]
}

func (s *ShardStore) GetByID(shardID string) (repository.Shard, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	projectID, ok := s.projectIndex[shardID]
	if !ok {
		return repository.Shard{}, false
	}
	shard, ok := s.data[projectID][shardID]
	return shard, ok
}

func (s *ShardStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id := range s.data[projectID] {
		delete(s.projectIndex, id)
	}
	delete(s.data, projectID)
	delete(s.loadedProjects, projectID)
}

func (s *ShardStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]map[string]repository.Shard)
	s.projectIndex = make(map[string]string)
	s.loadedProjects = make(map[string]bool)
}
