package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyCreateTable(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	tableName := stmt.Table.Key()

	if _, exists := schema.Tables[tableName]; exists {
		if stmt.IfNotExists {
			return nil
		}
		return fmt.Errorf("table %q already exists", tableName)
	}

	table := NewLogicalTable(stmt.Table)

	for _, column := range stmt.Columns {
		if err := addColumn(table, column); err != nil {
			return err
		}
	}

	for _, constraint := range stmt.Constraints {
		if err := addConstraint(table, constraint); err != nil {
			return err
		}
	}

	for _, index := range stmt.Indexes {
		if err := addIndex(table, index); err != nil {
			return err
		}
	}

	schema.Tables[tableName] = table

	return nil
}
