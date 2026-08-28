package cache

import (
	"sync"

	"hyperfulcrum/internal/repository"
)

type NodeRuntimeStateStore struct {
	mu       sync.RWMutex
	data     map[string]repository.NodeRuntimeState
	projects map[string]map[string]struct{}
}

func NewNodeRuntimeStateStore() *NodeRuntimeStateStore {
	return &NodeRuntimeStateStore{data: make(map[string]repository.NodeRuntimeState),
		projects: make(map[string]map[string]struct{})}
}

func (s *NodeRuntimeStateStore) ReplaceProject(projectID string, states []repository.NodeRuntimeState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for nodeID := range s.projects[projectID] {
		delete(s.data, nodeID)
	}
	ids := make(map[string]struct{}, len(states))
	for _, state := range states {
		s.data[state.NodeID] = state
		ids[state.NodeID] = struct{}{}
	}
	s.projects[projectID] = ids
}

func (s *NodeRuntimeStateStore) GetByNodeID(nodeID string) (repository.NodeRuntimeState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.data[nodeID]
	return state, ok
}

func (s *NodeRuntimeStateStore) GetByProject(projectID string) []repository.NodeRuntimeState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	states := make([]repository.NodeRuntimeState, 0, len(s.projects[projectID]))
	for nodeID := range s.projects[projectID] {
		states = append(states, s.data[nodeID])
	}
	return states
}

func (s *NodeRuntimeStateStore) DeleteNode(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, nodeID)
	for _, ids := range s.projects {
		delete(ids, nodeID)
	}
}

func (s *NodeRuntimeStateStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for nodeID := range s.projects[projectID] {
		delete(s.data, nodeID)
	}
	delete(s.projects, projectID)
}

func (s *NodeRuntimeStateStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]repository.NodeRuntimeState)
	s.projects = make(map[string]map[string]struct{})
}
