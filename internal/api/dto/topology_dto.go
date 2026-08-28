package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type TopologyCreateDto struct {
	ShardID       string `json:"shard_id"`
	StandbyNodeID string `json:"standby_node_id"`
}

type TopologyDeleteDto struct {
	RelationID string `json:"relation_id"`
	ProjectID  string `json:"project_id"`
}

func (dto *TopologyCreateDto) ValidateCreate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.ShardID, validation.Required, validUUID),
		validation.Field(&dto.StandbyNodeID, validation.Required, validUUID),
	)
}

func (dto *TopologyDeleteDto) ValidateDelete() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.RelationID, validation.Required, validUUID),
		validation.Field(&dto.ProjectID, validation.Required, validUUID),
	)
}
