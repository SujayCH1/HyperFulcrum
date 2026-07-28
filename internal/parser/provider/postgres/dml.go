package postgres

import (
	"strings"

	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

// convertDML converts a SELECT/INSERT/UPDATE/DELETE ParsedQuery into
// ir.DMLStatement.
//
// ValkDB's SELECT projection list and UPDATE's SET clauses are raw
// expression/text (SelectColumn.Expression, plain "col = expr" strings) —
// not pre-parsed columns — so the mappings below are best-effort text-level
// splits. See parseCondition and parseSetClause.
func convertDML(pq *postgresparser.ParsedQuery) (*ir.DMLStatement, error) {
	stmt := &ir.DMLStatement{
		Cmd:    mapDMLCommand(pq.Command),
		Tables: convertTables(pq.Tables),
		Joins:  convertJoins(pq.Tables),
	}

	switch pq.Command {
	case postgresparser.QueryCommandSelect:
		for _, col := range pq.Columns {
			stmt.Columns = append(stmt.Columns, ir.Column{
				Name:  col.Expression,
				Alias: col.Alias,
			})
		}

	case postgresparser.QueryCommandInsert:
		for _, col := range pq.InsertColumns {
			stmt.Columns = append(stmt.Columns, ir.Column{Name: col})
		}

	case postgresparser.QueryCommandUpdate:
		for _, set := range pq.SetClauses {
			stmt.Columns = append(stmt.Columns, parseSetClause(set))
		}
	}

	for _, w := range pq.Where {
		stmt.Conditions = append(stmt.Conditions, parseCondition(w))
	}

	stmt.GroupBy = toExpressions(pq.GroupBy)

	for _, o := range pq.OrderBy {
		raw := o.Expression
		if o.Direction != "" {
			raw = raw + " " + o.Direction
		}
		stmt.OrderBy = append(stmt.OrderBy, ir.Expression{Raw: raw})
	}

	if pq.Limit != nil {
		if n, ok := parseIntPtr(pq.Limit.Limit); ok {
			stmt.Limit = n
		}
		if n, ok := parseIntPtr(pq.Limit.Offset); ok {
			stmt.Offset = n
		}
	}

	return stmt, nil
}

// convertTables drops nested (CTE/subquery-internal) refs — only top-level
// FROM/JOIN tables are surfaced.
func convertTables(refs []postgresparser.TableRef) []ir.Table {
	var out []ir.Table
	for _, t := range refs {
		if t.Nested {
			continue
		}
		out = append(out, newTable(t.Schema, t.Name, t.Alias))
	}
	return out
}

// convertJoins pairs each joined TableRef with the table immediately before
// it in source order. This matches simple chained joins; genuinely branching
// join trees (e.g. joining back to an earlier alias) won't reconstruct
// correctly from this flat list alone — flag it if you hit that case.
func convertJoins(refs []postgresparser.TableRef) []ir.Join {
	var out []ir.Join
	var prev *postgresparser.TableRef

	for i := range refs {
		t := &refs[i]
		if t.Nested {
			continue
		}
		if t.JoinType != "" && prev != nil {
			out = append(out, ir.Join{
				Type:      mapJoinType(t.JoinType),
				Left:      newTable(prev.Schema, prev.Name, prev.Alias),
				Right:     newTable(t.Schema, t.Name, t.Alias),
				Condition: ir.Expression{Raw: t.JoinCondition},
			})
		}
		prev = t
	}

	return out
}

// parseSetClause splits a raw "col = expr" UPDATE SET fragment. The assigned
// value is stashed in DefaultValue since ir.Column has no dedicated
// "new value" field.
func parseSetClause(raw string) ir.Column {
	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return ir.Column{Name: strings.TrimSpace(raw)}
	}
	return ir.Column{
		Name:         strings.TrimSpace(parts[0]),
		DefaultValue: nonEmptyExpr(strings.TrimSpace(parts[1])),
	}
}
