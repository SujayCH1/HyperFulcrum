package metadata

import (
	"context"
	"database/sql"
	"errors"

	"hyperfulcrum/internal/cache"
	"hyperfulcrum/internal/connections"
	"hyperfulcrum/internal/repository"
)

type ProjectService struct {
	repo              *repository.ProjectRepository
	cache             *cache.CacheManager
	refresher         *cache.CacheRefresher
	connectionManager *connections.ConnectionManager
}

func NewProjectService(
	repo *repository.ProjectRepository,
	cache *cache.CacheManager,
	refresher *cache.CacheRefresher,
	connectionManager *connections.ConnectionManager,
) *ProjectService {

	return &ProjectService{
		repo:              repo,
		cache:             cache,
		refresher:         refresher,
		connectionManager: connectionManager,
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

	projects, loaded := s.cache.Projects.GetAll()

	if loaded {
		return projects, nil
	}

	if err := s.refresher.RefreshProjects(ctx); err != nil {
		return nil, err
	}

	projects, _ = s.cache.Projects.GetAll()

	return projects, nil
}

func (s *ProjectService) GetProjectByID(
	ctx context.Context,
	projectID string,
) (repository.Project, error) {

	project, ok := s.cache.Projects.Get(projectID)
	if ok {
		return project, nil
	}

	if s.cache.Projects.Loaded() {
		return repository.Project{}, sql.ErrNoRows
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

	connectionErr := s.connectionManager.RemoveProject(projectID)
	s.cache.DeleteProject(projectID)

	return errors.Join(
		connectionErr,
		s.refresher.RefreshProjects(ctx),
	)
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
