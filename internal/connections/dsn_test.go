package connections

import (
	"net/url"
	"testing"

	"hyperfulcrum/internal/repository"
)

func TestBuildDSNEscapesConnectionValues(t *testing.T) {
	connection := repository.NodeConnection{
		Host:         "[::1]",
		Port:         5432,
		DatabaseName: "database name",
		Username:     "user@example.com",
		Password:     "p@ss:/?#word",
	}

	dsn, err := url.Parse(buildDSN(connection))
	if err != nil {
		t.Fatal(err)
	}

	password, _ := dsn.User.Password()

	if dsn.User.Username() != connection.Username || password != connection.Password {
		t.Fatal("connection credentials were not preserved")
	}

	if dsn.Hostname() != "::1" || dsn.Port() != "5432" {
		t.Fatal("connection host was not encoded correctly")
	}

	if dsn.Path != "/"+connection.DatabaseName {
		t.Fatal("database name was not encoded correctly")
	}

	if dsn.Query().Get("sslmode") != "disable" {
		t.Fatal("ssl mode was not included")
	}
}
