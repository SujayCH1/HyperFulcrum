package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSchemaLocked    = errors.New("project schema is locked")
	ErrSchemaActivated = errors.New("project schema has already been activated")
	ErrSchemaEmpty     = errors.New("project schema is empty")
	ErrSchemaHasKeys   = errors.New("project schema has shard keys")
	ErrSchemaNotLocked = errors.New("project schema is not locked")
	ErrSchemaRevision  = errors.New("project schema revision conflict")
)

type SchemaVersion struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	RawSQL      string     `json:"raw_sql"`
	Revision    int64      `json:"revision"`
	Locked      bool       `json:"locked"`
	ActivatedAt *time.Time `json:"activated_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (r *SchemaVersionRepository) LockSchema(
	ctx context.Context,
	projectID string,
) (SchemaVersion, error) {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return SchemaVersion{}, err
	}
	defer tx.Rollback()

	var existingProjectID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM projects WHERE id = $1 FOR UPDATE`,
		parsedProjectID,
	).Scan(&existingProjectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	schema := SchemaVersion{}
	selectQuery := `
		SELECT id, project_id, raw_sql, revision, locked, activated_at,
			created_at, updated_at
		FROM schema_versions
		WHERE project_id = $1
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, selectQuery, parsedProjectID).Scan(
		&schema.ID,
		&schema.ProjectID,
		&schema.RawSQL,
		&schema.Revision,
		&schema.Locked,
		&schema.ActivatedAt,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}
	if schema.ActivatedAt != nil {
		return SchemaVersion{}, ErrSchemaActivated
	}
	if schema.Locked {
		return schema, nil
	}

	var hasColumns bool
	err = tx.QueryRowContext(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM columns WHERE project_id = $1)`,
		parsedProjectID,
	).Scan(&hasColumns)
	if err != nil {
		return SchemaVersion{}, err
	}
	if !hasColumns {
		return SchemaVersion{}, ErrSchemaEmpty
	}

	updateQuery := `
		UPDATE schema_versions
		SET locked = TRUE,
			revision = revision + 1,
			updated_at = NOW()
		WHERE project_id = $1
		RETURNING revision, locked, updated_at
	`
	err = tx.QueryRowContext(ctx, updateQuery, parsedProjectID).Scan(
		&schema.Revision,
		&schema.Locked,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	err = tx.Commit()
	if err != nil {
		return SchemaVersion{}, err
	}

	return schema, nil
}

func (r *SchemaVersionRepository) UnlockSchema(
	ctx context.Context,
	projectID string,
) (SchemaVersion, error) {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return SchemaVersion{}, err
	}
	defer tx.Rollback()

	var existingProjectID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM projects WHERE id = $1 FOR UPDATE`,
		parsedProjectID,
	).Scan(&existingProjectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	schema := SchemaVersion{}
	query := `
		SELECT id, project_id, raw_sql, revision, locked, activated_at,
			created_at, updated_at
		FROM schema_versions
		WHERE project_id = $1
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, query, parsedProjectID).Scan(
		&schema.ID,
		&schema.ProjectID,
		&schema.RawSQL,
		&schema.Revision,
		&schema.Locked,
		&schema.ActivatedAt,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}
	if schema.ActivatedAt != nil {
		return SchemaVersion{}, ErrSchemaActivated
	}
	if !schema.Locked {
		return schema, nil
	}

	var hasKeys bool
	err = tx.QueryRowContext(
		ctx,
		`SELECT EXISTS (SELECT 1 FROM shard_keys WHERE project_id = $1)`,
		parsedProjectID,
	).Scan(&hasKeys)
	if err != nil {
		return SchemaVersion{}, err
	}
	if hasKeys {
		return SchemaVersion{}, ErrSchemaHasKeys
	}

	updateQuery := `
		UPDATE schema_versions
		SET locked = FALSE,
			revision = revision + 1,
			updated_at = NOW()
		WHERE project_id = $1
		RETURNING revision, locked, updated_at
	`
	err = tx.QueryRowContext(ctx, updateQuery, parsedProjectID).Scan(
		&schema.Revision,
		&schema.Locked,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	err = tx.Commit()
	if err != nil {
		return SchemaVersion{}, err
	}

	return schema, nil
}

