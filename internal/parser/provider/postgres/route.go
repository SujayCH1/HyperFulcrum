package postgres

import (
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"

	"hyperfulcrum/internal/parser/ir"
)

type routeBuilder struct {
	statement *ir.RouteStatement
	tables    map[string]struct{}
	ctes      map[string]struct{}
}

func convertRoute(
	node *pg_query.Node,
	command ir.Command,
	metadata ir.Metadata,
) (*ir.RouteStatement, error) {
	builder := &routeBuilder{
		statement: &ir.RouteStatement{
			Metadata:        metadata,
			Cmd:             command,
			RoutingComplete: true,
		},
		tables: make(map[string]struct{}),
		ctes:   make(map[string]struct{}),
	}

	switch command {
	case ir.Select:
		builder.selectStatement(node.GetSelectStmt())
	case ir.Insert:
		builder.insertStatement(node.GetInsertStmt())
	case ir.Update:
		builder.updateStatement(node.GetUpdateStmt())
	case ir.Delete:
		builder.deleteStatement(node.GetDeleteStmt())
	case ir.Merge:
		builder.mergeStatement(node.GetMergeStmt())
	default:
		return nil, ir.ErrUnsupportedSQL
	}

	return builder.statement, nil
}

func (b *routeBuilder) selectStatement(statement *pg_query.SelectStmt) {
	if statement == nil {
		b.statement.RoutingComplete = false
		return
	}

	b.withClause(statement.WithClause)

	if statement.Larg != nil || statement.Rarg != nil {
		b.statement.RoutingComplete = false
		b.selectStatement(statement.Larg)
		b.selectStatement(statement.Rarg)
	}

	for _, node := range statement.FromClause {
		b.fromNode(node)
	}
	b.predicate(statement.WhereClause)
}

func (b *routeBuilder) insertStatement(statement *pg_query.InsertStmt) {
	if statement == nil || statement.Relation == nil {
		b.statement.RoutingComplete = false
		return
	}

	b.withClause(statement.WithClause)
	b.addTable(convertTable(statement.Relation))

	for _, node := range statement.Cols {
		target := node.GetResTarget()
		if target == nil || target.Name == "" {
			b.statement.RoutingComplete = false
			continue
		}
		b.statement.InsertColumns = append(b.statement.InsertColumns, target.Name)
	}

	selectStatement := statement.SelectStmt.GetSelectStmt()
	if selectStatement == nil {
		b.statement.RoutingComplete = false
		return
	}

	if len(selectStatement.ValuesLists) == 0 {
		b.statement.RoutingComplete = false
		b.selectStatement(selectStatement)
		return
	}

	for _, node := range selectStatement.ValuesLists {
		list := node.GetList()
		if list == nil {
			b.statement.RoutingComplete = false
			continue
		}

		row := make([]ir.RouteValue, 0, len(list.Items))
		for _, item := range list.Items {
			value := convertValue(item)
			if value.Kind == ir.UnknownValue {
				b.statement.RoutingComplete = false
			}
			row = append(row, value)
		}
		if len(b.statement.InsertColumns) == 0 || len(row) != len(b.statement.InsertColumns) {
			b.statement.RoutingComplete = false
		}
		b.statement.InsertRows = append(b.statement.InsertRows, row)
	}

	if statement.OnConflictClause != nil {
		b.statement.RoutingComplete = false
		b.predicate(statement.OnConflictClause.WhereClause)
	}
}

func (b *routeBuilder) updateStatement(statement *pg_query.UpdateStmt) {
	if statement == nil || statement.Relation == nil {
		b.statement.RoutingComplete = false
		return
	}

	b.withClause(statement.WithClause)
	b.addTable(convertTable(statement.Relation))
	for _, node := range statement.FromClause {
		b.fromNode(node)
	}
	for _, node := range statement.TargetList {
		target := node.GetResTarget()
		if target == nil || target.Name == "" {
			b.statement.RoutingComplete = false
			continue
		}

		value := convertValue(target.Val)
		if value.Kind == ir.UnknownValue {
			b.statement.RoutingComplete = false
		}
		b.statement.Assignments = append(b.statement.Assignments, ir.RouteAssignment{
			Column: target.Name,
			Value:  value,
		})
	}
	b.predicate(statement.WhereClause)
}

func (b *routeBuilder) deleteStatement(statement *pg_query.DeleteStmt) {
	if statement == nil || statement.Relation == nil {
		b.statement.RoutingComplete = false
		return
	}

	b.withClause(statement.WithClause)
	b.addTable(convertTable(statement.Relation))
	for _, node := range statement.UsingClause {
		b.fromNode(node)
	}
	b.predicate(statement.WhereClause)
}

func (b *routeBuilder) mergeStatement(statement *pg_query.MergeStmt) {
	if statement == nil || statement.Relation == nil {
		b.statement.RoutingComplete = false
		return
	}

	b.withClause(statement.WithClause)
	b.addTable(convertTable(statement.Relation))
	b.fromNode(statement.SourceRelation)
	b.predicate(statement.JoinCondition)
	b.statement.RoutingComplete = false
}

