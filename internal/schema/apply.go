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
	if schema == nil {
		return fmt.Errorf("schema is nil")
	}

	working := schema.Clone()

	for _, stmt := range batch.Statements {
		if err := ApplyStatement(working, stmt); err != nil {
			return err
		}
	}

	*schema = *working

	return nil
}

func ApplyStatement(
	schema *LogicalSchema,
	statement ir.Statement,
) error {

	if statement == nil {
		return ir.ErrNilStatement
	}
	if schema == nil {
		return fmt.Errorf("schema is nil")
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

	case ir.CreateIndex:
		return applyCreateIndex(schema, stmt)

	case ir.DropIndex:
		return applyDropIndex(schema, stmt)

	default:
		return fmt.Errorf(
			"schema: unsupported DDL command %s",
			stmt.Cmd,
		)
	}
}
