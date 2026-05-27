package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	ShardCount  int    `json:"shard_count"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) ProjectAdd(
	ctx context.Context, name string, description string) (*Project, error) {

	project := &Project{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		ShardCount:  0,
	}
	query := `INSERT INTO projects (id,name,description,shard_count) VALUES ($1,$2,$3,$4) RETURNING status, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, project.ID, project.Name, project.Description, project.ShardCount).Scan(&project.Status, &project.CreatedAt, &project.UpdatedAt)
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
	result, err := r.db.ExecContext(ctx, query, projectID)
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

	query := `SELECT id,name,description,shard_count,status,created_at,updated_at FROM projects ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		project := &Project{}
		err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.ShardCount, &project.Status, &project.CreatedAt, &project.UpdatedAt)
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

func (r *ProjectRepository) GetProjectByID(ctx context.Context, id string) (Project, error) {
	projectID, err := uuid.Parse(id)
	if err != nil {
		return Project{}, err
	}

	query := `SELECT id,name,description,shard_count,status,created_at,updated_at FROM projects WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, projectID)

	project := &Project{}
	err = row.Scan(&project.ID, &project.Name, &project.Description, &project.ShardCount, &project.Status, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return Project{}, err
	}
	return *project, nil
}
