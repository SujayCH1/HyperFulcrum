package postgres

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"

	"hyperfulcrum/internal/parser/ir"
)

func convertCreateTable(
	statement *pg_query.CreateStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	if statement.Relation == nil {
		return nil, ir.ErrInvalidAST
	}

	result := &ir.DDLStatement{
		Metadata:    metadata,
		Cmd:         ir.CreateTable,
		Table:       convertTable(statement.Relation),
		IfNotExists: statement.IfNotExists,
	}

	for _, element := range statement.TableElts {
		switch {
		case element.GetColumnDef() != nil:
			column, constraints, err := convertColumn(element.GetColumnDef())
			if err != nil {
				return nil, err
			}
			result.Columns = append(result.Columns, column)
			result.Constraints = append(result.Constraints, constraints...)
		case element.GetConstraint() != nil:
			constraint, err := convertConstraint(element.GetConstraint(), "")
			if err != nil {
				return nil, err
			}
			if constraint != nil {
				result.Constraints = append(result.Constraints, *constraint)
			}
		default:
			return nil, ir.ErrUnsupportedSQL
		}
	}

	for _, node := range statement.Constraints {
		constraint, err := convertConstraint(node.GetConstraint(), "")
		if err != nil {
			return nil, err
		}
		if constraint != nil {
			result.Constraints = append(result.Constraints, *constraint)
		}
	}

	return result, nil
}

func convertColumn(definition *pg_query.ColumnDef) (ir.Column, []ir.Constraint, error) {
	dataType, err := convertDataType(definition.TypeName)
	if err != nil {
		return ir.Column{}, nil, err
	}

	column := ir.Column{
		Name:     definition.Colname,
		DataType: dataType,
		Nullable: !definition.IsNotNull,
	}

	if definition.RawDefault != nil {
		value, err := expressionSQL(definition.RawDefault)
		if err != nil {
			return ir.Column{}, nil, err
		}
		column.DefaultValue = &ir.Expression{Raw: value}
	}

	constraints := make([]ir.Constraint, 0, len(definition.Constraints))
	for _, node := range definition.Constraints {
		constraint := node.GetConstraint()
		if constraint == nil {
			return ir.Column{}, nil, ir.ErrInvalidAST
		}

		switch constraint.Contype {
		case pg_query.ConstrType_CONSTR_NOTNULL:
			column.Nullable = false
		case pg_query.ConstrType_CONSTR_NULL:
			column.Nullable = true
		case pg_query.ConstrType_CONSTR_DEFAULT:
			value, err := expressionSQL(constraint.RawExpr)
			if err != nil {
				return ir.Column{}, nil, err
			}
			column.DefaultValue = &ir.Expression{Raw: value}
		case pg_query.ConstrType_CONSTR_IDENTITY:
			column.Identity = identityType(constraint.GeneratedWhen)
		case pg_query.ConstrType_CONSTR_GENERATED:
			value, err := expressionSQL(constraint.RawExpr)
			if err != nil {
				return ir.Column{}, nil, err
			}
			column.Generated = &ir.Expression{Raw: value}
		default:
			converted, err := convertConstraint(constraint, definition.Colname)
			if err != nil {
				return ir.Column{}, nil, err
			}
			if converted != nil {
				constraints = append(constraints, *converted)
			}
		}
	}

	return column, constraints, nil
}

