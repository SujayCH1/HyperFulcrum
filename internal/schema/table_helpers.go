package schema

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
)

func addColumn(
	table *LogicalTable,
	column ir.Column,
) error {

	if _, exists := table.Columns[column.Name]; exists {
		return fmt.Errorf(
			"column %q already exists",
			column.Name,
		)
	}

	column.Table = table.Name

	c := column

	table.Columns[c.Name] = &c

	return nil
}

func dropColumn(
	table *LogicalTable,
	columnName string,
) error {

	if _, exists := table.Columns[columnName]; !exists {
		return fmt.Errorf("column %q does not exist", columnName)
	}

	delete(table.Columns, columnName)

	return nil
}

func addConstraint(
	table *LogicalTable,
	constraint ir.Constraint,
) error {

	ensureConstraintName(table.Name, &constraint)

	if _, exists := table.Constraints[constraint.Name]; exists {
		return fmt.Errorf(
			"constraint %q already exists",
			constraint.Name,
		)
	}

	c := constraint

	table.Constraints[c.Name] = &c

	return nil
}

func dropConstraint(
	table *LogicalTable,
	key string,
) {

	delete(table.Constraints, key)
}

func addIndex(
	table *LogicalTable,
	index ir.Index,
) error {

	if _, exists := table.Indexes[index.Name]; exists {
		return fmt.Errorf("index %q already exists", index.Name)
	}

	idx := index

	table.Indexes[index.Name] = &idx

	return nil
}

func dropIndex(
	table *LogicalTable,
	name string,
) {

	delete(table.Indexes, name)
}

func renameColumn(
	table *LogicalTable,
	oldName string,
	newName string,
) error {

	column, exists := table.Columns[oldName]
	if !exists {
		return fmt.Errorf("column %q does not exist", oldName)
	}

	if _, exists := table.Columns[newName]; exists {
		return fmt.Errorf("column %q already exists", newName)
	}

	delete(table.Columns, oldName)

	column.Name = newName

	table.Columns[newName] = column

	// Update every constraint that references this column.
	for _, constraint := range table.Constraints {

		for i, columnName := range constraint.Columns {
			if columnName == oldName {
				constraint.Columns[i] = newName
			}
		}
	}

	return nil
}

func renameTable(
	schema *LogicalSchema,
	table *LogicalTable,
	newName string,
) error {

	if _, exists := schema.Tables[newName]; exists {
		return fmt.Errorf("table %q already exists", newName)
	}

	oldName := table.Name

	delete(schema.Tables, oldName)

	table.Name = newName

	schema.Tables[newName] = table

	// Update every column.
	for _, column := range table.Columns {
		column.Table = newName
	}

	// Update every index.
	for _, index := range table.Indexes {
		index.Table = newName
	}

	// Update FK references across the entire schema.
	for _, tbl := range schema.Tables {

		for _, constraint := range tbl.Constraints {

			if constraint.Type != ir.ForeignKey {
				continue
			}

			if constraint.Reference == nil {
				continue
			}

			if constraint.Reference.Table == oldName {
				constraint.Reference.Table = newName
			}
		}

	}

	return nil
}

func dropConstraintByName(
	table *LogicalTable,
	name string,
) error {

	if _, exists := table.Constraints[name]; !exists {
		return fmt.Errorf(
			"constraint %q does not exist",
			name,
		)
	}

	delete(table.Constraints, name)

	return nil
}
