package shardkey

import (
	"context"
	"hyperfulcrum/internal/repository"
	"hyperfulcrum/pkg/logger"
)

type InferenceService struct {
	columnRepo   *repository.ColumnRepository
	fkRepo       *repository.FKEdgesRepository
	shardKeyRepo *repository.ShardKeysRepository
}

func NewInferenceService(
	columnRepo *repository.ColumnRepository,
	fkRepo *repository.FKEdgesRepository,
	shardKeyRepo *repository.ShardKeysRepository,
) *InferenceService {
	return &InferenceService{
		columnRepo:   columnRepo,
		fkRepo:       fkRepo,
		shardKeyRepo: shardKeyRepo,
	}
}

func (s *InferenceService) ApplyShardKeyInference(
	ctx context.Context,
	projectID string,
) error {
	return s.RunHeuristicInference(ctx, projectID)
}

func (s *InferenceService) RunHeuristicInference(ctx context.Context, projectID string) error {
	logger.Logger.Info("heuristic inference started")

	logicalSchema, err := s.buildSchema(ctx, projectID)
	if err != nil {
		return err
	}

	inferenceResult := BuildShardKeyPlan(&logicalSchema)
	inferred := convertDecisionsToShardKeyRecords(inferenceResult.Decisions)

	return s.shardKeyRepo.ReplaceShardKeysForProject(ctx, projectID, inferred)
}

func convertDecisionsToShardKeyRecords(
	decisions []ShardKeyDecision,
) []repository.ShardKeyRecord {
	records := make([]repository.ShardKeyRecord, 0, len(decisions))

	for _, d := range decisions {
		records = append(records, repository.ShardKeyRecord{
			TableName:      d.Table,
			ShardKeyColumn: d.Column.Column,
			IsManual:       false,
		})
	}

	return records
}
