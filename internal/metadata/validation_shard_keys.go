package metadata

import (
	"context"
	"database/sql"
)

func (s *ShardKeysService) validateProject(
	ctx context.Context,
	projectID string,
) error {
	_, ok := s.cache.Projects.Get(projectID)
	if ok {
		return nil
	}

	if s.cache.Projects.Loaded() {
		return sql.ErrNoRows
	}

	err := s.refresher.RefreshProject(ctx, projectID)
	if err != nil {
		return err
	}

	_, ok = s.cache.Projects.Get(projectID)
	if !ok {
		return sql.ErrNoRows
	}

	return nil
}

func (s *ShardKeysService) validateShardKeyColumn(
	ctx context.Context,
	projectID string,
	tableName string,
	columnName string,
) error {
	_, ok := s.cache.Columns.Get(projectID, tableName, columnName)
	if ok {
		return nil
	}

	_, loaded := s.cache.Columns.GetByProject(projectID)
	if !loaded {
		err := s.refresher.RefreshColumns(ctx, projectID)
		if err != nil {
			return err
		}

		_, ok = s.cache.Columns.Get(projectID, tableName, columnName)
	}

	if !ok {
		return sql.ErrNoRows
	}

	return nil
}

func (s *ShardKeysService) validateSchemaLocked(
	ctx context.Context,
	projectID string,
) error {
	schema, ok := s.cache.SchemaVersion.Get(projectID)
	if !ok {
		err := s.refresher.RefreshSchemaVersion(ctx, projectID)
		if err != nil {
			return err
		}
		schema, ok = s.cache.SchemaVersion.Get(projectID)
	}

	if !ok {
		return sql.ErrNoRows
	}
	if !schema.Locked {
		return ErrSchemaNotLocked
	}

	return nil
}
