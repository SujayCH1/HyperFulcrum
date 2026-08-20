package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyAlterTable(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	table, exists := schema.Tables[stmt.Table.Key()]
	if !exists {
		if stmt.IfExists {
			return nil
		}
		return fmt.Errorf("table %q does not exist", stmt.Table.Key())
	}

	for _, operation := range stmt.AlterOperations {

		switch operation.Type {

		case ir.AddColumn:

			if operation.Column == nil {
				return fmt.Errorf("missing column for ADD COLUMN")
			}

			if err := addColumn(table, *operation.Column); err != nil {
				return err
			}

		case ir.DropColumn:

			if operation.OldName == "" {
				return fmt.Errorf("missing column name for DROP COLUMN")
			}

			if err := dropColumn(schema, table, operation.OldName); err != nil {
				return err
			}

		case ir.RenameColumn:

			if err := renameColumn(
				schema,
				table,
				operation.OldName,
				operation.NewName,
			); err != nil {
				return err
			}

		case ir.AddConstraint:

			if operation.Constraint == nil {
				return fmt.Errorf("missing constraint for ADD CONSTRAINT")
			}

			if err := addConstraint(table, *operation.Constraint); err != nil {
				return err
			}

		case ir.DropConstraint:

			if operation.OldName == "" {
				return fmt.Errorf("missing constraint name")
			}

			if err := dropConstraintByName(
				table,
				operation.OldName,
			); err != nil {
				return err
			}

		case ir.RenameTable:

			if err := renameTable(
				schema,
				table,
				operation.NewName,
			); err != nil {
				return err
			}

		case ir.AlterColumnType:

			if operation.DataType == nil {
				return fmt.Errorf("missing data type for ALTER COLUMN")
			}

			column, exists := table.Columns[operation.OldName]
			if !exists {
				return fmt.Errorf("column %q does not exist", operation.OldName)
			}

			column.DataType = *operation.DataType

		case ir.SetNotNull:

			column, exists := table.Columns[operation.OldName]
			if !exists {
				return fmt.Errorf("column %q does not exist", operation.OldName)
			}

			column.Nullable = false

		case ir.DropNotNull:

			column, exists := table.Columns[operation.OldName]
			if !exists {
				return fmt.Errorf("column %q does not exist", operation.OldName)
			}

			column.Nullable = true

		case ir.SetDefault:

			column, exists := table.Columns[operation.OldName]
			if !exists {
				return fmt.Errorf("column %q does not exist", operation.OldName)
			}

			column.DefaultValue = operation.Expression

		case ir.DropDefault:

			column, exists := table.Columns[operation.OldName]
			if !exists {
				return fmt.Errorf("column %q does not exist", operation.OldName)
			}

			column.DefaultValue = nil

		default:

			return fmt.Errorf(
				"unsupported alter operation %s",
				operation.Type,
			)
		}
	}

	return nil
}
