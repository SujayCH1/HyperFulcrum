package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func ApplyBatch(
	schema *LogicalSchema,
	batch *ir.Batch,
) error {

	if batch == nil {
		return ir.ErrNilStatement
	}

	for _, stmt := range batch.Statements {
		if err := ApplyStatement(schema, stmt); err != nil {
			return err
		}
	}

	return nil
}

func ApplyStatement(
	schema *LogicalSchema,
	statement ir.Statement,
) error {

	if statement == nil {
		return ir.ErrNilStatement
	}

	switch stmt := statement.(type) {

	case *ir.DDLStatement:
		return applyDDLStatement(schema, stmt)

	default:
		return fmt.Errorf(
			"schema: unsupported statement kind %s",
			statement.Kind(),
		)
	}
}

func applyDDLStatement(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	switch stmt.Cmd {

	case ir.CreateTable:
		return applyCreateTable(schema, stmt)

	case ir.AlterTable:
		return applyAlterTable(schema, stmt)

	case ir.DropTable:
		return applyDropTable(schema, stmt)

	default:
		return fmt.Errorf(
			"schema: unsupported DDL command %s",
			stmt.Cmd,
		)
	}
}