func (b *routeBuilder) withClause(clause *pg_query.WithClause) {
	if clause == nil {
		return
	}

	b.statement.RoutingComplete = false
	for _, node := range clause.Ctes {
		cte := node.GetCommonTableExpr()
		if cte == nil {
			continue
		}
		b.ctes[cte.Ctename] = struct{}{}
		b.statementNode(cte.Ctequery)
	}
}

func (b *routeBuilder) statementNode(node *pg_query.Node) {
	if node == nil {
		return
	}

	switch {
	case node.GetSelectStmt() != nil:
		b.selectStatement(node.GetSelectStmt())
	case node.GetInsertStmt() != nil:
		b.insertStatement(node.GetInsertStmt())
	case node.GetUpdateStmt() != nil:
		b.updateStatement(node.GetUpdateStmt())
	case node.GetDeleteStmt() != nil:
		b.deleteStatement(node.GetDeleteStmt())
	default:
		b.statement.RoutingComplete = false
	}
}

func (b *routeBuilder) fromNode(node *pg_query.Node) {
	if node == nil {
		return
	}

	switch {
	case node.GetRangeVar() != nil:
		table := convertTable(node.GetRangeVar())
		if _, exists := b.ctes[table.Name]; !exists {
			b.addTable(table)
		}
	case node.GetJoinExpr() != nil:
		join := node.GetJoinExpr()
		b.fromNode(join.Larg)
		b.fromNode(join.Rarg)
		b.predicate(join.Quals)
		if join.IsNatural || len(join.UsingClause) > 0 {
			b.statement.RoutingComplete = false
		}
	case node.GetRangeSubselect() != nil:
		b.statement.RoutingComplete = false
		b.statementNode(node.GetRangeSubselect().Subquery)
	default:
		b.statement.RoutingComplete = false
	}
}

func (b *routeBuilder) predicate(node *pg_query.Node) {
	if node == nil {
		return
	}

	if expression := node.GetBoolExpr(); expression != nil {
		if expression.Boolop != pg_query.BoolExprType_AND_EXPR {
			b.statement.RoutingComplete = false
		}
		for _, argument := range expression.Args {
			b.predicate(argument)
		}
		return
	}

	if expression := node.GetAExpr(); expression != nil {
		b.comparison(expression)
		b.subqueries(expression.Lexpr)
		b.subqueries(expression.Rexpr)
		return
	}

	if link := node.GetSubLink(); link != nil {
		b.statement.RoutingComplete = false
		b.statementNode(link.Subselect)
		return
	}

	b.statement.RoutingComplete = false
}

func (b *routeBuilder) subqueries(node *pg_query.Node) {
	if node == nil {
		return
	}

	if link := node.GetSubLink(); link != nil {
		b.statement.RoutingComplete = false
		b.statementNode(link.Subselect)
		return
	}

	if expression := node.GetAExpr(); expression != nil {
		b.subqueries(expression.Lexpr)
		b.subqueries(expression.Rexpr)
		return
	}

	if expression := node.GetBoolExpr(); expression != nil {
		for _, argument := range expression.Args {
			b.subqueries(argument)
		}
		return
	}

	if cast := node.GetTypeCast(); cast != nil {
		b.subqueries(cast.Arg)
	}
}

func (b *routeBuilder) comparison(expression *pg_query.A_Expr) {
	operator := strings.Join(nodeStrings(expression.Name), " ")
	leftTable, leftColumn, leftOK := columnReference(expression.Lexpr)
	rightTable, rightColumn, rightOK := columnReference(expression.Rexpr)

	if leftOK {
		value := convertValue(expression.Rexpr)
		if value.Kind == ir.UnknownValue {
			b.statement.RoutingComplete = false
			return
		}
		b.statement.Predicates = append(b.statement.Predicates, ir.RoutePredicate{
			Table:    leftTable,
			Column:   leftColumn,
			Operator: operator,
			Value:    value,
		})
		return
	}

	if rightOK {
		value := convertValue(expression.Lexpr)
		if value.Kind == ir.UnknownValue {
			b.statement.RoutingComplete = false
			return
		}
		b.statement.Predicates = append(b.statement.Predicates, ir.RoutePredicate{
			Table:    rightTable,
			Column:   rightColumn,
			Operator: reverseOperator(operator),
			Value:    value,
		})
		return
	}

	if leftTable != "" || leftColumn != "" || rightTable != "" || rightColumn != "" {
		return
	}

	b.statement.RoutingComplete = false
}

func reverseOperator(operator string) string {
	switch operator {
	case ">":
		return "<"
	case ">=":
		return "<="
	case "<":
		return ">"
	case "<=":
		return ">="
	default:
		return operator
	}
}

func (b *routeBuilder) addTable(table ir.Table) {
	if table.Name == "" {
		b.statement.RoutingComplete = false
		return
	}

	key := table.Key() + ":" + table.Alias
	if _, exists := b.tables[key]; exists {
		return
	}

	b.tables[key] = struct{}{}
	b.statement.Tables = append(b.statement.Tables, table)
}
