package metadata

import (
	"context"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

func ensureNodeOutsideTopology(
	ctx context.Context,
	refresher *cache.CacheRefresher,
	cacheManager *cache.CacheManager,
	node repository.Node,
) error {
	topologies, loaded := cacheManager.Topology.GetByProjectID(node.ProjectID)
	if !loaded {
		err := refresher.RefreshTopology(ctx, node.ProjectID)
		if err != nil {
			return err
		}
		topologies, _ = cacheManager.Topology.GetByProjectID(node.ProjectID)
	}

	for _, topology := range topologies {
		if topology.PrimaryNodeID == node.ID || topology.StandbyNodeID == node.ID {
			return ErrNodeInTopology
		}
	}

	return nil
}
