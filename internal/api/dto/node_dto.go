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
		validation.Field(&dto.Index, validation.Skip),
		validation.Field(&dto.Status, validation.Skip),
		validation.Field(&dto.Type, validation.Required, validation.In("shard", "replica")),
	)
}

func ValidateNodeName(name string) error {
	return validation.Validate(name, validation.Required, validation.Length(3, 100))
}

func ValidateNodeType(nodeType string) error {
	return validation.Validate(nodeType, validation.Required, validation.In("shard", "replica"))
}
