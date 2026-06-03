package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type NodeConnectionDto struct {
	NodeID       string `json:"node_id"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

func (d *NodeConnectionDto) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.NodeID, validation.Required),
		validation.Field(&d.Host, validation.Required),
		validation.Field(&d.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&d.DatabaseName, validation.Required),
		validation.Field(&d.Username, validation.Required),
		validation.Field(&d.Password, validation.Required),
	)
}
