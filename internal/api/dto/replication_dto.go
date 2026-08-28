package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateReplicationDto struct {
	ProjectID     string `json:"project_id"`
	ShardID       string `json:"shard_id"`
	StandbyNodeID string `json:"standby_node_id"`
}

type DeleteReplicationDto struct {
	ProjectID  string `json:"project_id"`
	RelationID string `json:"relation_id"`
}

type PromoteReplicaDto struct {
	RelationID    string `json:"relation_id"`
	ShardID       string `json:"shard_id"`
	StandbyNodeID string `json:"standby_node_id"`
}

func (dto *CreateReplicationDto) ValidateCreate() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.ProjectID, validation.Required, validUUID),
		validation.Field(&dto.ShardID, validation.Required, validUUID),
		validation.Field(&dto.StandbyNodeID, validation.Required, validUUID),
	)
}

func (dto *DeleteReplicationDto) ValidateDelete() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.ProjectID, validation.Required, validUUID),
		validation.Field(&dto.RelationID, validation.Required, validUUID),
	)
}

func (dto *PromoteReplicaDto) ValidatePromote() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.RelationID, validation.Required, validUUID),
		validation.Field(&dto.ShardID, validation.Required, validUUID),
		validation.Field(&dto.StandbyNodeID, validation.Required, validUUID),
	)
}
