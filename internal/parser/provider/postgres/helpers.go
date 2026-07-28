package postgres

import (
	"strconv"

	"hyperfulcrum/internal/parser/ir"
)

func ptr[T any](v T) *T {
	return &v
}

func newIdentifier(schema, name, alias string) ir.Identifier {
	return ir.Identifier{
		Schema: schema,
		Name:   name,
		Alias:  alias,
	}
}

func newTable(schema, name, alias string) ir.Table {
	return ir.Table{
		Schema: schema,
		Name:   name,
		Alias:  alias,
	}
}

// nonEmptyExpr returns nil when raw is empty, so DefaultValue stays nil
// instead of pointing at an empty Expression.
func nonEmptyExpr(raw string) *ir.Expression {
	if raw == "" {
		return nil
	}
	return &ir.Expression{Raw: raw}
}

// containsFlag checks DDLAction.Flags (e.g. "IF_EXISTS", "CONCURRENTLY", "UNIQUE").
func containsFlag(flags []string, needle string) bool {
	for _, f := range flags {
		if f == needle {
			return true
		}
	}
	return false
}

func firstOrEmpty(vals []string) string {
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// parseIntPtr parses a LIMIT/OFFSET value. ValkDB keeps these as raw strings
// since they may be placeholders like "$1"; ok is false when it isn't a
// plain integer literal.
func parseIntPtr(raw string) (*int, bool) {
	if raw == "" {
		return nil, false
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil, false
	}
	return &n, true
}
