package connections

import (
	"net"
	"net/url"
	"strconv"
	"strings"

	"hyperfulcrum/internal/repository"
)

func buildDSN(conn repository.NodeConnection) string {
	host := conn.Host
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	}

	dsn := &url.URL{
		Scheme:  "postgres",
		User:    url.UserPassword(conn.Username, conn.Password),
		Host:    net.JoinHostPort(host, strconv.Itoa(conn.Port)),
		Path:    "/" + conn.DatabaseName,
		RawPath: "/" + url.PathEscape(conn.DatabaseName),
	}

	query := dsn.Query()
	query.Set("sslmode", "disable")
	dsn.RawQuery = query.Encode()

	return dsn.String()
}
