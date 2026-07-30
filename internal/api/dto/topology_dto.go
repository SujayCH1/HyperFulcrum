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
		validation.Field(&dto.ProjectID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ShardNodeID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ReplicaNodeID, validation.Required, validation.Length(3, 100)),
	)
}

func (dto *TopologyDeleteDto) ValidateDelete() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.RelationID, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.ProjectID, validation.Required, validation.Length(3, 100)),
	)
}
