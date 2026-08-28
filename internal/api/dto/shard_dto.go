package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateShardDto struct {
	Name          string `json:"name"`
	PrimaryNodeID string `json:"primary_node_id"`
}

func (dto *CreateShardDto) Validate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.PrimaryNodeID, validation.Required, validUUID),
	)
}

func ValidateShardName(name string) error {
	return validation.Validate(name, validation.Required, validation.Length(3, 100))
}
