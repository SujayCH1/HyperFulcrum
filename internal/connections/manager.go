package connections

import "hyperfulcrum/internal/repository"

type ConnectionManager struct {
	store       *ConnectionStore
	projectRepo *repository.ProjectRepository
}

func (m *ConnectionManager) InitializeConnections() error {
	return nil
}
