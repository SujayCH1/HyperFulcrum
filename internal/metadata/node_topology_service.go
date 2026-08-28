package metadata

import (
	"context"
	"database/sql"
	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type TopologyService struct {
	repo      *repository.NodeTopologyRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func topologyPrimaryNodeID(cacheManager *cache.CacheManager, shardID string) string {
	shard, _ := cacheManager.Shards.GetByID(shardID)
	return shard.PrimaryNodeID
}

func NewTopologyService(
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

// NewTpologyService is retained for compatibility with older callers.
func NewTpologyService(repo *repository.NodeTopologyRepository, cacheManager *cache.CacheManager,
	refresher *cache.CacheRefresher) *TopologyService {
	return NewTopologyService(repo, cacheManager, refresher)
}

func (s *TopologyService) CreateTopology(
	ctx context.Context,
	projectID string,
	shardID string,
	standbyID string,
) (repository.NodeTopology, error) {

	if err := s.validateCreateTopology(
		ctx,
		projectID,
		shardID,
		standbyID,
	); err != nil {
		return repository.NodeTopology{}, err
	}

	topology, err := s.repo.TopologyAdd(
		ctx,
		projectID,
		topologyPrimaryNodeID(s.cache, shardID),
		standbyID,
	)
	if err != nil {
		return repository.NodeTopology{}, err
	}

	if err := s.refresher.RefreshTopology(
		ctx,
		projectID,
	); err != nil {
		return repository.NodeTopology{}, err
	}

	return topology, nil
}

func (s *TopologyService) DeleteTopology(
	ctx context.Context,
	relationID string,
	projectID string,
) error {

	topology, err := s.GetTopologyByID(
		ctx,
		projectID,
		relationID,
	)
	if err != nil {
		return err
	}

	if err := s.validateDeleteTopology(
		ctx,
		topology,
	); err != nil {
		return err
	}

	if err := s.repo.TopologyRemove(
		ctx,
		relationID,
	); err != nil {
		return err
	}

	return s.refresher.RefreshTopology(
		ctx,
		projectID,
	)
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

	return repository.NodeTopology{}, sql.ErrNoRows
}
