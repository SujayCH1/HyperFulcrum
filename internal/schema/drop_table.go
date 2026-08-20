package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyDropTable(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	tables := stmt.Tables
	if len(tables) == 0 {
		tables = []ir.Table{stmt.Table}
	}

	for _, target := range tables {
		tableName := target.Key()

		if _, exists := schema.Tables[tableName]; !exists {
			if stmt.IfExists {
				continue
			}
			return fmt.Errorf(
				"table %q does not exist",
				tableName,
			)
		}
	}

	for _, target := range tables {
		tableName := target.Key()

		if _, exists := schema.Tables[tableName]; !exists {
			continue
		}

		delete(schema.Tables, tableName)

		for _, table := range schema.Tables {
			for name, constraint := range table.Constraints {
				if constraint.Type == ir.ForeignKey &&
					constraint.Reference != nil &&
					constraint.Reference.Table.Key() == tableName {
					delete(table.Constraints, name)
				}
			}
		}
	}

	return nil
}
