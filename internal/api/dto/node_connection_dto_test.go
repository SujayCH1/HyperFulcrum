package dto

import (
	"encoding/json"
	"strings"
	"testing"

	"hyperfulcrum/internal/repository"
)

func TestNodeConnectionResponseDoesNotIncludePassword(t *testing.T) {
	connection := repository.NodeConnection{
		NodeId:       "node-id",
		Host:         "localhost",
		Port:         5432,
		DatabaseName: "postgres",
		Username:     "postgres",
		Password:     "secret",
	}

	data, err := json.Marshal(NewNodeConnectionResponse(connection))
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(string(data), "secret") || strings.Contains(string(data), "password") {
		t.Fatal("node connection response includes the password")
	}
}
