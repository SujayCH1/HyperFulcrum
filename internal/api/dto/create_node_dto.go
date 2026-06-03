package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NodeDto struct {
	Name   string `json:"name"`
	Index  int    `json:"index"`
	Status bool   `json:"status"`
	Type   string `json:"type"`
}

func (dto *NodeDto) Validate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.Index, validation.Required),
		validation.Field(&dto.Status, validation.Required),
		validation.Field(&dto.Type, validation.Required, validation.In("shard", "replica")),
	)
}