func convertConstraint(
	constraint *pg_query.Constraint,
	column string,
) (*ir.Constraint, error) {
	if constraint == nil {
		return nil, ir.ErrInvalidAST
	}

	result := &ir.Constraint{
		Name:              constraint.Conname,
		Columns:           nodeStrings(constraint.Keys),
		Deferrable:        constraint.Deferrable,
		InitiallyDeferred: constraint.Initdeferred,
		NullsNotDistinct:  constraint.NullsNotDistinct,
	}

	if len(result.Columns) == 0 && column != "" {
		result.Columns = []string{column}
	}

	switch constraint.Contype {
	case pg_query.ConstrType_CONSTR_PRIMARY:
		result.Type = ir.PrimaryKey
	case pg_query.ConstrType_CONSTR_UNIQUE:
		result.Type = ir.Unique
	case pg_query.ConstrType_CONSTR_CHECK:
		value, err := expressionSQL(constraint.RawExpr)
		if err != nil {
			return nil, err
		}
		result.Type = ir.Check
		result.Expression = &ir.Expression{Raw: value}
	case pg_query.ConstrType_CONSTR_FOREIGN:
		result.Type = ir.ForeignKey
		result.Columns = nodeStrings(constraint.FkAttrs)
		if len(result.Columns) == 0 && column != "" {
			result.Columns = []string{column}
		}
		result.Reference = &ir.Reference{
			Table:    convertTable(constraint.Pktable),
			Columns:  nodeStrings(constraint.PkAttrs),
			OnUpdate: foreignKeyAction(constraint.FkUpdAction),
			OnDelete: foreignKeyAction(constraint.FkDelAction),
		}
	case pg_query.ConstrType_CONSTR_NULL,
		pg_query.ConstrType_CONSTR_NOTNULL,
		pg_query.ConstrType_CONSTR_DEFAULT:
		return nil, nil
	default:
		return nil, unsupported(constraint.Contype.String())
	}

	return result, nil
}

func identityType(value string) string {
	if value == "a" {
		return "ALWAYS"
	}
	if value == "d" {
		return "BY DEFAULT"
	}
	return ""
}

func foreignKeyAction(action string) string {
	switch action {
	case "a":
		return "NO ACTION"
	case "r":
		return "RESTRICT"
	case "c":
		return "CASCADE"
	case "n":
		return "SET NULL"
	case "d":
		return "SET DEFAULT"
	default:
		return ""
	}
}

func convertAlterTable(
	statement *pg_query.AlterTableStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	if statement.Relation == nil {
		return nil, ir.ErrInvalidAST
	}

	result := &ir.DDLStatement{
		Metadata: metadata,
		Cmd:      ir.AlterTable,
		Table:    convertTable(statement.Relation),
		IfExists: statement.MissingOk,
	}

	for _, node := range statement.Cmds {
		command := node.GetAlterTableCmd()
		if command == nil {
			return nil, ir.ErrInvalidAST
		}

		operations, err := convertAlterOperations(command)
		if err != nil {
			return nil, err
		}
		result.AlterOperations = append(result.AlterOperations, operations...)
	}

	return result, nil
}

func convertAlterOperations(command *pg_query.AlterTableCmd) ([]ir.AlterOperation, error) {
	switch command.Subtype {
	case pg_query.AlterTableType_AT_AddColumn:
		column, constraints, err := convertColumn(command.Def.GetColumnDef())
		if err != nil {
			return nil, err
		}
		operations := []ir.AlterOperation{{Type: ir.AddColumn, Column: &column}}
		for i := range constraints {
			constraint := constraints[i]
			operations = append(operations, ir.AlterOperation{
				Type:       ir.AddConstraint,
				Constraint: &constraint,
			})
		}
		return operations, nil
	case pg_query.AlterTableType_AT_DropColumn:
		return []ir.AlterOperation{{Type: ir.DropColumn, OldName: command.Name}}, nil
	case pg_query.AlterTableType_AT_AddConstraint:
		constraint, err := convertConstraint(command.Def.GetConstraint(), "")
		if err != nil {
			return nil, err
		}
		if constraint == nil {
			return nil, ir.ErrInvalidAST
		}
		return []ir.AlterOperation{{Type: ir.AddConstraint, Constraint: constraint}}, nil
	case pg_query.AlterTableType_AT_DropConstraint:
		return []ir.AlterOperation{{Type: ir.DropConstraint, OldName: command.Name}}, nil
	case pg_query.AlterTableType_AT_AlterColumnType:
		definition := command.Def.GetColumnDef()
		if definition == nil {
			return nil, ir.ErrInvalidAST
		}
		dataType, err := convertDataType(definition.TypeName)
		if err != nil {
			return nil, err
		}
		return []ir.AlterOperation{{
			Type:     ir.AlterColumnType,
			OldName:  command.Name,
			DataType: &dataType,
		}}, nil
	case pg_query.AlterTableType_AT_SetNotNull:
		return []ir.AlterOperation{{Type: ir.SetNotNull, OldName: command.Name}}, nil
	case pg_query.AlterTableType_AT_DropNotNull:
		return []ir.AlterOperation{{Type: ir.DropNotNull, OldName: command.Name}}, nil
	case pg_query.AlterTableType_AT_ColumnDefault:
		if command.Def == nil {
			return []ir.AlterOperation{{Type: ir.DropDefault, OldName: command.Name}}, nil
		}
		value, err := expressionSQL(command.Def)
		if err != nil {
			return nil, err
		}
		return []ir.AlterOperation{{
			Type:       ir.SetDefault,
			OldName:    command.Name,
			Expression: &ir.Expression{Raw: value},
		}}, nil
	default:
		return nil, unsupported(command.Subtype.String())
	}
}

