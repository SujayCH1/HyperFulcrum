package schema

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"hyperfulcrum/internal/metadata"
	"hyperfulcrum/internal/parser"
	"hyperfulcrum/internal/repository"
)

type SchemaMetadata struct {
	RawSQL  string
	Columns []repository.Column
	FKEdges []repository.FkEdges
}

type SchemaService struct {
	schemaVersionService *metadata.SchemaVersionService
}

func NewSchemaService(
	schemaVersionService *metadata.SchemaVersionService,
) *SchemaService {
	return &SchemaService{
		schemaVersionService: schemaVersionService,
	}
}

// ApplyProjectDDL appends a DDL fragment to the editable project draft.
func (s *SchemaService) ApplyProjectDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
) error {
	_, err := s.AppendProjectDDL(ctx, projectID, rawSQL)
	return err
}

func (s *SchemaService) AppendProjectDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
) (repository.SchemaVersion, error) {
	currentSQL := ""
	expectedRevision := int64(0)

	current, err := s.schemaVersionService.GetSchemaVersion(ctx, projectID)
	if err == nil {
		currentSQL = current.RawSQL
		expectedRevision = current.Revision
	} else if !errors.Is(err, sql.ErrNoRows) {
		return repository.SchemaVersion{}, err
	}

	combinedSQL := appendDDL(currentSQL, rawSQL)

	return s.replaceProjectSchema(
		ctx,
		projectID,
		combinedSQL,
		expectedRevision,
	)
}

func (s *SchemaService) ReplaceProjectDDL(
	ctx context.Context,
	projectID string,
	rawSQL string,
) (repository.SchemaVersion, error) {
	expectedRevision := int64(0)
	current, err := s.schemaVersionService.GetSchemaVersion(ctx, projectID)
	if err == nil {
		expectedRevision = current.Revision
	} else if !errors.Is(err, sql.ErrNoRows) {
		return repository.SchemaVersion{}, err
	}

	return s.replaceProjectSchema(
		ctx,
		projectID,
		strings.TrimSpace(rawSQL),
		expectedRevision,
	)
}

func (s *SchemaService) LockProjectSchema(
	ctx context.Context,
	projectID string,
) (repository.SchemaVersion, error) {
	return s.schemaVersionService.LockSchema(ctx, projectID)
}

func (s *SchemaService) UnlockProjectSchema(
	ctx context.Context,
	projectID string,
) (repository.SchemaVersion, error) {
	return s.schemaVersionService.UnlockSchema(ctx, projectID)
}

func (s *SchemaService) replaceProjectSchema(
	ctx context.Context,
	projectID string,
	rawSQL string,
	expectedRevision int64,
) (repository.SchemaVersion, error) {
	metadataSnapshot, err := BuildSchemaMetadata(projectID, rawSQL)
	if err != nil {
		return repository.SchemaVersion{}, err
	}

	return s.schemaVersionService.ReplaceProjectSchema(
		ctx,
		projectID,
		metadataSnapshot.RawSQL,
		metadataSnapshot.Columns,
		metadataSnapshot.FKEdges,
		expectedRevision,
	)
}

func BuildSchemaMetadata(
	projectID string,
	rawSQL string,
) (SchemaMetadata, error) {
	batch, err := parser.ParseDDLBatch(rawSQL)
	if err != nil {
		return SchemaMetadata{}, err
	}

	logicalSchema := NewLogicalSchema()
	logicalSchema.ProjectID = projectID

	err = ApplyBatch(logicalSchema, batch)
	if err != nil {
		return SchemaMetadata{}, err
	}

	columns, edges := FlattenLogicalSchema(logicalSchema)

	return SchemaMetadata{
		RawSQL:  rawSQL,
		Columns: columns,
		FKEdges: edges,
	}, nil
}

func appendDDL(currentSQL string, rawSQL string) string {
	currentSQL = strings.TrimSpace(currentSQL)
	rawSQL = strings.TrimSpace(rawSQL)

	if currentSQL == "" {
		return rawSQL
	}
	if rawSQL == "" {
		return currentSQL
	}
	if !strings.HasSuffix(currentSQL, ";") {
		currentSQL += ";"
	}

	return currentSQL + "\n" + rawSQL
}
