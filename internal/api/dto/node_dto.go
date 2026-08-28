package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NodeDto struct {
	Name   string `json:"name"`
	Index  int    `json:"index"`
	Status bool   `json:"status"`
	Role   string `json:"role"`
}

func (dto *NodeDto) Validate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.Index, validation.Skip),
		validation.Field(&dto.Status, validation.Skip),
		validation.Field(&dto.Role, validation.Required, validation.In("primary", "standby", "unassigned")),
	)
}

func ValidateNodeName(name string) error {
	return validation.Validate(name, validation.Required, validation.Length(3, 100))
}

func ValidateNodeRole(nodeRole string) error {
	return validation.Validate(nodeRole, validation.Required, validation.In("primary", "standby", "unassigned"))
}
