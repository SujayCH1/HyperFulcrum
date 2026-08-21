package metadata

import (
	"context"
	"database/sql"
	"errors"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type SchemaVersionService struct {
	repo      *repository.SchemaVersionRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func (s *SchemaVersionService) ReplaceProjectSchema(
	ctx context.Context,
	projectID string,
	rawSQL string,
	columns []repository.Column,
	edges []repository.FkEdges,
	expectedRevision int64,
) (repository.SchemaVersion, error) {
	schema, err := s.repo.ReplaceProjectSchema(
		ctx,
		projectID,
		rawSQL,
		columns,
		edges,
		expectedRevision,
	)
	if err != nil {
		if errors.Is(err, repository.ErrSchemaRevision) {
			refreshErr := s.refresher.RefreshSchemaVersion(ctx, projectID)
			if refreshErr != nil {
				return repository.SchemaVersion{}, errors.Join(err, refreshErr)
			}
		}
		return repository.SchemaVersion{}, err
	}

	s.cache.Columns.ReplaceProject(projectID, columns)
	s.cache.FKEdges.ReplaceProject(projectID, edges)
	s.cache.SchemaVersion.Set(schema)

	return schema, nil
}

func (s *SchemaVersionService) LockSchema(
	ctx context.Context,
	projectID string,
) (repository.SchemaVersion, error) {
	schema, err := s.repo.LockSchema(ctx, projectID)
	if err != nil {
		return repository.SchemaVersion{}, err
	}

	s.cache.SchemaVersion.Set(schema)

	return schema, nil
}

func (s *SchemaVersionService) UnlockSchema(
	ctx context.Context,
	projectID string,
) (repository.SchemaVersion, error) {
	schema, err := s.repo.UnlockSchema(ctx, projectID)
	if err != nil {
		return repository.SchemaVersion{}, err
	}

	s.cache.SchemaVersion.Set(schema)

	return schema, nil
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
