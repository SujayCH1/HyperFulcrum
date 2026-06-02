package connections

import (
	"fmt"
	"hyperfulcrum/internal/repository"
)

func buildDSN(conn repository.NodeConnection) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conn.Username,
		conn.Password,
		conn.Host,
		conn.Port,
		conn.DatabaseName,
	)
}
