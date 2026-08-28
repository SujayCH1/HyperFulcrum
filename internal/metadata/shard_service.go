package metadata

import (
	"context"
	"database/sql"
	"errors"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type ShardService struct {
	repo      *repository.ShardRepository
	nodes     *repository.NodeRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewShardService(repo *repository.ShardRepository, nodes *repository.NodeRepository,
	cacheManager *cache.CacheManager, refresher *cache.CacheRefresher) *ShardService {
	return &ShardService{repo: repo, nodes: nodes, cache: cacheManager, refresher: refresher}
}

func (s *ShardService) AddShard(ctx context.Context, projectID, name, primaryNodeID string) (repository.Shard, error) {
	if err := s.validateAdd(ctx, projectID, name, primaryNodeID); err != nil {
		return repository.Shard{}, err
	}
	shard, err := s.repo.ShardAdd(ctx, projectID, name, primaryNodeID)
	if err != nil {
		return repository.Shard{}, err
	}
	if err := s.refresher.RefreshShards(ctx, projectID); err != nil {
		return repository.Shard{}, err
	}
	return shard, nil
}

func (s *ShardService) ListShards(ctx context.Context, projectID string) ([]repository.Shard, error) {
	shards, loaded := s.cache.Shards.GetByProject(projectID)
	if loaded {
		return shards, nil
	}
	if err := s.refresher.RefreshShards(ctx, projectID); err != nil {
		return nil, err
	}
	shards, _ = s.cache.Shards.GetByProject(projectID)
	return shards, nil
}

func (s *ShardService) GetShard(ctx context.Context, shardID string) (repository.Shard, error) {
	if shard, ok := s.cache.Shards.GetByID(shardID); ok {
		return shard, nil
	}
	shard, err := s.repo.ShardGetByID(ctx, shardID)
	if err != nil {
		return repository.Shard{}, err
	}
	if err := s.refresher.RefreshShards(ctx, shard.ProjectID); err != nil {
		return repository.Shard{}, err
	}
	return shard, nil
}

func (s *ShardService) RenameShard(ctx context.Context, shardID, name string) error {
	shard, err := s.GetShard(ctx, shardID)
	if err != nil {
		return err
	}
	if err := s.ensureProjectStopped(ctx, shard.ProjectID); err != nil {
		return err
	}
	shards, err := s.ListShards(ctx, shard.ProjectID)
	if err != nil {
		return err
	}
	for _, existing := range shards {
		if existing.ID != shardID && existing.Name == name {
			return ErrDuplicateShardName
		}
	}
	if err := s.repo.ShardUpdateName(ctx, shardID, name); err != nil {
		return err
	}
	return s.refresher.RefreshShards(ctx, shard.ProjectID)
}

func (s *ShardService) RemoveShard(ctx context.Context, shardID string) error {
	shard, err := s.GetShard(ctx, shardID)
	if err != nil {
		return err
	}
	if err := s.ensureProjectStopped(ctx, shard.ProjectID); err != nil {
		return err
	}
	topologies, err := s.topologies(ctx, shard.ProjectID)
	if err != nil {
		return err
	}
	for _, topology := range topologies {
		if topology.ShardID == shardID {
			return ErrShardHasStandbys
		}
	}
	if err := s.repo.ShardRemove(ctx, shardID); err != nil {
		return err
	}
	return errors.Join(s.refresher.RefreshShards(ctx, shard.ProjectID),
		s.refresher.RefreshNodes(ctx, shard.ProjectID))
}

func (s *ShardService) validateAdd(ctx context.Context, projectID, name, primaryNodeID string) error {
	if err := s.ensureProjectStopped(ctx, projectID); err != nil {
		return err
	}
	node, err := s.nodes.NodeGetByID(ctx, primaryNodeID)
	if err != nil {
		return err
	}
	if node.ProjectID != projectID {
		return sql.ErrNoRows
	}
	if node.Role != repository.NodeRolePrimary {
		return ErrPrimaryRoleRequired
	}
	shards, err := s.ListShards(ctx, projectID)
	if err != nil {
		return err
	}
	for _, shard := range shards {
		if shard.Name == name {
			return ErrDuplicateShardName
		}
		if shard.PrimaryNodeID == primaryNodeID {
			return ErrPrimaryNodeAlreadyUsed
		}
	}
	return nil
}

func (s *ShardService) ensureProjectStopped(ctx context.Context, projectID string) error {
	project, ok := s.cache.Projects.Get(projectID)
	if !ok {
		if err := s.refresher.RefreshProject(ctx, projectID); err != nil {
			return err
		}
		project, ok = s.cache.Projects.Get(projectID)
	}
	if !ok {
		return sql.ErrNoRows
	}
	if project.Running {
		return ErrProjectRunning
	}
	return nil
}

func (s *ShardService) topologies(ctx context.Context, projectID string) ([]repository.NodeTopology, error) {
	values, loaded := s.cache.Topology.GetByProjectID(projectID)
	if !loaded {
		if err := s.refresher.RefreshTopology(ctx, projectID); err != nil {
			return nil, err
		}
		values, _ = s.cache.Topology.GetByProjectID(projectID)
	}
	return values, nil
}
