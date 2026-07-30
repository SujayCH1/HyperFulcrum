package schema

import (
	"context"
	"hyperfulcrum/internal/repository"
)

type SchemaService struct {
	columnRepo  *repository.ColumnRepository
	fkEdgesRepo *repository.FKEdgesRepository
}

func NewSchemaService(columnRepo *repository.ColumnRepository, fkEdgesRepo *repository.FKEdgesRepository) *SchemaService {
	return &SchemaService{
		columnRepo:  columnRepo,
		fkEdgesRepo: fkEdgesRepo,
	}
}

func (s *SchemaService) ApplyDDL(ctx context.Context, projectID string,
	ddl string) error {

	//build schema from incoming ddl
	deltaSchema, err := BuildLogicalSchemaFromDDL(ddl)
	if err != nil {
		return err
	}
	//build schema from existing metadata stored
	columns, err := s.columnRepo.ColumnsListByProjectID(ctx, projectID)
	if err != nil {
		return err
	}
	fkEdges, err := s.fkEdgesRepo.EdgesListByProjectID(ctx, projectID)
	if err != nil {
		return err
	}

	baseSchema, err := BuildLogicalSchemaFromMetaData(projectID, columns, fkEdges)
	if err != nil {
		return err
	}

	//merge both schema
	mergedSchema, err := MergeLogicalSchema(baseSchema, deltaSchema)
	if err != nil {
		return err
	}

	newColumns, newFKEdges, err := FlattenLogicalSchema(mergedSchema)
	if err != nil {
		return err
	}

	//save to db
	if err := s.columnRepo.ColumnReplace(ctx, projectID, newColumns); err != nil {
		return err
	}

	if err := s.fkEdgesRepo.FKEdgesReplaceForProject(ctx, projectID, newFKEdges); err != nil {
		return err
	}
	return nil
}
