package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ShardKey struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	TableName string `json:"table_name"`
	KeyColumn string `json:"key_column"`

	UpdatedAt time.Time `json:"updated_at"`
}

type ShardKeyRepository struct {
	conn *sql.DB
}

func NewShardKeyRepository(conn *sql.DB) *ShardKeyRepository {
	return &ShardKeyRepository{conn: conn}
}

func (r *ShardKeyRepository) AddShardKey(
	ctx context.Context,
	projectID string,
	tableName string,
	keyColumn string,
) (ShardKey, error) {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return ShardKey{}, err
	}

	key := ShardKey{
		ID:        uuid.NewString(),
		ProjectID: parsedProjectID.String(),
		TableName: tableName,
		KeyColumn: keyColumn,
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO shard_keys (id, project_id, table_name, key_column, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING updated_at
	`

	err = r.conn.QueryRowContext(
		ctx,
		query,
		key.ID,
		key.ProjectID,
		key.TableName,
		key.KeyColumn,
		key.UpdatedAt,
	).Scan(&key.UpdatedAt)
	if err != nil {
		return ShardKey{}, err
	}

	return key, nil
}

func (r *ShardKeyRepository) DeleteShardKey(
	ctx context.Context,
	projectID string,
	id string,
) error {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return err
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	query := `
		DELETE FROM shard_keys
		WHERE project_id = $1
		AND id = $2
	`
	result, err := r.conn.ExecContext(
		ctx,
		query,
		parsedProjectID,
		parsedID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *ShardKeyRepository) FetchShardKeys(
	ctx context.Context,
	projectID string,
) ([]ShardKey, error) {
	parsedProjectID, err := uuid.Parse(projectID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, project_id, table_name, key_column, updated_at
		FROM shard_keys
		WHERE project_id = $1
		ORDER BY table_name
	`

	rows, err := r.conn.QueryContext(ctx, query, parsedProjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	keys := make([]ShardKey, 0)
	for rows.Next() {
		var key ShardKey
		err := rows.Scan(
			&key.ID,
			&key.ProjectID,
			&key.TableName,
			&key.KeyColumn,
			&key.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return keys, nil
}
