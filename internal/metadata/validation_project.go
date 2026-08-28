package metadata

import "context"

// Checks if the project can be deleted
func (s *ProjectService) validateDeleteProject(
	ctx context.Context,
	projectID string,
) error {

	project, err := s.GetProjectByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project.Running {
		return ErrProjectRunning
	}

	err = s.refresher.RefreshNodes(ctx, projectID)
	if err != nil {
		return err
	}

	nodes, _ := s.cache.Nodes.GetByProject(projectID)
	if len(nodes) != 0 {
		return ErrProjectHasNodes
	}

	err = s.refresher.RefreshTopology(ctx, projectID)
	if err != nil {
		return err
	}

	topologies, _ := s.cache.Topology.GetByProjectID(projectID)
	if len(topologies) != 0 {
		return ErrProjectHasTopology
	}

	if err := s.refresher.RefreshShards(ctx, projectID); err != nil {
		return err
	}
	shards, _ := s.cache.Shards.GetByProject(projectID)
	if len(shards) != 0 {
		return ErrProjectHasShards
	}

	// Deferred until the corresponding state exists:
	// - Project must not have pending schema execution.
	// - Project must not have attached agents.

	return nil
}
