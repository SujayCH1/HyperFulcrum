package shardkey

import (
	"strings"
	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/schema"
)



func ExtractCandidates(s *schema.LogicalSchema) CandidateSet {
	candidates := make(CandidateSet)

	for tableName, table := range s.Tables {
		var tableCandidates []ColumnRef

		for _, column := range table.Columns {
			eliminated, _ := isEliminated(column)
			if eliminated {
				continue
			}
			tableCandidates = append(tableCandidates, ColumnRef{
				Table:  tableName,
				Column: column.Name,
			})
		}

		if len(tableCandidates) > 0 {
			candidates[tableName] = tableCandidates
		}
	}
	return candidates
}

func isEliminated(col *ir.Column) (bool, string) {
	if isNullable(col) {
		return true, "column is nullable"
	}
	if isTechnicalColumn(col) {
		return true, "technical metadata column"
	}
	if isLowCardinality(col) {
		return true, "low cardinality column"
	}
	return false, ""
}

func isLowCardinality(col *ir.Column) bool {
	switch strings.ToLower(col.DataType.Name) {
	case "bool", "boolean":
		return true
	}

	name := strings.ToLower(col.Name)
	if strings.HasPrefix(name, "is_") ||
		strings.Contains(name, "flag") ||
		strings.Contains(name, "status") {
		return true
	}

	return false
}

func isNullable(col *ir.Column) bool {
	return col.Nullable
}

func isTechnicalColumn(col *ir.Column) bool {
	switch strings.ToLower(col.Name) {
	case "created_at", "updated_at", "deleted_at", "version":
		return true
	}
	return false
}