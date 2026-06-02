package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CreateProjectDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (dto *CreateProjectDto) Validate() error {
	return validation.ValidateStruct(dto,
		validation.Field(&dto.Name, validation.Required, validation.Length(3, 100)),
		validation.Field(&dto.Description, validation.Required, validation.Length(5, 500)),
	)
}
