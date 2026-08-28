package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type NodeRuntimeStateService struct {
	repo      *repository.NodeRuntimeStateRepository
	nodes     *repository.NodeRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewNodeRuntimeStateService(repo *repository.NodeRuntimeStateRepository,
	nodes *repository.NodeRepository, cacheManager *cache.CacheManager,
	refresher *cache.CacheRefresher) *NodeRuntimeStateService {
	return &NodeRuntimeStateService{repo: repo, nodes: nodes, cache: cacheManager, refresher: refresher}
}

func (s *NodeRuntimeStateService) GetByNodeID(ctx context.Context, nodeID string) (repository.NodeRuntimeState, error) {
	if state, ok := s.cache.Runtime.GetByNodeID(nodeID); ok {
		return state, nil
	}
	node, err := s.nodes.NodeGetByID(ctx, nodeID)
	if err != nil {
		return repository.NodeRuntimeState{}, err
	}
	if err := s.refresher.RefreshNodeRuntimeStates(ctx, node.ProjectID); err != nil {
		return repository.NodeRuntimeState{}, err
	}
	state, ok := s.cache.Runtime.GetByNodeID(nodeID)
	if !ok {
		return repository.NodeRuntimeState{}, sql.ErrNoRows
	}
	return state, nil
}

func (s *NodeRuntimeStateService) ListByProject(ctx context.Context, projectID string) ([]repository.NodeRuntimeState, error) {
	if _, ok := s.cache.Projects.Get(projectID); !ok {
		if err := s.refresher.RefreshProject(ctx, projectID); err != nil {
			return nil, err
		}
		if _, ok := s.cache.Projects.Get(projectID); !ok {
			return nil, sql.ErrNoRows
		}
	}
	if err := s.refresher.RefreshNodeRuntimeStates(ctx, projectID); err != nil {
		return nil, err
	}
	return s.cache.Runtime.GetByProject(projectID), nil
}
