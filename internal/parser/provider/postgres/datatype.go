package postgres

import (
	"strconv"
	"strings"

	"hyperfulcrum/internal/parser/ir"
)

// parseDataType converts ValkDB's raw column type text (e.g. "varchar(255)",
// "numeric(10,2)", "text[]", "bigint") into the structured ir.DataType.
func parseDataType(raw string) ir.DataType {
	t := strings.TrimSpace(raw)

	dt := ir.DataType{}

	if strings.HasSuffix(t, "[]") {
		dt.Array = true
		t = strings.TrimSuffix(t, "[]")
	}

	name := t
	if open := strings.Index(t, "("); open != -1 && strings.HasSuffix(t, ")") {
		name = strings.TrimSpace(t[:open])
		args := t[open+1 : len(t)-1]
		parts := strings.Split(args, ",")

		switch len(parts) {
		case 1:
			if n, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				dt.Length = &n
			}
		case 2:
			if p, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				dt.Precision = &p
			}
			if s, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				dt.Scale = &s
			}
		}
	}

	dt.Name = strings.ToLower(name)
	return dt
}
