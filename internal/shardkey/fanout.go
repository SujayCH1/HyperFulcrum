package shardkey

import (
	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/schema"
)

func ComputeFanout(s *schema.LogicalSchema, candidates CandidateSet) map[ColumnRef]FanoutStats {
	candidateIndex := indexCandidates(candidates)
	fanout := make(map[ColumnRef]FanoutStats)
	seenTables := make(map[ColumnRef]map[string]struct{})

	for _, table := range s.Tables {
		for _, constraint := range table.Constraints {
			if constraint.Type != ir.ForeignKey || constraint.Reference == nil {
				continue
			}

			if len(constraint.Columns) == 0 || len(constraint.Reference.Columns) == 0 {
				continue
			}

			parentCol := ColumnRef{
				Table:  constraint.Reference.Table,
				Column: constraint.Reference.Columns[0],
			}

			if _, ok := candidateIndex[parentCol]; !ok {
				continue
			}

			stats := fanout[parentCol]
			stats.IncomingFKs++

			if _, ok := seenTables[parentCol]; !ok {
				seenTables[parentCol] = make(map[string]struct{})
			}

			childTable := table.Name
			if _, seen := seenTables[parentCol][childTable]; !seen {
				seenTables[parentCol][childTable] = struct{}{}
				stats.ReferencingTables++
			}

			fanout[parentCol] = stats
		}
	}

	return fanout
}

func indexCandidates(candidates CandidateSet) map[ColumnRef]struct{} {
	index := make(map[ColumnRef]struct{})

	for tableName, cols := range candidates {
		for _, col := range cols {
			index[ColumnRef{
				Table:  tableName,
				Column: col.Column,
			}] = struct{}{}
		}
	}

	return index
}