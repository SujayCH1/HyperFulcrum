package metadata

import (
	"context"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type FKEdgesService struct {
	repo      *repository.FKEdgesRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewFKEdgesService(
	repo *repository.FKEdgesRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *FKEdgesService {

	return &FKEdgesService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *FKEdgesService) ReplaceEdges(
	ctx context.Context,
	projectID string,
	edges []repository.FkEdges,
) error {

	if err := s.repo.FKEdgesReplaceForProject(
		ctx,
		projectID,
		edges,
	); err != nil {
		return err
	}

	return s.refresher.RefreshFKEdges(
		ctx,
		projectID,
	)
}

func (s *FKEdgesService) ListEdges(
	ctx context.Context,
	projectID string,
) ([]repository.FkEdges, error) {

	cachedEdges := s.cache.FKEdges.GetAll()

	var result []repository.FkEdges

	for _, edge := range cachedEdges {
		if edge.ProjectId == projectID {
			result = append(result, edge)
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	if err := s.refresher.RefreshFKEdges(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	cachedEdges = s.cache.FKEdges.GetAll()

	result = nil

	for _, edge := range cachedEdges {
		if edge.ProjectId == projectID {
			result = append(result, edge)
		}
	}

	return result, nil
}

func (s *FKEdgesService) ListByParentTable(
	ctx context.Context,
	projectID string,
	tableName string,
) ([]repository.FkEdges, error) {

	return s.repo.EdgesListByParentTable(
		ctx,
		projectID,
		tableName,
	)
}

func (s *FKEdgesService) ListByChildTable(
	ctx context.Context,
	projectID string,
	tableName string,
) ([]repository.FkEdges, error) {

	return s.repo.EdgesListByChildTable(
		ctx,
		projectID,
		tableName,
	)
}
