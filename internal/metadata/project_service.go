package metadata

import (
	"context"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/repository"
)

type ProjectService struct {
	repo      *repository.ProjectRepository
	cache     *cache.CacheManager
	refresher *cache.CacheRefresher
}

func NewProjectService(
	repo *repository.ProjectRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
) *ProjectService {

	return &ProjectService{
		repo:      repo,
		cache:     cache,
		refresher: refresher,
	}
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	name string,
	description string,
) (repository.Project, error) {

	project, err := s.repo.ProjectAdd(
		ctx,
		name,
		description,
	)
	if err != nil {
		return repository.Project{}, err
	}

	s.cache.Projects.Set(project)

	return project, nil
}

func (s *ProjectService) ListProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	projects, err := s.repo.ProjectList(ctx)
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		s.cache.Projects.Set(project)
	}

	return projects, nil
}

func (s *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
) (repository.Project, error) {

	if project, ok := s.cache.Projects.Get(projectID); ok {
		return project, nil
	}

	project, err := s.repo.ProjectGetByID(
		ctx,
		projectID,
	)
	if err != nil {
		return repository.Project{}, err
	}

	s.cache.Projects.Set(project)

	return project, nil
}

func (s *ProjectService) DeleteProject(
	ctx context.Context,
	projectID string,
) error {

	err := s.repo.ProjectRemove(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	s.cache.Projects.Delete(projectID)

	return nil
}

func (s *ProjectService) GetReadyProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	return s.repo.ProjectGetReady(ctx)
}
