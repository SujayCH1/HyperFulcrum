package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func applyAlterTable(
	schema *LogicalSchema,
	stmt *ir.DDLStatement,
) error {

	table, exists := schema.Tables[stmt.Table.Name]
	if !exists {
		return fmt.Errorf("table %q does not exist", stmt.Table.Name)
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

			if err := dropColumn(table, operation.OldName); err != nil {
				return err
			}

		case ir.RenameColumn:

			if err := renameColumn(
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

		default:

			return fmt.Errorf(
				"unsupported alter operation %s",
				operation.Type,
			)
		}
	}

	return nil
}
