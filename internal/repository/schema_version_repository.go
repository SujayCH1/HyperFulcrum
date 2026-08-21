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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query1 := `
		DELETE FROM schema_versions
		WHERE project_id = $1
	`

	_, err = tx.ExecContext(
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

	now := time.Now()

	_, err = tx.ExecContext(
		ctx,
		query2,
		uuid.New(),
		projectID,
		rawSQL,
		now,
		now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()

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
		return SchemaVersion{}, err
	}

	return schema, nil
}

func (r *SchemaVersionRepository) DeleteSchema(
	ctx context.Context,
	projectID string,
) error {
	query := `
		DELETE FROM schema_versions
		WHERE project_id = $1
	`

	res, err := r.db.ExecContext(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
