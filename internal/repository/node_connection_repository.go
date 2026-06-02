package repository

import (
	"context"
	"database/sql"
	"time"
)

type NodeConnection struct {
	NodeId       string `json:"node_id"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NodeConnectionRepository struct {
	conn *sql.DB
}

func NewNodeConnectionRepository(connConfig *sql.DB) *NodeConnectionRepository {
	return &NodeConnectionRepository{conn: connConfig}
}

func (r *NodeConnectionRepository) ConnectionCreate(ctx context.Context, node_conn *NodeConnection) error {

	query := `INSERT INTO node_connections (node_id,host,port,database_name,username,password) VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := r.conn.ExecContext(ctx, query, node_conn.NodeId, node_conn.Host, node_conn.Port, node_conn.DatabaseName, node_conn.Username, node_conn.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *NodeConnectionRepository) ConnectionRemove(ctx context.Context, nodeId string) error {

	query := `DELETE FROM node_connections WHERE node_id = $1`
	res, err := r.conn.ExecContext(ctx, query, nodeId)
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

func (r *NodeConnectionRepository) ConnectionUpdate(ctx context.Context, nodeConn *NodeConnection) error {

	query := `UPDATE node_connections SET host = $1, port = $2, database_name = $3, username = $4, password = $5, updated_at = NOW()  WHERE node_id = $6`
	res, err := r.conn.ExecContext(ctx, query, nodeConn.Host, nodeConn.Port, nodeConn.DatabaseName, nodeConn.Username, nodeConn.Password, nodeConn.NodeId)
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

func (r *NodeConnectionRepository) GetConnectionByNodeId(ctx context.Context, nodeId string) (*NodeConnection, error) {

	query := `SELECT node_id, host, port, database_name, username, password, created_at, updated_at FROM node_connections WHERE node_id = $1`
	row := r.conn.QueryRowContext(ctx, query, nodeId)

	var nodeConn NodeConnection
	err := row.Scan(&nodeConn.NodeId, &nodeConn.Host, &nodeConn.Port, &nodeConn.DatabaseName, &nodeConn.Username, &nodeConn.Password, &nodeConn.CreatedAt, &nodeConn.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &nodeConn, nil
}
