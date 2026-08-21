package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"hyperfulcrum/internal/repository"
)

type ShardKeyCreateDto struct {
	TableName string `json:"table_name"`
	KeyColumn string `json:"key_column"`
}

func (dto *ShardKeyCreateDto) Validate() error {
	return validation.ValidateStruct(
		dto,
		validation.Field(
			&dto.TableName,
			validation.Required,
			validation.Length(1, 255),
		),
		validation.Field(
			&dto.KeyColumn,
			validation.Required,
			validation.Length(1, 255),
		),
	)
}

type ShardKeyResponseDto struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	TableName string    `json:"table_name"`
	KeyColumn string    `json:"key_column"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewShardKeyResponse(key repository.ShardKey) ShardKeyResponseDto {
	return ShardKeyResponseDto{
		ID:        key.ID,
		ProjectID: key.ProjectID,
		TableName: key.TableName,
		KeyColumn: key.KeyColumn,
		UpdatedAt: key.UpdatedAt,
	}
}

func NewShardKeyListResponse(keys []repository.ShardKey) []ShardKeyResponseDto {
	response := make([]ShardKeyResponseDto, 0, len(keys))

	for _, key := range keys {
		response = append(response, NewShardKeyResponse(key))
	}

	return response
}
