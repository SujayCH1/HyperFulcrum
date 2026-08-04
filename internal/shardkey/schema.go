package shardkey

import (
	"context"
	"hyperfulcrum/internal/schema"
)


func (s *InferenceService) buildSchema(ctx context.Context, projectID string) (schema.LogicalSchema, error) {
	columns, err := s.columnRepo.ColumnsListByProjectID(ctx, projectID)
	if err != nil {
		return schema.LogicalSchema{}, err
	}

	fkEdges, err := s.fkRepo.EdgesListByProjectID(ctx, projectID)
	if err != nil {
		return schema.LogicalSchema{}, err
	}

	logicalSchema, err := schema.BuildLogicalSchemaFromMetadata(projectID, columns, fkEdges)
	if err != nil {
		return schema.LogicalSchema{}, err
	}

	return *logicalSchema, nil
}
