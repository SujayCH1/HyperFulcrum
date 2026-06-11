package repository

import (
	"context"
	"database/sql"
	"time"
)

type NodeKey struct {
	ProjectID     string `json:"project_id"`
	TableName     string `json:"table_name"`
	NodeKeyColumn string `json:"node_key_column"`

	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

type NodeKeyRepository struct {
	conn *sql.DB
}

func NewNodeKeyRepository(connConfig *sql.DB) *NodeKeyRepository {
	return &NodeKeyRepository{conn: connConfig}
}

func (r *NodeKeyRepository) NodeKeyFetchByProjectID(ctx context.Context, projectID string) ([]NodeKey, error) {
	query := `SELECT project_id,table_name,node_key_column,updated_at,created_at FROM node_keys WHERE project_id = $1`
	rows, err := r.conn.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodeKeys []NodeKey

	for rows.Next() {
		var nk NodeKey
		if err := rows.Scan(&nk.ProjectID, &nk.TableName, &nk.NodeKeyColumn, &nk.UpdatedAt, &nk.CreatedAt); err != nil {
			return nil, err
		}
		nodeKeys = append(nodeKeys, nk)
	}

	return nodeKeys, nil
}

func (r *NodeKeyRepository) NodeKeyReplaceByProjectID(ctx context.Context, projectID string, records []NodeKey) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	deleteQuery := `DELETE FROM node_keys WHERE project_id = $1`

	if _, err := tx.ExecContext(ctx, deleteQuery, projectID); err != nil {
		return err
	}

	//UPSERT
	insertQuery := `INSERT INTO node_keys (project_id, table_name, node_key_column, updated_at, created_at) 
	VALUES ($1, $2, $3, $4, $5) ON CONFLICT 
	(project_id, table_name) DO UPDATE SET 
	node_key_column = EXCLUDED.node_key_column, 
	updated_at = EXCLUDED.updated_at`

	for _, record := range records {
		_, err := tx.ExecContext(ctx, insertQuery, record.ProjectID, record.TableName, record.NodeKeyColumn, record.UpdatedAt, record.CreatedAt)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
