package cache

import (
	"hyperfulcrum/internal/repository"
	"sync"
)

type NodeStore struct {
	mu sync.RWMutex

	// projectID -> nodeID -> Node
	data map[string]map[string]repository.Node

	// nodeID -> projectID
	projectIndex map[string]string

	loadedProjects map[string]bool
}

func NewNodeStore() *NodeStore {
	return &NodeStore{
		data:           make(map[string]map[string]repository.Node),
		projectIndex:   make(map[string]string),
		loadedProjects: make(map[string]bool),
	}
}

func (s *NodeStore) Set(node repository.Node) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[node.ProjectID]; !ok {
		s.data[node.ProjectID] = make(map[string]repository.Node)
	}

	s.data[node.ProjectID][node.ID] = node
	s.projectIndex[node.ID] = node.ProjectID
}

func (s *NodeStore) ReplaceProject(
	projectID string,
	nodes []repository.Node,
) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if projectNodes, ok := s.data[projectID]; ok {
		for nodeID := range projectNodes {
			delete(s.projectIndex, nodeID)
		}
	}

	projectNodes := make(map[string]repository.Node, len(nodes))

	for _, node := range nodes {
		node.ProjectID = projectID
		projectNodes[node.ID] = node
		s.projectIndex[node.ID] = projectID
	}

	s.data[projectID] = projectNodes
	s.loadedProjects[projectID] = true
}

func (s *NodeStore) Get(nodeID string) (repository.Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectID, ok := s.projectIndex[nodeID]
	if !ok {
		return repository.Node{}, false
	}

	node, ok := s.data[projectID][nodeID]
	return node, ok
}

func (s *NodeStore) GetProjectID(nodeID string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectID, ok := s.projectIndex[nodeID]
	return projectID, ok
}

func (s *NodeStore) GetByProject(projectID string) ([]repository.Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectNodes, ok := s.data[projectID]
	if !ok {
		return nil, s.loadedProjects[projectID]
	}

	nodes := make([]repository.Node, 0, len(projectNodes))

	for _, node := range projectNodes {
		nodes = append(nodes, node)
	}

	return nodes, s.loadedProjects[projectID]
}

func (s *NodeStore) GetAll() []repository.Node {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nodes := make([]repository.Node, 0)

	for _, projectNodes := range s.data {
		for _, node := range projectNodes {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

func (s *NodeStore) Delete(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectID, ok := s.projectIndex[nodeID]
	if !ok {
		return
	}

	delete(s.data[projectID], nodeID)

	if len(s.data[projectID]) == 0 {
		delete(s.data, projectID)
	}

	delete(s.projectIndex, nodeID)
}

func (s *NodeStore) DeleteProject(projectID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	projectNodes, ok := s.data[projectID]
	if ok {
		for nodeID := range projectNodes {
			delete(s.projectIndex, nodeID)
		}
	}

	delete(s.data, projectID)
	delete(s.loadedProjects, projectID)
}

func (s *NodeStore) Exists(nodeID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.projectIndex[nodeID]
	return ok
}

func (s *NodeStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]map[string]repository.Node)
	s.projectIndex = make(map[string]string)
	s.loadedProjects = make(map[string]bool)
}
