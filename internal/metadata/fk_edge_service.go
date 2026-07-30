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

	edges := s.cache.FKEdges.GetByProject(projectID)
	if len(edges) > 0 {
		return edges, nil
	}

	if err := s.refresher.RefreshFKEdges(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	edges = s.cache.FKEdges.GetByProject(projectID)

	return edges, nil
}

func (s *FKEdgesService) ListByParentTable(
	ctx context.Context,
	projectID string,
	tableName string,
) ([]repository.FkEdges, error) {

	edges := s.cache.FKEdges.GetByTable(
		projectID,
		tableName,
	)

	if len(edges) == 0 {
		if err := s.refresher.RefreshFKEdges(
			ctx,
			projectID,
		); err != nil {
			return nil, err
		}

		edges = s.cache.FKEdges.GetByTable(
			projectID,
			tableName,
		)
	}

	result := make([]repository.FkEdges, 0)

	for _, edge := range edges {
		if edge.ParentTable == tableName {
			result = append(result, edge)
		}
	}

	return result, nil
}

func (s *FKEdgesService) ListByChildTable(
	ctx context.Context,
	projectID string,
	tableName string,
) ([]repository.FkEdges, error) {

	edges := s.cache.FKEdges.GetByTable(
		projectID,
		tableName,
	)

	if len(edges) == 0 {
		if err := s.refresher.RefreshFKEdges(
			ctx,
			projectID,
		); err != nil {
			return nil, err
		}

		edges = s.cache.FKEdges.GetByTable(
			projectID,
			tableName,
		)
	}

	result := make([]repository.FkEdges, 0)

	for _, edge := range edges {
		if edge.ChildTable == tableName {
			result = append(result, edge)
		}
	}

	return result, nil
}
