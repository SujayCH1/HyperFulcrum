package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Ready       bool   `json:"ready"`
	Running     bool   `json:"running"`
	NodeCount   int    `json:"node_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectRepository struct {
	conn *sql.DB
}

func NewProjectRepository(connConfig *sql.DB) *ProjectRepository {
	return &ProjectRepository{conn: connConfig}
}

func (r *ProjectRepository) ProjectAdd(
	ctx context.Context, name string, description string) (*Project, error) {

	project := &Project{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		NodeCount:   0,
		Ready:       false,
		Running:     false,
	}

	query := `
		INSERT INTO projects (id, name, description, node_count)
		VALUES ($1, $2, $3, $4)
		RETURNING ready, running, created_at, updated_at
	`

	err := r.conn.QueryRowContext(
		ctx,
		query,
		project.ID,
		project.Name,
		project.Description,
		project.NodeCount,
	).Scan(
		&project.Ready,
		&project.Running,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) ProjectRemove(ctx context.Context, id string) error {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM projects WHERE id = $1`
	result, err := r.conn.ExecContext(ctx, query, projectID)
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

func (r *ProjectRepository) ProjectList(ctx context.Context) ([]*Project, error) {
	query := `
		SELECT id, name, description, node_count, ready, running, created_at, updated_at
		FROM projects
		ORDER BY created_at DESC
	`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.NodeCount,
			&project.Ready,
			&project.Running,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepository) ProjectGetByID(ctx context.Context, id string) (*Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, name, description, node_count, ready, running, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	row := r.conn.QueryRowContext(ctx, query, projectID)

	project := &Project{}
	err = row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.NodeCount,
		&project.Ready,
		&project.Running,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) ProjectUpdateReady(ctx context.Context, id string, ready bool) (*Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	project := &Project{}

	query := `
		UPDATE projects
		SET ready = $1
		WHERE id = $2
		RETURNING id, name, description, node_count, ready, running, created_at, updated_at
	`

	err = r.conn.QueryRowContext(ctx, query, ready, projectID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.NodeCount,
		&project.Ready,
		&project.Running,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) ProjectUpdateRunning(ctx context.Context, id string, running bool) (*Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	project := &Project{}

	query := `
		UPDATE projects
		SET running = $1
		WHERE id = $2
		RETURNING id, name, description, node_count, ready, running, created_at, updated_at
	`

	err = r.conn.QueryRowContext(ctx, query, running, projectID).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.NodeCount,
		&project.Ready,
		&project.Running,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return project, nil
}


func (r *ProjectRepository) ProjectGetReady(ctx context.Context) ([]*Project, error) {
	query := `SELECT id,name,description,node_count,ready,running,created_at,updated_at 
	FROM projects WHERE ready=TRUE`

	rows, err := r.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(&project.ID,
			&project.Name,
			&project.Description,
			&project.NodeCount,
			&project.Ready,
			&project.Running,
			&project.CreatedAt,
			&project.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, sql.ErrNoRows
	}
	return projects, nil
}
