package metadata

import (
	"context"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type ColumnService struct {
	repo      *repository.ColumnRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewColumnService(
	repo *repository.ColumnRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *ColumnService {

	return &ColumnService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *ColumnService) ReplaceColumns(
	ctx context.Context,
	projectID string,
	columns []repository.Column,
) error {

	if err := s.repo.ColumnReplace(
		ctx,
		projectID,
		columns,
	); err != nil {
		return err
	}

	return s.refresher.RefreshColumns(
		ctx,
		projectID,
	)
}

func (s *ColumnService) ListColumns(
	ctx context.Context,
	projectID string,
) ([]repository.Column, error) {

	columns, loaded := s.cache.Columns.GetByProject(projectID)
	if loaded {
		return columns, nil
	}

	if err := s.refresher.RefreshColumns(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	columns, _ = s.cache.Columns.GetByProject(projectID)

	return columns, nil
}
