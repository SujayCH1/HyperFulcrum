package connections

import (
	"database/sql"
	"errors"
	"sync"
)

var (
	ErrProjectPoolsNotFound = errors.New("connection pools for project not found")
	ErrNodePoolNotFound     = errors.New("connection pool for node not found")
	ErrNilPool              = errors.New("connection pool cannot be nil")
)

type PoolStore struct {
	mu    sync.RWMutex
	pools map[string]map[string]*sql.DB
}

func NewPoolStore() *PoolStore {
	return &PoolStore{
		pools: make(map[string]map[string]*sql.DB),
	}
}

func (s *PoolStore) Set(projectID string, nodeID string, db *sql.DB) error {
	if db == nil {
		return ErrNilPool
	}

	s.mu.Lock()

	projectPools, ok := s.pools[projectID]
	if !ok {
		projectPools = make(map[string]*sql.DB)
		s.pools[projectID] = projectPools
	}

	oldDB := projectPools[nodeID]
	projectPools[nodeID] = db

	s.mu.Unlock()

	if oldDB != nil && oldDB != db {
		return oldDB.Close()
	}

	return nil
}

func (s *PoolStore) Get(projectID string, nodeID string) (*sql.DB, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	projectPools, ok := s.pools[projectID]
	if !ok {
		return nil, ErrProjectPoolsNotFound
	}

	db, ok := projectPools[nodeID]
	if !ok {
		return nil, ErrNodePoolNotFound
	}

	return db, nil
}

func (s *PoolStore) Remove(projectID string, nodeID string) error {
	s.mu.Lock()

	projectPools, ok := s.pools[projectID]
	if !ok {
		s.mu.Unlock()
		return nil
	}

	db := projectPools[nodeID]
	delete(projectPools, nodeID)

	if len(projectPools) == 0 {
		delete(s.pools, projectID)
	}

	s.mu.Unlock()

	if db != nil {
		return db.Close()
	}

	return nil
}

func (s *PoolStore) RemoveProject(projectID string) error {
	s.mu.Lock()
	projectPools := s.pools[projectID]
	delete(s.pools, projectID)
	s.mu.Unlock()

	return closePools(projectPools)
}

func (s *PoolStore) ReplaceAll(pools map[string]map[string]*sql.DB) error {
	if pools == nil {
		pools = make(map[string]map[string]*sql.DB)
	}

	s.mu.Lock()
	oldPools := s.pools
	s.pools = pools
	s.mu.Unlock()

	return closePoolStore(oldPools, pools)
}

func (s *PoolStore) Close() error {
	return s.ReplaceAll(make(map[string]map[string]*sql.DB))
}

func closePoolStore(
	pools map[string]map[string]*sql.DB,
	keep map[string]map[string]*sql.DB,
) error {
	var closeErr error

	for projectID, projectPools := range pools {
		for nodeID, db := range projectPools {
			if keepPools, ok := keep[projectID]; ok && keepPools[nodeID] == db {
				continue
			}

			closeErr = errors.Join(closeErr, db.Close())
		}
	}

	return closeErr
}

func closePools(pools map[string]*sql.DB) error {
	var closeErr error

	for _, db := range pools {
		closeErr = errors.Join(closeErr, db.Close())
	}

	return closeErr
}
