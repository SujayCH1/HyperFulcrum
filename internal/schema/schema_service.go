package schema

import (
	"context"

	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/parser"
	"hyperfulcrum/internal/parser/ir"
)

type SchemaService struct {
	columnService        *metadata.ColumnService
	fkEdgeService        *metadata.FKEdgesService
	schemaVersionService *metadata.SchemaVersionService
}

func NewSchemaService(
	columnService *metadata.ColumnService,
	fkEdgeService *metadata.FKEdgesService,
	schemaVersionService *metadata.SchemaVersionService,
) *SchemaService {
	return &SchemaService{
		columnService:        columnService,
		fkEdgeService:        fkEdgeService,
		schemaVersionService: schemaVersionService,
	}
}

func (s *SchemaService) ApplyProjectDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
) error {

	batch, err := parser.ParseDDLBatch(rawSQL)
	if err != nil {
		return err
	}

	if err := s.ConvertDDL(
		ctx,
		projectID,
		batch,
	); err != nil {
		return err
	}

	if err := s.StoreRawDDL(
		ctx,
		projectID,
		rawSQL,
	); err != nil {
		return err
	}

	return nil
}

func (s *SchemaService) ConvertDDL(
	ctx context.Context,
	projectID string,
	batch *ir.Batch,
) error {

	columns, err := s.columnService.ListColumns(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	fkEdges, err := s.fkEdgeService.ListEdges(
		ctx,
		projectID,
	)
	if err != nil {
		return err
	}

	schema, err := BuildLogicalSchemaFromMetadata(
		projectID,
		columns,
		fkEdges,
	)
	if err != nil {
		return err
	}

	if err := ApplyBatch(schema, batch); err != nil {
		return err
	}

	updatedColumns, updatedEdges := FlattenLogicalSchema(schema)

	if err := s.columnService.ReplaceColumns(
		ctx,
		projectID,
		updatedColumns,
	); err != nil {
		return err
	}

	if err := s.fkEdgeService.ReplaceEdges(
		ctx,
		projectID,
		updatedEdges,
	); err != nil {
		return err
	}

	return nil
}

func (s *SchemaService) StoreRawDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
) error {

	return s.schemaVersionService.ReplaceSchemaVersion(
		ctx,
		projectID,
		rawSQL,
	)
}
