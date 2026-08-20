package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyCreateIndex(
	schema *LogicalSchema,
	statement *ir.DDLStatement,
) error {
	table, exists := schema.Tables[statement.Table.Key()]
	if !exists {
		return fmt.Errorf("table %q does not exist", statement.Table.Key())
	}

	for _, index := range statement.Indexes {
		if _, exists := table.Indexes[index.Name]; exists && statement.IfNotExists {
			continue
		}
		if err := addIndex(table, index); err != nil {
			return err
		}
	}

	return nil
}

func applyDropIndex(
	schema *LogicalSchema,
	statement *ir.DDLStatement,
) error {
	for _, target := range statement.Tables {
		found := false
		for _, table := range schema.Tables {
			if _, exists := table.Indexes[target.Name]; exists {
				dropIndex(table, target.Name)
				found = true
				break
			}
		}

		if !found {
			if statement.IfExists {
				continue
			}
			return fmt.Errorf("index %q does not exist", target.Name)
		}
	}

	return nil
}
