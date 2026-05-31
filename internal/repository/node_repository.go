package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Node struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"node_name"`
	Index     int    `json:"node_index"`
	Status    bool   `json:"node_status"`
	Type      string `json:"node_type"`
	CreatedAt string `json:"created_at"`
}

type NodeRepository struct {
	conn *sql.DB
}

func NewNodeRepository(connConfig *sql.DB) *NodeRepository {
	return &NodeRepository{conn: connConfig}
}

// Main functions

func (r *NodeRepository) NodeList(ctx context.Context, projectID string) ([]Node, error) {

	query := `
		SELECT
		id, project_id, node_name, node_index, node_status, node_type, created_at
		FROM nodes
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

	nodes := make([]Node, 0)

	for rows.Next() {
		var node Node
		err = rows.Scan(
			&node.ID,
			&node.ProjectID,
			&node.Name,
			&node.Index,
			&node.Status,
			&node.Type,
			&node.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	return nodes, nil

}

func (r *NodeRepository) NodeAdd(ctx context.Context, projectID string, nodeType string, nodeName string) error {
	var index int

	queryIndex := `
		SELECT COALESCE(MAX(node_index), -1) + 1
		FROM nodes
		WHERE project_id = $1
	`

	err := r.conn.QueryRowContext(ctx, queryIndex, projectID).Scan(&index)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO nodes 
		(id, project_id, node_name, node_index, node_status, node_type)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = r.conn.ExecContext(
		ctx,
		query,
		uuid.New(),
		projectID,
		nodeName,
		index,
		false,
		nodeType,
	)
	if err != nil {
		return err
	}

	return nil

}

func (r *NodeRepository) NodeRemove(ctx context.Context, nodeID string) error {

	query := `
		DELETE FROM nodes
		WHERE id = $1
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		nodeID,
	)
	if err != nil {
		return err
	}

	return nil

}

func (r *NodeRepository) NodeRemoveAll(ctx context.Context, projectID string) error {

	query := `
		DELETE FROM nodes
		WHERE project_id = $1
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return err
	}

	return nil

}

func (r *NodeRepository) NodeUpdateName(ctx context.Context, nodeID string, name string) error {
	query := `
		UPDATE nodes
		SET node_name = $1
		WHERE id = $2
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		name,
		nodeID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *NodeRepository) NodeUpdateStatus(ctx context.Context, nodeID string, status bool) error {
	query := `
		UPDATE nodes
		SET node_status = $1
		WHERE id = $2
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		status,
		nodeID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *NodeRepository) UpdateType(ctx context.Context, nodeID string, nodeType string) error {
	if nodeType != "shard" && nodeType != "replica" {
		return fmt.Errorf("invalid node type: %s", nodeType)
	}

	query := `
		UPDATE nodes
		SET node_type = $1
		WHERE id = $2
	`

	_, err := r.conn.ExecContext(
		ctx,
		query,
		nodeType,
		nodeID,
	)
	if err != nil {
		return err
	}

	return nil
}

// helpers

// func (r *NodeRepository) fetchNodeIndexes(ctx context.Context, projectID string) ([]int, error) {

// 	query := `
// 		SELECT
// 		node_index
// 		FROM nodes
// 		WHERE project_id = $1
// 	`

// 	rows, err := r.conn.QueryContext(
// 		ctx,
// 		query,
// 		projectID,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	indexes := make([]int, 0)

// 	for rows.Next() {
// 		var index int
// 		err = rows.Scan(
// 			&index,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}

// 		indexes = append(indexes, index)
// 	}

// 	return indexes, nil
// }
