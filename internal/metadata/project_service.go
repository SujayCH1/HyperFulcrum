package metadata

import (
	"context"

	"hyperfulcrum/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(
	repo *repository.ProjectRepository,
) *ProjectService {

	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	name string,
	description string,
) (repository.Project, error) {

	return s.repo.ProjectAdd(
		ctx,
		name,
		description,
	)
}

func (s *ProjectService) ListProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	return s.repo.ProjectList(ctx)
}

func (s *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
) (repository.Project, error) {

	return s.repo.ProjectGetByID(
		ctx,
		projectID,
	)
}

func (s *ProjectService) DeleteProject(
	ctx context.Context,
	projectID string,
) error {

	return s.repo.ProjectRemove(
		ctx,
		projectID,
	)
}

func (s *ProjectService) GetReadyProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	return s.repo.ProjectGetReady(ctx)
}
