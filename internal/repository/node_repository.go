package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const (
	NodeRolePrimary    = "primary"
	NodeRoleStandby    = "standby"
	NodeRoleUnassigned = "unassigned"
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

func (r *NodeRepository) NodeList(
	ctx context.Context,
	projectID string,
) ([]Node, error) {

	query := `
		SELECT
			id,
			project_id,
			node_name,
			node_index,
			node_status,
			node_type,
			created_at
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

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *NodeRepository) NodeAdd(
	ctx context.Context,
	projectID string,
	nodeRole string,
	nodeName string,
) (Node, error) {
	if err := validateNodeRole(nodeRole); err != nil {
		return Node{}, err
	}

	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return Node{}, err
	}
	defer tx.Rollback()

	queryProject := `
		SELECT id
		FROM projects
		WHERE id = $1
		FOR UPDATE
	`

	var project string
	err = tx.QueryRowContext(ctx, queryProject, projectID).Scan(&project)
	if err != nil {
		return Node{}, err
	}

	var index int

	queryIndex := `
		SELECT COALESCE(MAX(node_index), -1) + 1
		FROM nodes
		WHERE project_id = $1
	`

	err = tx.QueryRowContext(ctx, queryIndex, projectID).Scan(&index)
	if err != nil {
		return Node{}, err
	}

	node := Node{
		ID:        uuid.NewString(),
		ProjectID: projectID,
		Name:      nodeName,
		Index:     index,
		Status:    false,
		Type:      nodeRole,
	}

	query := `
		INSERT INTO nodes 
		(id, project_id, node_name, node_index, node_status, node_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at
	`

	err = tx.QueryRowContext(
		ctx,
		query,
		node.ID,
		node.ProjectID,
		node.Name,
		node.Index,
		node.Status,
		node.Type,
	).Scan(&node.CreatedAt)
	if err != nil {
		return Node{}, err
	}

	err = tx.Commit()
	if err != nil {
		return Node{}, err
	}

	return node, nil
}

func (r *NodeRepository) NodeRemove(
	ctx context.Context,
	nodeID string,
) error {

	query := `
		DELETE FROM nodes
		WHERE id = $1
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		nodeID,
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

	res, err := r.conn.ExecContext(
		ctx,
		query,
		name,
		nodeID,
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

func (r *NodeRepository) NodeUpdateStatus(
	ctx context.Context,
	nodeID string,
	status bool,
) error {

	query := `
		UPDATE nodes
		SET node_status = $1
		WHERE id = $2
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		status,
		nodeID,
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

func (r *NodeRepository) NodeUpdateType(
	ctx context.Context,
	nodeID string,
	nodeRole string,
) error {
	return r.NodeUpdateRole(ctx, nodeID, nodeRole)
}

// NodeUpdateRole updates the desired PostgreSQL role of a physical node.
// NodeUpdateType remains as a compatibility wrapper while callers migrate to
// the new terminology.
func (r *NodeRepository) NodeUpdateRole(
	ctx context.Context,
	nodeID string,
	nodeRole string,
) error {
	if err := validateNodeRole(nodeRole); err != nil {
		return err
	}

	query := `
		UPDATE nodes
		SET node_type = $1
		WHERE id = $2
	`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		nodeRole,
		nodeID,
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

func validateNodeRole(nodeRole string) error {
	switch nodeRole {
	case NodeRolePrimary, NodeRoleStandby, NodeRoleUnassigned:
		return nil
	default:
		return fmt.Errorf("invalid node role: %s", nodeRole)
	}
}

// func (r *NodeRepository) NodesGetByPorjectID(ctx context.Context, projectID string) ([]Node, error) {

// 	query := `
// 		SELECT
// 		id, project_id, node_name, node_index, node_status, node_type, created_at
// 		FROM nodes
// 		WHERE project_id = $1
// 	`

// 	rows, err := r.conn.QueryContext(
// 		ctx,
// 		query,
// 		projectID,
// 	)
// 	if err != nil {
// 		return []Node{}, err
// 	}

// 	var nodes []Node
// 	var node Node

// 	for rows.Next() {
// 		err := rows.Scan(
// 			&node.ID,
// 			&node.ProjectID,
// 			&node.Name,
// 			&node.Index,
// 			&node.Status,
// 			&node.Type,
// 			&node.CreatedAt,
// 		)
// 		if err != nil {
// 			return []Node{}, err
// 		}

// 		nodes = append(nodes, node)
// 	}

// 	return nodes, nil

// }

func (r *NodeRepository) NodeGetByID(ctx context.Context, nodeID string) (Node, error) {

	query := `
		SELECT 
		id, project_id, node_name, node_index, node_status, node_type, created_at
		FROM nodes
		WHERE id = $1
	`

	row := r.conn.QueryRowContext(
		ctx,
		query,
		nodeID,
	)

	var node Node

	err := row.Scan(
		&node.ID,
		&node.ProjectID,
		&node.Name,
		&node.Index,
		&node.Status,
		&node.Type,
		&node.CreatedAt,
	)
	if err != nil {
		return Node{}, err
	}

	return node, nil

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
