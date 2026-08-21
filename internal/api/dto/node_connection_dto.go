package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"hyperfulcrum/internal/repository"
)

type NodeConnectionDto struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func (d *NodeConnectionDto) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Host, validation.Required),
		validation.Field(&d.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&d.DatabaseName, validation.Required),
		validation.Field(&d.Username, validation.Required),
		validation.Field(&d.Password, validation.Required),
	)
}

type NodeConnectionResponseDto struct {
	NodeID       string     `json:"node_id"`
	Host         string     `json:"host"`
	Port         int        `json:"port"`
	DatabaseName string     `json:"database_name"`
	Username     string     `json:"username"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func NewNodeConnectionResponse(connection repository.NodeConnection) NodeConnectionResponseDto {
	response := NodeConnectionResponseDto{
		NodeID:       connection.NodeId,
		Host:         connection.Host,
		Port:         connection.Port,
		DatabaseName: connection.DatabaseName,
		Username:     connection.Username,
	}

	if !connection.CreatedAt.IsZero() {
		response.CreatedAt = &connection.CreatedAt
	}

	if !connection.UpdatedAt.IsZero() {
		response.UpdatedAt = &connection.UpdatedAt
	}

	return response
}
