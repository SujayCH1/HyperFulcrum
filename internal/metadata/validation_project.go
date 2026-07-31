package metadata

import "context"

// Checks if the project can be deleted
func (s *ProjectService) validateDeleteProject(
	ctx context.Context,
	projectID string,
) error {

	// Ensure the project exists.
	if _, err := s.GetProjectByID(ctx, projectID); err != nil {
		return err
	}

	// Future validations:
	// - Project must not be active.
	// - Project must not contain any nodes.
	// - Project must not have replication topology.
	// - Project must not have pending schema execution.
	// - Project must not have attached agents.

	return nil
}
