package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SchemaVersion struct {
	ID        string
	ProjectID string
	RawSQL    string
	CreatedAt string
	UpdatedAt string
}

type SchemaVersionRepository struct {
	db *sql.DB
}

func NewSchemaVersionRepository(db *sql.DB) *SchemaVersionRepository {
	return &SchemaVersionRepository{
		db: db,
	}
}

func (r *SchemaVersionRepository) ReplaceSchema(
	ctx context.Context,
	projectID string,
	rawSQL string,
) error {
	query1 := `
		DELETE FROM schema_versions
		WHERE project_id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query1,
		projectID,
	)
	if err != nil {
		return err
	}

	query2 := `
		INSERT INTO schema_versions
		(id, project_id, raw_sql, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.ExecContext(
		ctx,
		query2,
		uuid.New(),
		projectID,
		rawSQL,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil

}

func (s *SchemaVersionRepository) FetchSchema(
	ctx context.Context,
	projectID string,
) (SchemaVersion, error) {
	query := `
		SELECT
		id, project_id, raw_sql, created_at, updated_at
		FROM schema_versions
		WHERE project_id = $1
	`

	row := s.db.QueryRowContext(
		ctx,
		query,
		projectID,
	)

	var schema SchemaVersion

	err := row.Scan(
		&schema.ID,
		&schema.ProjectID,
		&schema.RawSQL,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return SchemaVersion{}, nil
		}

		return SchemaVersion{}, err
	}

	return schema, nil
}

func (r *SchemaVersionRepository) DeleteSchema(
	ctx context.Context,
	projectID string,
) error {
	query := `
		DELETE FROM project_schema
		WHERE project_id = $1
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return err
	}

	return nil
}
