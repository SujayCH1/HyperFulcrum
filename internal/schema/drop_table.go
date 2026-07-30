package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyDropTable(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	tableName := stmt.Table.Name

	if _, exists := schema.Tables[tableName]; !exists {
		return fmt.Errorf(
			"table %q does not exist",
			tableName,
		)
	}

	delete(schema.Tables, tableName)

	return nil
}
