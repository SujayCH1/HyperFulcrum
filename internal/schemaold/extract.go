package schema

import pg_query "github.com/pganalyze/pg_query_go/v5"

func extractSchemaFromAST(ast *pg_query.ParseResult, logicalSchema *LogicalSchema) error {

	for _, stmt := range ast.Stmts {
		node := stmt.Stmt.Node
		switch node := node.(type) {
		case *pg_query.Node_CreateStmt:
			if err := extractCreateTable(node.CreateStmt, logicalSchema); err != nil {
				return err
			}
		case *pg_query.Node_AlterTableStmt:
			if err := extractAlterTable(node.AlterTableStmt, logicalSchema); err != nil {
				return err
			}
		}
	}
	return nil
}

func extractCreateTable(stmt *pg_query.CreateStmt, logicalSchema *LogicalSchema) error {
	tableName := stmt.Relation.Relname
	if _, ok := logicalSchema.Tables[tableName]; !ok {
		logicalSchema.Tables[tableName] = &Table{
			Columns: make(map[string]*Column),
			Fks:     make(map[Fkkey]*Fk),
		}
	}
	for _, elt := range stmt.TableElts {
		switch e := elt.Node.(type) {

		case *pg_query.Node_ColumnDef:
			extractColumnDef(tableName, e.ColumnDef, logicalSchema)

		case *pg_query.Node_Constraint:
			conType := e.Constraint.Contype
			if conType == pg_query.ConstrType_CONSTR_FOREIGN {
				if err := extractFKConstraint(tableName, e.Constraint, logicalSchema); err != nil {
					return err
				}
			}

			if conType == pg_query.ConstrType_CONSTR_PRIMARY {
				for _, key := range e.Constraint.Keys {
					colName := getStringFromNode(key)
					col, ok := logicalSchema.Tables[tableName].Columns[colName]
					if ok {
						col.IsPk = true
					}
				}
			}
		}
	}
	return nil
}

func extractAlterTable(stmt *pg_query.AlterTableStmt, schema *LogicalSchema) error {

	tableName := stmt.Relation.Relname

	if _, ok := schema.Tables[tableName]; !ok {
		schema.Tables[tableName] = &Table{
			Columns: make(map[string]*Column),
			Fks:     make(map[Fkkey]*Fk),
		}
	}

	for _, cmd := range stmt.Cmds {
		c := cmd.Node.(*pg_query.Node_AlterTableCmd).AlterTableCmd

		switch c.Subtype {

		case pg_query.AlterTableType_AT_AddColumn:
			colDef := c.Def.Node.(*pg_query.Node_ColumnDef).ColumnDef
			extractColumnDef(tableName, colDef, schema)

		case pg_query.AlterTableType_AT_AddConstraint:
			con := c.Def.Node.(*pg_query.Node_Constraint).Constraint
			conType := con.Contype
			if conType == pg_query.ConstrType_CONSTR_FOREIGN {
				if err := extractFKConstraint(tableName, con, schema); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func extractColumnDef(tableName string, colDef *pg_query.ColumnDef, schema *LogicalSchema) {
	colName := colDef.Colname
	nullable := true
	isPK := false
	for _, c := range colDef.Constraints {
		con := c.Node.(*pg_query.Node_Constraint).Constraint
		conType := con.Contype
		if conType == pg_query.ConstrType_CONSTR_NOTNULL {
			nullable = false
		}
		if conType == pg_query.ConstrType_CONSTR_PRIMARY {
			isPK = true
		}
	}
	dataType := extractTypeName(colDef.TypeName)
	schema.Tables[tableName].Columns[colName] = &Column{
		Name:     colName,
		DataType: dataType,
		Nullable: nullable,
		IsPk:     isPK,
	}
}

func extractFKConstraint(tableName string, constraint *pg_query.Constraint, schema *LogicalSchema) error {
	parentTable := constraint.Pktable.Relname
	for i, key := range constraint.FkAttrs {
		childCol := getStringFromNode(key)
		parentCol := getStringFromNode(constraint.PkAttrs[i])
		fkKey := Fkkey{
			ChildColumn:  childCol,
			ParentTable:  parentTable,
			ParentColumn: parentCol,
		}
		schema.Tables[tableName].Fks[fkKey] = &Fk{
			ParentTable:  parentTable,
			ParentColumn: parentCol,
			ChildColumn:  childCol,
		}
	}
	return nil
}

func extractTypeName(typeName *pg_query.TypeName) string {
	if typeName == nil || len(typeName.Names) == 0 {
		return ""
	}
	lastNode := typeName.Names[len(typeName.Names)-1]
	return getStringFromNode(lastNode)
}

func getStringFromNode(n *pg_query.Node) string {
	if n == nil {
		return ""
	}
	switch v := n.Node.(type) {
	case *pg_query.Node_String_:
		return v.String_.Sval
	default:
		return ""
	}
}
