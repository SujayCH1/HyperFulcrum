package metadata

import (
	"context"
	"database/sql"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type SchemaVersionService struct {
	repo      *repository.SchemaVersionRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewSchemaVersionService(
	repo *repository.SchemaVersionRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *SchemaVersionService {

	return &SchemaVersionService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *SchemaVersionService) ReplaceSchemaVersion(
	ctx context.Context,
	projectID string,
	rawSQL string,
) error {

	err := s.repo.ReplaceSchema(
		ctx,
		projectID,
		rawSQL,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshSchemaVersion(
		ctx,
		projectID,
	)
}

func (s *SchemaVersionService) GetSchemaVersion(
	ctx context.Context,
	projectID string,
) (repository.SchemaVersion, error) {

	schema, ok := s.cache.SchemaVersion.Get(projectID)
	if ok {
		return schema, nil
	}

	if s.cache.SchemaVersion.Loaded(projectID) {
		return repository.SchemaVersion{}, sql.ErrNoRows
	}

	if err := s.refresher.RefreshSchemaVersion(
		ctx,
		projectID,
	); err != nil {
		return repository.SchemaVersion{}, err
	}

	schema, ok = s.cache.SchemaVersion.Get(projectID)
	if !ok {
		return repository.SchemaVersion{}, sql.ErrNoRows
	}

	return schema, nil
}

func (s *SchemaVersionService) DeleteSchemaVersion(
	ctx context.Context,
	projectID string,
) error {

	err := s.repo.DeleteSchema(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	return s.refresher.RefreshSchemaVersion(
		ctx,
		projectID,
	)
}
