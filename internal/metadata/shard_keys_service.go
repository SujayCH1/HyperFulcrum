package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type ShardKeysService struct {
	repo      *repository.ShardKeyRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewShardKeysService(
	repo *repository.ShardKeyRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *ShardKeysService {

	return &ShardKeysService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *ShardKeysService) AddShardKey(
	ctx context.Context,
	projectID string,
	tableName string,
	keyColumn string,
) (repository.ShardKey, error) {
	err := s.validateProject(ctx, projectID)
	if err != nil {
		return repository.ShardKey{}, err
	}

	err = s.validateShardKeyColumn(
		ctx,
		projectID,
		tableName,
		keyColumn,
	)
	if err != nil {
		return repository.ShardKey{}, err
	}

	key, err := s.repo.AddShardKey(ctx, projectID, tableName, keyColumn)
	if err != nil {
		return repository.ShardKey{}, err
	}

	s.cache.ShardKeys.Set(key)

	return key, nil
}

func (s *ShardKeysService) DeleteShardKey(
	ctx context.Context,
	projectID string,
	keyID string,
) error {
	err := s.validateProject(ctx, projectID)
	if err != nil {
		return err
	}

	err = s.repo.DeleteShardKey(ctx, projectID, keyID)
	if err != nil {
		return err
	}

	return s.refresher.RefreshShardKeys(ctx, projectID)
}

func (s *ShardKeysService) GetShardKey(
	ctx context.Context,
	projectID string,
	tableName string,
) (repository.ShardKey, error) {
	err := s.validateProject(ctx, projectID)
	if err != nil {
		return repository.ShardKey{}, err
	}

	key, ok := s.cache.ShardKeys.Get(projectID, tableName)
	if ok {
		return key, nil
	}

	_, loaded := s.cache.ShardKeys.GetByProject(projectID)
	if !loaded {
		err = s.refresher.RefreshShardKeys(ctx, projectID)
		if err != nil {
			return repository.ShardKey{}, err
		}

		key, ok = s.cache.ShardKeys.Get(projectID, tableName)
	}

	if !ok {
		return repository.ShardKey{}, sql.ErrNoRows
	}

	return key, nil
}

func (s *ShardKeysService) ListShardKeys(
	ctx context.Context,
	projectID string,
) ([]repository.ShardKey, error) {
	err := s.validateProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	keys, loaded := s.cache.ShardKeys.GetByProject(projectID)
	if loaded {
		return keys, nil
	}

	err = s.refresher.RefreshShardKeys(ctx, projectID)
	if err != nil {
		return nil, err
	}

	keys, _ = s.cache.ShardKeys.GetByProject(projectID)

	return keys, nil
}
