package repository

import (
	"context"
	"database/sql"
)

type ShardKeyRecord struct {
	TableName      string `json:"table_name"`
	ShardKeyColumn string `json:"shard_key_column"`
	IsManual       bool   `json:"is_manual"`
}

type ShardKeysRepository struct {
	db *sql.DB
}

func NewShardKeysRepository(db *sql.DB) *ShardKeysRepository {
	return &ShardKeysRepository{db: db}
}

func (r *ShardKeysRepository) GetShardKeysByProjectID(
	ctx context.Context,
	projectID string,
) ([]ShardKeyRecord, error) {
	query := `
		SELECT table_name, shard_key_column, is_manual
		FROM table_shard_keys
		WHERE project_id = $1`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ShardKeyRecord
	for rows.Next() {
		var rec ShardKeyRecord
		if err := rows.Scan(&rec.TableName, &rec.ShardKeyColumn, &rec.IsManual); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

func (r *ShardKeysRepository) ReplaceShardKeysForProject(
	ctx context.Context,
	projectID string,
	records []ShardKeyRecord,
) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`DELETE FROM table_shard_keys WHERE project_id = $1 AND is_manual = FALSE`,
		projectID,
	)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO table_shard_keys (project_id, table_name, shard_key_column, is_manual)
		VALUES ($1, $2, $3, $4)`

	for _, rec := range records {
		_, err = tx.ExecContext(ctx, insertQuery, projectID, rec.TableName, rec.ShardKeyColumn, rec.IsManual)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}