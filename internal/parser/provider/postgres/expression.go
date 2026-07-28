package postgres

import (
	"strings"

	"hyperfulcrum/internal/parser/ir"
)

// conditionOperators is checked in order, longest-first, so "<=" isn't
// mis-split as "<" + "=".
var conditionOperators = []string{
	"!=", "<>", ">=", "<=", "=", ">", "<",
	" NOT LIKE ", " LIKE ", " NOT IN ", " IN ", " IS NOT ", " IS ",
}

// parseCondition does a best-effort split of a raw predicate string (as
// returned by ValkDB in ParsedQuery.Where / JoinConditions — plain text, not
// a structured expression tree) into an ir.Condition. If no known operator
// is found, the whole fragment becomes the left-hand expression rather than
// being dropped.
func parseCondition(raw string) ir.Condition {
	raw = strings.TrimSpace(raw)

	for _, op := range conditionOperators {
		if idx := indexOperator(raw, op); idx != -1 {
			left := strings.TrimSpace(raw[:idx])
			right := strings.TrimSpace(raw[idx+len(op):])
			return ir.Condition{
				Left:     ir.Expression{Raw: left},
				Operator: strings.TrimSpace(op),
				Right:    ir.Expression{Raw: right},
			}
		}
	}

	return ir.Condition{Left: ir.Expression{Raw: raw}}
}

func indexOperator(s, op string) int {
	return strings.Index(strings.ToUpper(s), strings.ToUpper(op))
}

func toExpressions(raw []string) []ir.Expression {
	exprs := make([]ir.Expression, 0, len(raw))
	for _, r := range raw {
		exprs = append(exprs, ir.Expression{Raw: strings.TrimSpace(r)})
	}
	return exprs
}
