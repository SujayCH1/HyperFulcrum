package metadata

import (
	"context"
	"database/sql"

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

	if err := s.refresher.RefreshProjects(ctx); err != nil {
		return repository.Project{}, err
	}

	return project, nil
}

func (s *ProjectService) ListProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	projects := s.cache.Projects.GetAll()

	if len(projects) > 0 {
		return projects, nil
	}

	if err := s.refresher.RefreshProjects(ctx); err != nil {
		return nil, err
	}

	return s.cache.Projects.GetAll(), nil
}

func (s *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
) (repository.Project, error) {

	project, ok := s.cache.Projects.Get(projectID)
	if ok {
		return project, nil
	}

	if err := s.refresher.RefreshProject(
		ctx,
		projectID,
	); err != nil {
		return repository.Project{}, err
	}

	project, ok = s.cache.Projects.Get(projectID)
	if !ok {
		return repository.Project{}, sql.ErrNoRows
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(
	ctx context.Context,
	projectID string,
) error {

	if err := s.validateDeleteProject(
		ctx,
		projectID,
	); err != nil {
		return err
	}

	if err := s.repo.ProjectRemove(
		ctx,
		projectID,
	); err != nil {
		return err
	}

	return s.refresher.RefreshProjects(ctx)
}

func (s *ProjectService) GetReadyProjects(
	ctx context.Context,
) ([]repository.Project, error) {

	return s.repo.ProjectGetReady(ctx)
}

// Remaining functions
// ActivateProjecy()
// DeactivateProject()
// UpdateProject()
// ArchiveProject()
