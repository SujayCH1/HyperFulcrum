package repository

import (
	"context"
	"database/sql"
	"time"
)

type Column struct {
	ProjectID  string `json:"project_id"`
	TableName  string `json:"table_name"`
	ColumnName string `json:"column_name"`
	DataType   string `json:"data_type"`
	IsNullable bool   `json:"is_nullable"`
	IsPk       bool   `json:"is_pk"`

	CreatedAt time.Time `json:"created_at"`
}

type ColumnRepository struct {
	conn *sql.DB
}

func NewColumnRepository(connConfig *sql.DB) *ColumnRepository {
	return &ColumnRepository{conn: connConfig}
}

func (r *ColumnRepository) ColumnAdd(ctx context.Context, column *Column) error {
	query := `
		INSERT INTO columns (project_id,table_name,column_name,data_type,nullable,is_primary_key)
		VALUES ($1,$2,$3,$4,$5,$6)
	`

	_, err := r.conn.ExecContext(ctx,
		query, column.ProjectID, column.TableName,
		column.ColumnName, column.DataType, column.IsNullable, column.IsPk)

	if err != nil {
		return err
	}
	return nil
}

func (r *ColumnRepository) ColumnReplace(ctx context.Context, projectID string,
	cols []Column) error {
	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `DELETE FROM columns WHERE project_id = $1`
	_, err = tx.ExecContext(ctx, query, projectID)
	if err != nil {
		return err
	}

	for _, col := range cols {
		query := `
			INSERT INTO columns (project_id,table_name,column_name,data_type,nullable,is_primary_key)
			VALUES ($1,$2,$3,$4,$5,$6)
		`
		_, err = tx.ExecContext(ctx, query, projectID, col.TableName, col.ColumnName, col.DataType, col.IsNullable, col.IsPk)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ColumnRepository) ColumnsListByProjectID(ctx context.Context, projectID string) ([]Column, error) {
	query := `SELECT project_id,table_name,column_name,data_type,nullable,is_primary_key
	FROM columns WHERE project_id = $1`

	rows, err := r.conn.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]Column, 0)
	for rows.Next() {
		var column Column
		err := rows.Scan(&column.ProjectID, &column.TableName, &column.ColumnName, &column.DataType, &column.IsNullable, &column.IsPk)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}
