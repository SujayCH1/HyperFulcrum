package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type TopologyCreateDto struct {
	ProjectID     string `json:"project_id"`
	ShardNodeID   string `json:"shard_node_id"`
	ReplicaNodeID string `json:"replica_node_id"`
}

type TopologyDeleteDto struct {
	RelationID string `json:"relation_id"`
	ProjectID  string `json:"project_id"`
}

func (dto *TopologyCreateDto) ValidateCreate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.ProjectID, validation.Required, validUUID),
		validation.Field(&dto.ShardNodeID, validation.Required, validUUID),
		validation.Field(&dto.ReplicaNodeID, validation.Required, validUUID),
	)
}

func (dto *TopologyDeleteDto) ValidateDelete() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.RelationID, validation.Required, validUUID),
		validation.Field(&dto.ProjectID, validation.Required, validUUID),
	)
}
