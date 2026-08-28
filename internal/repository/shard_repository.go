package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	ShardStatusProvisioning  = "provisioning"
	ShardStatusActive        = "active"
	ShardStatusReconfiguring = "reconfiguring"
	ShardStatusUnavailable   = "unavailable"
)

type Shard struct {
	ID                 string    `json:"id"`
	ProjectID          string    `json:"project_id"`
	Name               string    `json:"shard_name"`
	Index              int       `json:"shard_index"`
	PrimaryNodeID      string    `json:"primary_node_id"`
	Status             string    `json:"status"`
	TopologyGeneration int64     `json:"topology_generation"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ShardRepository struct{ conn *sql.DB }

func NewShardRepository(conn *sql.DB) *ShardRepository { return &ShardRepository{conn: conn} }

func (r *ShardRepository) ShardAdd(ctx context.Context, projectID, name, primaryNodeID string) (Shard, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return Shard{}, err
	}
	defer tx.Rollback()

	var project string
	if err := tx.QueryRowContext(ctx, `SELECT id FROM projects WHERE id = $1 FOR UPDATE`, projectID).Scan(&project); err != nil {
		return Shard{}, err
	}
	var index int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(shard_index), -1) + 1 FROM shards WHERE project_id = $1`, projectID).Scan(&index); err != nil {
		return Shard{}, err
	}

	shard := Shard{ID: uuid.NewString(), ProjectID: projectID, Name: name, Index: index,
		PrimaryNodeID: primaryNodeID, Status: ShardStatusProvisioning, TopologyGeneration: 1}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO shards (id, project_id, shard_name, shard_index, primary_node_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at`, shard.ID, shard.ProjectID, shard.Name, shard.Index,
		shard.PrimaryNodeID, shard.Status).Scan(&shard.CreatedAt, &shard.UpdatedAt)
	if err != nil {
		return Shard{}, err
	}
	if err := tx.Commit(); err != nil {
		return Shard{}, err
	}
	return shard, nil
}

func (r *ShardRepository) ShardList(ctx context.Context, projectID string) ([]Shard, error) {
	rows, err := r.conn.QueryContext(ctx, shardSelect+` WHERE project_id = $1 ORDER BY shard_index`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shards := make([]Shard, 0)
	for rows.Next() {
		shard, err := scanShard(rows)
		if err != nil {
			return nil, err
		}
		shards = append(shards, shard)
	}
	return shards, rows.Err()
}

func (r *ShardRepository) ShardGetByID(ctx context.Context, shardID string) (Shard, error) {
	return scanShard(r.conn.QueryRowContext(ctx, shardSelect+` WHERE id = $1`, shardID))
}

func (r *ShardRepository) ShardGetByPrimaryNodeID(ctx context.Context, nodeID string) (Shard, error) {
	return scanShard(r.conn.QueryRowContext(ctx, shardSelect+` WHERE primary_node_id = $1`, nodeID))
}

func (r *ShardRepository) ShardUpdateName(ctx context.Context, shardID, name string) error {
	return requireAffected(r.conn.ExecContext(ctx, `UPDATE shards SET shard_name = $1, updated_at = NOW() WHERE id = $2`, name, shardID))
}

func (r *ShardRepository) ShardUpdateStatus(ctx context.Context, shardID, status string) error {
	if err := validateShardStatus(status); err != nil {
		return err
	}
	return requireAffected(r.conn.ExecContext(ctx, `UPDATE shards SET status = $1, updated_at = NOW() WHERE id = $2`, status, shardID))
}

func (r *ShardRepository) ShardSetPrimary(ctx context.Context, shardID, primaryNodeID string, expectedGeneration int64) (int64, error) {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var oldPrimaryNodeID string
	if err := tx.QueryRowContext(ctx, `
		SELECT primary_node_id FROM shards WHERE id = $1 AND topology_generation = $2 FOR UPDATE`,
		shardID, expectedGeneration).Scan(&oldPrimaryNodeID); err != nil {
		return 0, err
	}
	var promotable bool
	if err := tx.QueryRowContext(ctx, `
		SELECT EXISTS (SELECT 1 FROM node_topology WHERE shard_id = $1 AND replica_node_id = $2)`,
		shardID, primaryNodeID).Scan(&promotable); err != nil {
		return 0, err
	}
	if !promotable {
		return 0, sql.ErrNoRows
	}

	var generation int64
	err = tx.QueryRowContext(ctx, `
		UPDATE shards SET primary_node_id = $1, topology_generation = topology_generation + 1,
			status = $2, updated_at = NOW()
		WHERE id = $3 AND topology_generation = $4
		RETURNING topology_generation`, primaryNodeID, ShardStatusReconfiguring, shardID, expectedGeneration).Scan(&generation)
	if err != nil {
		return 0, err
	}

	// Swap the promoted standby and former primary while keeping all other
	// standby relationships attached to the same stable logical shard.
	if _, err := tx.ExecContext(ctx, `
		UPDATE node_topology
		SET shard_node_id = $1,
			replica_node_id = CASE WHEN replica_node_id = $1 THEN $2 ELSE replica_node_id END,
			updated_at = NOW()
		WHERE shard_id = $3`, primaryNodeID, oldPrimaryNodeID, shardID); err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE nodes SET node_type = CASE
			WHEN id = $1 THEN $3::ntype
			WHEN id = $2 THEN $4::ntype
		END
		WHERE id IN ($1, $2)`, oldPrimaryNodeID, primaryNodeID,
		NodeRoleStandby, NodeRolePrimary); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return generation, nil
}

func (r *ShardRepository) ShardRemove(ctx context.Context, shardID string) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var primaryNodeID string
	if err := tx.QueryRowContext(ctx, `DELETE FROM shards WHERE id = $1 RETURNING primary_node_id`, shardID).Scan(&primaryNodeID); err != nil {
		return err
	}
	if err := requireAffected(tx.ExecContext(ctx, `UPDATE nodes SET node_type = $1 WHERE id = $2`, NodeRoleUnassigned, primaryNodeID)); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ShardRepository) ShardRemoveAll(ctx context.Context, projectID string) error {
	_, err := r.conn.ExecContext(ctx, `DELETE FROM shards WHERE project_id = $1`, projectID)
	return err
}

const shardSelect = `SELECT id, project_id, shard_name, shard_index, primary_node_id,
	status, topology_generation, created_at, updated_at FROM shards`

type rowScanner interface{ Scan(...any) error }

func scanShard(row rowScanner) (Shard, error) {
	var shard Shard
	err := row.Scan(&shard.ID, &shard.ProjectID, &shard.Name, &shard.Index, &shard.PrimaryNodeID,
		&shard.Status, &shard.TopologyGeneration, &shard.CreatedAt, &shard.UpdatedAt)
	return shard, err
}

func validateShardStatus(status string) error {
	switch status {
	case ShardStatusProvisioning, ShardStatusActive, ShardStatusReconfiguring, ShardStatusUnavailable:
		return nil
	default:
		return fmt.Errorf("invalid shard status: %s", status)
	}
}

func requireAffected(result sql.Result, err error) error {
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil
}
