package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type NodeTopology struct {
	RelationID          string `json:"relation_id"`
	ProjectID           string `json:"project_id"`
	ShardID             string `json:"shard_id"`
	PrimaryNodeID       string `json:"primary_node_id"`
	StandbyNodeID       string `json:"standby_node_id"`
	Status              string `json:"relationship_status"`
	ReplicationSlotName string `json:"replication_slot_name,omitempty"`
	ApplicationName     string `json:"application_name,omitempty"`
	PromotionPriority   int    `json:"promotion_priority"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NodeTopologyRepository struct {
	conn *sql.DB
}

func NewNodeTopologyRepo(conn *sql.DB) *NodeTopologyRepository {
	return &NodeTopologyRepository{
		conn: conn,
	}
}

func (r *NodeTopologyRepository) TopologyAdd(
	ctx context.Context,
	projectID string,
	primaryNodeID string,
	replicaID string,
) (NodeTopology, error) {

	query := `
		INSERT INTO node_topology
		(
			relation_id,
			project_id,
			shard_id,
			shard_node_id,
			replica_node_id,
			created_at,
			updated_at
		)
		SELECT $1, $2, shard.id, $3, $4, $5, $5
		FROM shards AS shard
		WHERE shard.project_id = $2 AND shard.primary_node_id = $3
		RETURNING shard_id, relationship_status, updated_at
	`

	newUUID := uuid.New().String()
	currTime := time.Now()

	var logicalShardID, relationshipStatus string
	var updatedAt time.Time
	err := r.conn.QueryRowContext(
		ctx,
		query,
		newUUID,
		projectID,
		primaryNodeID,
		replicaID,
		currTime,
	).Scan(&logicalShardID, &relationshipStatus, &updatedAt)
	if err != nil {
		return NodeTopology{}, err
	}

	return NodeTopology{
		RelationID:    string(newUUID),
		ProjectID:     projectID,
		ShardID:       logicalShardID,
		PrimaryNodeID: primaryNodeID,
		StandbyNodeID: replicaID,
		Status:        relationshipStatus,
		CreatedAt:     currTime,
		UpdatedAt:     updatedAt,
	}, err

}

func (r *NodeTopologyRepository) TopologyRemove(
	ctx context.Context,
	relationID string,
) error {

	query := `
		DELETE FROM node_topology
		WHERE relation_id = $1
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		relationID,
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

func (r *NodeTopologyRepository) TopologyUpdate(
	ctx context.Context,
	relationID string,
	projectID string,
	primaryNodeID string,
	replicaID string,
) error {

	query := `
		UPDATE node_topology
		SET
			project_id = $2,
			shard_id = (SELECT id FROM shards
				WHERE project_id = $2 AND primary_node_id = $3),
			shard_node_id = $3,
			replica_node_id = $4,
			updated_at = NOW()
		WHERE relation_id = $1
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		relationID,
		projectID,
		primaryNodeID,
		replicaID,
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

func (r *NodeTopologyRepository) TopologyGetAll(
	ctx context.Context,
	projectID string,
) ([]NodeTopology, error) {

	query := `
		SELECT
			relation_id,
			project_id,
			shard_id,
			shard_node_id,
			replica_node_id,
			relationship_status,
			COALESCE(replication_slot_name, ''),
			COALESCE(application_name, ''),
			promotion_priority,
			created_at,
			updated_at
		FROM node_topology
		WHERE project_id = $1
	`

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topologies := make([]NodeTopology, 0)

	for rows.Next() {

		var topology NodeTopology

		err := rows.Scan(
			&topology.RelationID,
			&topology.ProjectID,
			&topology.ShardID,
			&topology.PrimaryNodeID,
			&topology.StandbyNodeID,
			&topology.Status,
			&topology.ReplicationSlotName,
			&topology.ApplicationName,
			&topology.PromotionPriority,
			&topology.CreatedAt,
			&topology.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		topologies = append(topologies, topology)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topologies, nil
}

func (r *NodeTopologyRepository) TopologyGetByID(
	ctx context.Context,
	relationID string,
) (NodeTopology, error) {

	query := `
		SELECT
			relation_id,
			project_id,
			shard_id,
			shard_node_id,
			replica_node_id,
			relationship_status,
			COALESCE(replication_slot_name, ''),
			COALESCE(application_name, ''),
			promotion_priority,
			created_at,
			updated_at
		FROM node_topology
		WHERE relation_id = $1
	`

	var topology NodeTopology

	err := r.conn.QueryRowContext(
		ctx,
		query,
		relationID,
	).Scan(
		&topology.RelationID,
		&topology.ProjectID,
		&topology.ShardID,
		&topology.PrimaryNodeID,
		&topology.StandbyNodeID,
		&topology.Status,
		&topology.ReplicationSlotName,
		&topology.ApplicationName,
		&topology.PromotionPriority,
		&topology.CreatedAt,
		&topology.UpdatedAt,
	)

	if err != nil {
		return NodeTopology{}, err
	}

	return topology, nil
}

func (r *NodeTopologyRepository) GetReplicasForShard(
	ctx context.Context,
	primaryNodeID string,
) ([]string, error) {

	query := `
		SELECT replica_node_id
		FROM node_topology
		WHERE shard_node_id = $1
	`

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		primaryNodeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	replicas := make([]string, 0)

	for rows.Next() {

		var replicaID string

		err := rows.Scan(&replicaID)
		if err != nil {
			return nil, err
		}

		replicas = append(replicas, replicaID)
	}

	return replicas, rows.Err()
}

func (r *NodeTopologyRepository) GetStandbysForLogicalShard(
	ctx context.Context,
	shardID string,
) ([]string, error) {
	query := `
		SELECT replica_node_id
		FROM node_topology
		WHERE shard_id = $1
		ORDER BY promotion_priority, created_at
	`
	rows, err := r.conn.QueryContext(ctx, query, shardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	standbys := make([]string, 0)
	for rows.Next() {
		var standbyID string
		if err := rows.Scan(&standbyID); err != nil {
			return nil, err
		}
		standbys = append(standbys, standbyID)
	}
	return standbys, rows.Err()
}

func (r *NodeTopologyRepository) GetAllMappings(
	ctx context.Context,
) ([]NodeTopology, error) {

	query := `
		SELECT
			relation_id,
			project_id,
			shard_id,
			shard_node_id,
			replica_node_id,
			relationship_status,
			COALESCE(replication_slot_name, ''),
			COALESCE(application_name, ''),
			promotion_priority,
			created_at,
			updated_at
		FROM node_topology
	`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mappings := make([]NodeTopology, 0)

	for rows.Next() {

		var topology NodeTopology

		err := rows.Scan(
			&topology.RelationID,
			&topology.ProjectID,
			&topology.ShardID,
			&topology.PrimaryNodeID,
			&topology.StandbyNodeID,
			&topology.Status,
			&topology.ReplicationSlotName,
			&topology.ApplicationName,
			&topology.PromotionPriority,
			&topology.CreatedAt,
			&topology.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		mappings = append(mappings, topology)
	}

	return mappings, rows.Err()
}

// func (r *NodeTopologyRepository) getAllReplicas(ctx context.Context) ([]string, error) {

// 	query := `
// 		SELECT replica_node_id FROM node_topology
// 	`
// 	results, err := r.conn.QueryContext(
// 		ctx,
// 		query,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var replicaNode string
// 	var replicaNodes []string

// 	for results.Next() {

// 		err = results.Scan(
// 			&replicaNode,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		replicaNodes = append(replicaNodes, replicaNode)

// 	}

// 	return replicaNodes, nil

// }
