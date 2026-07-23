package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type NodeTopology struct {
	RelationID    string `json:"relation_id"`
	ProjectID     string `json:"project_id"`
	ShardNodeID   string `json:"shard_node_id"`
	ReplicaNodeID string `json:"replica_node_id"`

	CreatedAt time.Time `json:"created_at"`
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
	shardID string,
	replicaID string,
) (NodeTopology, error) {

	query := `
		INSERT INTO node_topology
		(
			relation_id,
			project_id,
			shard_node_id,
			replica_node_id,
			created_at
		)
		VALUES
		($1, $2, $3, $4, $5)

	`

	newUUID := uuid.New().String()
	currTime := time.Now()

	_, err := r.conn.ExecContext(
		ctx,
		query,
		newUUID,
		projectID,
		shardID,
		replicaID,
		currTime,
	)
	if err != nil {
		return NodeTopology{}, err
	}

	return NodeTopology{
		RelationID:    string(newUUID),
		ProjectID:     projectID,
		ShardNodeID:   shardID,
		ReplicaNodeID: replicaID,
		CreatedAt:     currTime,
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

	_, err := r.conn.ExecContext(
		ctx,
		query,
		relationID,
	)

	return err
}

func (r *NodeTopologyRepository) TopologyUpdate(
	ctx context.Context,
	relationID string,
	projectID string,
	shardID string,
	replicaID string,
) error {

	query := `
		UPDATE node_topology
		SET
			project_id = $2,
			shard_node_id = $3,
			replica_node_id = $4
		WHERE relation_id = $1
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		relationID,
		projectID,
		shardID,
		replicaID,
	)

	return err
}

func (r *NodeTopologyRepository) TopologyGetAll(
	ctx context.Context,
	projectID string,
) ([]NodeTopology, error) {

	query := `
		SELECT
			relation_id,
			project_id,
			shard_node_id,
			replica_node_id,
			created_at
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

	var topologies []NodeTopology

	for rows.Next() {

		var topology NodeTopology

		err := rows.Scan(
			&topology.RelationID,
			&topology.ProjectID,
			&topology.ShardNodeID,
			&topology.ReplicaNodeID,
			&topology.CreatedAt,
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
			shard_node_id,
			replica_node_id,
			created_at
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
		&topology.ShardNodeID,
		&topology.ReplicaNodeID,
		&topology.CreatedAt,
	)

	if err != nil {
		return NodeTopology{}, err
	}

	return topology, nil
}

func (r *NodeTopologyRepository) GetReplicasForShard(
	ctx context.Context,
	shardID string,
) ([]string, error) {

	query := `
		SELECT replica_node_id
		FROM node_topology
		WHERE shard_node_id = $1
	`

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		shardID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replicas []string

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

func (r *NodeTopologyRepository) GetAllMappings(
	ctx context.Context,
) ([]NodeTopology, error) {

	query := `
		SELECT
			relation_id,
			project_id,
			shard_node_id,
			replica_node_id,
			created_at
		FROM node_topology
	`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mappings []NodeTopology

	for rows.Next() {

		var topology NodeTopology

		err := rows.Scan(
			&topology.RelationID,
			&topology.ProjectID,
			&topology.ShardNodeID,
			&topology.ReplicaNodeID,
			&topology.CreatedAt,
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
