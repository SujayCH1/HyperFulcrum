package postgres

import (
	"strconv"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"

	"hyperfulcrum/internal/parser/ir"
)

func convertTable(table *pg_query.RangeVar) ir.Table {
	if table == nil {
		return ir.Table{}
	}

	alias := ""
	if table.Alias != nil {
		alias = table.Alias.Aliasname
	}

	return ir.Table{
		Schema: table.Schemaname,
		Name:   table.Relname,
		Alias:  alias,
	}
}

func nodeString(node *pg_query.Node) string {
	if node == nil || node.GetString_() == nil {
		return ""
	}

	return node.GetString_().Sval
}

func nodeStrings(nodes []*pg_query.Node) []string {
	values := make([]string, 0, len(nodes))
	for _, node := range nodes {
		if value := nodeString(node); value != "" {
			values = append(values, value)
		}
	}
	return values
}

func objectTable(node *pg_query.Node) ir.Table {
	if node == nil || node.GetList() == nil {
		return ir.Table{}
	}

	values := nodeStrings(node.GetList().Items)
	if len(values) == 1 {
		return ir.Table{Name: values[0]}
	}
	if len(values) >= 2 {
		return ir.Table{Schema: values[len(values)-2], Name: values[len(values)-1]}
	}

	return ir.Table{}
}

func convertDataType(dataType *pg_query.TypeName) (ir.DataType, error) {
	if dataType == nil {
		return ir.DataType{}, ir.ErrInvalidAST
	}

	names := nodeStrings(dataType.Names)
	if len(names) > 1 && names[0] == "pg_catalog" {
		names = names[1:]
	}

	result := ir.DataType{
		Name:            normalizeDataType(strings.Join(names, ".")),
		ArrayDimensions: len(dataType.ArrayBounds),
	}

	for _, modifier := range dataType.Typmods {
		value, err := expressionSQL(modifier)
		if err != nil {
			return ir.DataType{}, err
		}
		result.Modifiers = append(result.Modifiers, value)
	}

	return result, nil
}

func normalizeDataType(value string) string {
	switch value {
	case "int2":
		return "smallint"
	case "int4":
		return "integer"
	case "int8":
		return "bigint"
	case "float4":
		return "real"
	case "float8":
		return "double precision"
	case "bool":
		return "boolean"
	case "bpchar":
		return "character"
	case "varchar":
		return "character varying"
	default:
		return value
	}
}

func expressionSQL(node *pg_query.Node) (string, error) {
	if node == nil {
		return "", ir.ErrInvalidAST
	}

	tree := &pg_query.ParseResult{
		Stmts: []*pg_query.RawStmt{{
			Stmt: &pg_query.Node{
				Node: &pg_query.Node_SelectStmt{
					SelectStmt: &pg_query.SelectStmt{
						TargetList: []*pg_query.Node{{
							Node: &pg_query.Node_ResTarget{
								ResTarget: &pg_query.ResTarget{Val: node},
							},
						}},
					},
				},
			},
		}},
	}

	value, err := pg_query.Deparse(tree)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.TrimPrefix(value, "SELECT")), nil
}

func convertValue(node *pg_query.Node) ir.RouteValue {
	if node == nil {
		return ir.RouteValue{Kind: ir.UnknownValue}
	}

	if cast := node.GetTypeCast(); cast != nil {
		return convertValue(cast.Arg)
	}

	if expression := node.GetAExpr(); expression != nil &&
		expression.Lexpr == nil &&
		len(expression.Name) == 1 {
		operator := nodeString(expression.Name[0])
		value := convertValue(expression.Rexpr)
		if value.Kind == ir.LiteralValue && (operator == "-" || operator == "+") {
			value.Value = operator + value.Value
			return value
		}
	}

	if parameter := node.GetParamRef(); parameter != nil {
		return ir.RouteValue{
			Kind:      ir.ParameterValue,
			Parameter: int(parameter.Number),
		}
	}

	if constant := node.GetAConst(); constant != nil {
		value := ir.RouteValue{Kind: ir.LiteralValue}

		switch {
		case constant.Isnull:
			value.Value = "NULL"
		case constant.GetIval() != nil:
			value.Value = strconv.FormatInt(int64(constant.GetIval().Ival), 10)
		case constant.GetFval() != nil:
			value.Value = constant.GetFval().Fval
		case constant.GetBoolval() != nil:
			value.Value = strconv.FormatBool(constant.GetBoolval().Boolval)
		case constant.GetSval() != nil:
			value.Value = constant.GetSval().Sval
		default:
			value.Kind = ir.UnknownValue
		}

		return value
	}

	if table, column, ok := columnReference(node); ok {
		return ir.RouteValue{
			Kind:   ir.ColumnValue,
			Table:  table,
			Column: column,
		}
	}

	return ir.RouteValue{Kind: ir.UnknownValue}
}

func columnReference(node *pg_query.Node) (string, string, bool) {
	if node == nil || node.GetColumnRef() == nil {
		return "", "", false
	}

	fields := nodeStrings(node.GetColumnRef().Fields)
	if len(fields) == 0 {
		return "", "", false
	}

	if len(fields) == 1 {
		return "", fields[0], true
	}

	return fields[len(fields)-2], fields[len(fields)-1], true
}
