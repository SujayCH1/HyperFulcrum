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

	cachedColumns := s.cache.Columns.GetAll()

	var result []repository.Column

	for _, column := range cachedColumns {
		if column.ProjectID == projectID {
			result = append(result, column)
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	if err := s.refresher.RefreshColumns(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	cachedColumns = s.cache.Columns.GetAll()

	result = nil

	for _, column := range cachedColumns {
		if column.ProjectID == projectID {
			result = append(result, column)
		}
	}

	return result, nil
}