type SchemaVersionRepository struct {
	db *sql.DB
}

func NewSchemaVersionRepository(db *sql.DB) *SchemaVersionRepository {
	return &SchemaVersionRepository{
		db: db,
	}
}

func (r *SchemaVersionRepository) ReplaceProjectSchema(
	ctx context.Context,
	projectID string,
	rawSQL string,
	columns []Column,
	edges []FkEdges,
	expectedRevision int64,
) (SchemaVersion, error) {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return SchemaVersion{}, err
	}
	defer tx.Rollback()

	var existingProjectID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT id FROM projects WHERE id = $1 FOR UPDATE`,
		parsedProjectID,
	).Scan(&existingProjectID)
	if err != nil {
		return SchemaVersion{}, err
	}

	var locked bool
	var activatedAt *time.Time
	var currentRevision int64
	stateQuery := `
		SELECT locked, activated_at, revision
		FROM schema_versions
		WHERE project_id = $1
		FOR UPDATE
	`
	err = tx.QueryRowContext(ctx, stateQuery, parsedProjectID).Scan(
		&locked,
		&activatedAt,
		&currentRevision,
	)
	if errors.Is(err, sql.ErrNoRows) {
		currentRevision = 0
	} else if err != nil {
		return SchemaVersion{}, err
	}
	if currentRevision != expectedRevision {
		return SchemaVersion{}, ErrSchemaRevision
	}
	if activatedAt != nil {
		return SchemaVersion{}, ErrSchemaActivated
	}
	if locked {
		return SchemaVersion{}, ErrSchemaLocked
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM fk_edges WHERE project_id = $1`,
		parsedProjectID,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM columns WHERE project_id = $1`,
		parsedProjectID,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	columnQuery := `
		INSERT INTO columns
		(project_id, table_name, column_name, data_type, nullable, is_primary_key)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, column := range columns {
		_, err = tx.ExecContext(
			ctx,
			columnQuery,
			parsedProjectID,
			column.TableName,
			column.ColumnName,
			column.DataType,
			column.IsNullable,
			column.IsPk,
		)
		if err != nil {
			return SchemaVersion{}, err
		}
	}

	edgeQuery := `
		INSERT INTO fk_edges
		(project_id, parent_table, parent_column, child_table, child_column)
		VALUES ($1, $2, $3, $4, $5)
	`
	for _, edge := range edges {
		_, err = tx.ExecContext(
			ctx,
			edgeQuery,
			parsedProjectID,
			edge.ParentTable,
			edge.ParentColumn,
			edge.ChildTable,
			edge.ChildColumn,
		)
		if err != nil {
			return SchemaVersion{}, err
		}
	}

	now := time.Now()
	schema := SchemaVersion{}
	schemaQuery := `
		INSERT INTO schema_versions
		(id, project_id, raw_sql, revision, locked, created_at, updated_at)
		VALUES ($1, $2, $3, 1, FALSE, $4, $4)
		ON CONFLICT (project_id)
		DO UPDATE SET
			raw_sql = EXCLUDED.raw_sql,
			revision = schema_versions.revision + 1,
			updated_at = EXCLUDED.updated_at
		RETURNING id, project_id, raw_sql, revision, locked, activated_at,
			created_at, updated_at
	`
	err = tx.QueryRowContext(
		ctx,
		schemaQuery,
		uuid.NewString(),
		parsedProjectID,
		rawSQL,
		now,
	).Scan(
		&schema.ID,
		&schema.ProjectID,
		&schema.RawSQL,
		&schema.Revision,
		&schema.Locked,
		&schema.ActivatedAt,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	err = tx.Commit()
	if err != nil {
		return SchemaVersion{}, err
	}

	return schema, nil
}

func (s *SchemaVersionRepository) FetchSchema(
	ctx context.Context,
	projectID string,
) (SchemaVersion, error) {
	query := `
		SELECT
		id, project_id, raw_sql, revision, locked, activated_at,
		created_at, updated_at
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
		&schema.Revision,
		&schema.Locked,
		&schema.ActivatedAt,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)
	if err != nil {
		return SchemaVersion{}, err
	}

	return schema, nil
}