func convertRename(
	statement *pg_query.RenameStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	result := &ir.DDLStatement{
		Metadata: metadata,
		Cmd:      ir.AlterTable,
		Table:    convertTable(statement.Relation),
		IfExists: statement.MissingOk,
	}

	switch statement.RenameType {
	case pg_query.ObjectType_OBJECT_COLUMN:
		result.AlterOperations = []ir.AlterOperation{{
			Type:    ir.RenameColumn,
			OldName: statement.Subname,
			NewName: statement.Newname,
		}}
	case pg_query.ObjectType_OBJECT_TABLE:
		result.AlterOperations = []ir.AlterOperation{{
			Type:    ir.RenameTable,
			NewName: statement.Newname,
		}}
	default:
		return nil, unsupported(statement.RenameType.String())
	}

	return result, nil
}

func convertDrop(
	statement *pg_query.DropStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	result := &ir.DDLStatement{
		Metadata: metadata,
		IfExists: statement.MissingOk,
		Cascade:  statement.Behavior == pg_query.DropBehavior_DROP_CASCADE,
	}

	switch statement.RemoveType {
	case pg_query.ObjectType_OBJECT_TABLE:
		result.Cmd = ir.DropTable
	case pg_query.ObjectType_OBJECT_INDEX:
		result.Cmd = ir.DropIndex
	default:
		return nil, unsupported(statement.RemoveType.String())
	}

	for _, object := range statement.Objects {
		result.Tables = append(result.Tables, objectTable(object))
	}
	if len(result.Tables) == 0 {
		return nil, ir.ErrInvalidAST
	}
	result.Table = result.Tables[0]

	return result, nil
}

func convertCreateIndex(
	statement *pg_query.IndexStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	table := convertTable(statement.Relation)
	index := ir.Index{
		Name:   statement.Idxname,
		Table:  table,
		Unique: statement.Unique,
	}

	for _, node := range statement.IndexParams {
		element := node.GetIndexElem()
		if element == nil {
			return nil, ir.ErrInvalidAST
		}
		if element.Name != "" {
			index.Columns = append(index.Columns, element.Name)
			continue
		}
		value, err := expressionSQL(element.Expr)
		if err != nil {
			return nil, err
		}
		index.Expressions = append(index.Expressions, ir.Expression{Raw: value})
	}

	for _, node := range statement.IndexIncludingParams {
		element := node.GetIndexElem()
		if element == nil || element.Name == "" {
			return nil, ir.ErrInvalidAST
		}
		index.Include = append(index.Include, element.Name)
	}

	if statement.WhereClause != nil {
		value, err := expressionSQL(statement.WhereClause)
		if err != nil {
			return nil, err
		}
		index.Predicate = &ir.Expression{Raw: value}
	}

	return &ir.DDLStatement{
		Metadata:    metadata,
		Cmd:         ir.CreateIndex,
		Table:       table,
		Indexes:     []ir.Index{index},
		IfNotExists: statement.IfNotExists,
	}, nil
}

func convertTruncate(
	statement *pg_query.TruncateStmt,
	metadata ir.Metadata,
) (*ir.DDLStatement, error) {
	result := &ir.DDLStatement{
		Metadata: metadata,
		Cmd:      ir.Truncate,
		Cascade:  statement.Behavior == pg_query.DropBehavior_DROP_CASCADE,
	}

	for _, node := range statement.Relations {
		if node.GetRangeVar() == nil {
			return nil, ir.ErrInvalidAST
		}
		result.Tables = append(result.Tables, convertTable(node.GetRangeVar()))
	}
	if len(result.Tables) == 0 {
		return nil, ir.ErrInvalidAST
	}
	result.Table = result.Tables[0]

	return result, nil
}
