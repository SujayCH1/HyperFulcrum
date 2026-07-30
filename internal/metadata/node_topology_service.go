package metadata

import (
	"context"
	"fmt"
	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type TopologyService struct {
	repo      *repository.NodeTopologyRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewTpologyService(
	repo *repository.NodeTopologyRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *TopologyService {
	return &TopologyService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *TopologyService) CreateTopology(
	ctx context.Context,
	projectID string,
	replicaID string,
	shardID string,
) (repository.NodeTopology, error) {

	topology, err := s.repo.TopologyAdd(
		ctx,
		projectID,
		shardID,
		replicaID,
	)
	if err != nil {
		return repository.NodeTopology{}, err
	}

	err = s.refresher.RefreshTopology(ctx, projectID)
	if err != nil {
		return repository.NodeTopology{}, err
	}

	return topology, err
}

func (s *TopologyService) DeleteTopology(
	ctx context.Context,
	relationID string,
	projectID string,
) error {
	err := s.repo.TopologyRemove(
		ctx,
		relationID,
	)
	if err != nil {
		return err
	}

	err = s.refresher.RefreshTopology(ctx, projectID)
	if err != nil {
		return err
	}

	return nil

}

func (s *TopologyService) ListTopologies(
	ctx context.Context,
	projectID string,
) ([]repository.NodeTopology, error) {

	topologies, ok := s.cache.Topology.GetByProjectID(projectID)
	if ok {
		return topologies, nil
	}

	if err := s.refresher.RefreshTopology(
		ctx,
		projectID,
	); err != nil {
		return nil, err
	}

	topologies, _ = s.cache.Topology.GetByProjectID(projectID)

	return topologies, nil
}

func (s *TopologyService) GetTopologyByID(
	ctx context.Context,
	projectID string,
	relationID string,
) (repository.NodeTopology, error) {

	topologies, ok := s.cache.Topology.GetByProjectID(projectID)

	if !ok {
		if err := s.refresher.RefreshTopology(
			ctx,
			projectID,
		); err != nil {
			return repository.NodeTopology{}, err
		}

		topologies, _ = s.cache.Topology.GetByProjectID(projectID)
	}

	for _, topology := range topologies {
		if topology.RelationID == relationID {
			return topology, nil
		}
	}

	return repository.NodeTopology{}, fmt.Errorf("topology not found")
}
