package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateReplicationDto struct {
	ProjectID     string `json:"project_id"`
	ShardNodeID   string `json:"shard_node_id"`
	ReplicaNodeID string `json:"replica_node_id"`
}

type DeleteReplicationDto struct {
	ProjectID  string `json:"project_id"`
	RelationID string `json:"relation_id"`
}

type PromoteReplicaDto struct {
	RelationID    string `json:"relation_id"`
	ShardNodeID   string `json:"shard_node_id"`
	ReplicaNodeID string `json:"replica_node_id"`
}

func (dto *CreateReplicationDto) ValidateCreate() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.ProjectID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ShardNodeID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ReplicaNodeID, validation.Required, validation.Length(3, 100)),
	)
}

func (dto *DeleteReplicationDto) ValidateDelete() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.ProjectID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.RelationID, validation.Required, validation.Length(3, 100)),
	)
}

func (dto *PromoteReplicaDto) ValidatePromote() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(&dto.RelationID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ShardNodeID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ReplicaNodeID, validation.Required, validation.Length(3, 100)),
	)
}
