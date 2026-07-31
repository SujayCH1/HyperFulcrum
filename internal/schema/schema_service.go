package schema

import (
	"context"

	"hyperfulcrum/internal/metadata"
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

// main orchestraction function
func (s *SchemaService) ApplyProjectDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
	batch *ir.Batch,
) error {

	if err := s.StoreRawDDL(
		ctx,
		projectID,
		rawSQL,
	); err != nil {
		return err
	}

	if err := s.ConvertDDL(
		ctx,
		projectID,
		batch,
	); err != nil {
		return err
	}

	if err := s.ExecuteDDL(
		ctx,
		projectID,
		batch,
	); err != nil {
		return err
	}

	return nil
}

// ApplyDDL updates the project's logical schema metadata by replaying a batch
// of parsed DDL statements.
func (s *SchemaService) ConvertDDL(
	ctx context.Context,
	projectID string,
	batch *ir.Batch,
) error {

	// Load current metadata.
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

	// Build logical schema.
	schema, err := BuildLogicalSchemaFromMetadata(
		projectID,
		columns,
		fkEdges,
	)
	if err != nil {
		return err
	}

	// Replay the parsed DDL.
	if err := ApplyBatch(schema, batch); err != nil {
		return err
	}

	// Flatten back into metadata.
	updatedColumns, updatedEdges := FlattenLogicalSchema(schema)

	// Persist.
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

func (s *SchemaService) ExecuteDDL(
	ctx context.Context,
	projectID string,
	batch *ir.Batch,
) error {

	// TODO:
	// planner
	// executor
	// replication handling
	// status updates

	return nil
}
