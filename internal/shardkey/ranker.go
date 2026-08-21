package shardkey

import (
	"fmt"
	"sort"
	"strings"

	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/schema"
)

func RankTableCandidates(
	tableName string,
	local []ColumnRef,
	fanout map[ColumnRef]FanoutStats,
	s *schema.LogicalSchema,
) []RankedCandidate {
	var ranked []RankedCandidate

	table := s.Tables[tableName]

	for _, col := range local {
		stats, ok := fanout[col]
		if !ok {
			stats = FanoutStats{}
		}

		column := table.Columns[col.Column]

		score, reasons := scoreColumn(col, column, stats, table, fanout)

		ranked = append(ranked, RankedCandidate{
			Column:  col,
			Score:   score,
			Reasons: reasons,
		})
	}

	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score != ranked[j].Score {
			return ranked[i].Score > ranked[j].Score
		}
		return tieBreak(ranked[i].Column, ranked[j].Column)
	})

	return ranked
}

func scoreColumn(
	col ColumnRef,
	column *ir.Column,
	stats FanoutStats,
	table *schema.LogicalTable,
	fanout map[ColumnRef]FanoutStats,
) (int, []string) {
	score := 0
	var reasons []string

	if stats.IncomingFKs > 0 {
		value := stats.IncomingFKs * 10
		score += value
		reasons = append(reasons,
			fmt.Sprintf("referenced by %d foreign keys", stats.IncomingFKs),
		)
	}

	if stats.ReferencingTables > 0 {
		value := stats.ReferencingTables * 5
		score += value
		reasons = append(reasons,
			fmt.Sprintf("shared across %d tables", stats.ReferencingTables),
		)
	}

	if isForeignKey(col.Column, table) {
		score += 20
		reasons = append(reasons, "foreign key (ownership column)")

		if bonus, reason := rootAffinityBonus(col.Column, table, fanout); bonus > 0 {
			score += bonus
			reasons = append(reasons, reason)
		}
	}

	if isColumnPrimaryKey(col.Column, table) {
		score += 10
		reasons = append(reasons, "primary key (identity column)")
	}

	if isTextualColumn(column) {
		score -= 15
		reasons = append(reasons, "textual/content column")
	}

	score += 1
	reasons = append(reasons, "local column")

	return score, reasons
}

func rootAffinityBonus(
	columnName string,
	table *schema.LogicalTable,
	fanout map[ColumnRef]FanoutStats,
) (int, string) {
	for _, constraint := range table.Constraints {
		if constraint.Type != ir.ForeignKey || constraint.Reference == nil {
			continue
		}
		if len(constraint.Columns) == 0 || constraint.Columns[0] != columnName {
			continue
		}

		parent := ColumnRef{
			Table:  constraint.Reference.Table,
			Column: constraint.Reference.Columns[0],
		}

		stats, ok := fanout[parent]
		if !ok {
			return 0, ""
		}

		if stats.IncomingFKs > 0 {
			bonus := stats.IncomingFKs * 5
			return bonus, fmt.Sprintf(
				"points to root table (%d incoming references)",
				stats.IncomingFKs,
			)
		}
	}

	return 0, ""
}

func isForeignKey(columnName string, table *schema.LogicalTable) bool {
	for _, constraint := range table.Constraints {
		if constraint.Type != ir.ForeignKey {
			continue
		}
		for _, col := range constraint.Columns {
			if col == columnName {
				return true
			}
		}
	}
	return false
}

func isColumnPrimaryKey(columnName string, table *schema.LogicalTable) bool {
	for _, constraint := range table.Constraints {
		if constraint.Type != ir.PrimaryKey {
			continue
		}
		for _, col := range constraint.Columns {
			if col == columnName {
				return true
			}
		}
	}
	return false
}

func isTextualColumn(col *ir.Column) bool {
	switch strings.ToLower(col.DataType.Name) {
	case "text", "varchar", "char", "character varying":
		return true
	}
	return false
}

func tieBreak(a, b ColumnRef) bool {
	if a.Table != b.Table {
		return a.Table < b.Table
	}
	return a.Column < b.Column
}