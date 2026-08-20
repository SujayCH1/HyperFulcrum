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

	column.Table = table.Table.Key()

	c := column

	table.Columns[c.Name] = &c

	return nil
}

func dropColumn(
	schema *LogicalSchema,
	table *LogicalTable,
	columnName string,
) error {

	if _, exists := table.Columns[columnName]; !exists {
		return fmt.Errorf("column %q does not exist", columnName)
	}

	delete(table.Columns, columnName)

	for name, constraint := range table.Constraints {
		if containsColumn(constraint.Columns, columnName) {
			delete(table.Constraints, name)
		}
	}

	for name, index := range table.Indexes {
		if containsColumn(index.Columns, columnName) {
			delete(table.Indexes, name)
		}
	}

	for _, candidate := range schema.Tables {
		for name, constraint := range candidate.Constraints {
			if constraint.Type != ir.ForeignKey || constraint.Reference == nil {
				continue
			}

			if constraint.Reference.Table.Key() == table.Table.Key() &&
				containsColumn(constraint.Reference.Columns, columnName) {
				delete(candidate.Constraints, name)
			}
		}
	}

	return nil
}

func addConstraint(
	table *LogicalTable,
	constraint ir.Constraint,
) error {

	ensureConstraintName(table.Table.Key(), &constraint)

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
	schema *LogicalSchema,
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

	for _, index := range table.Indexes {
		for i, columnName := range index.Columns {
			if columnName == oldName {
				index.Columns[i] = newName
			}
		}
	}

	for _, candidate := range schema.Tables {
		for _, constraint := range candidate.Constraints {
			if constraint.Type != ir.ForeignKey || constraint.Reference == nil {
				continue
			}

			if constraint.Reference.Table.Key() != table.Table.Key() {
				continue
			}

			for i, columnName := range constraint.Reference.Columns {
				if columnName == oldName {
					constraint.Reference.Columns[i] = newName
				}
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

	oldName := table.Table.Key()
	newTable := table.Table
	newTable.Name = newName
	newKey := newTable.Key()

	if _, exists := schema.Tables[newKey]; exists {
		return fmt.Errorf("table %q already exists", newKey)
	}

	delete(schema.Tables, oldName)

	table.Table = newTable

	schema.Tables[newKey] = table

	// Update every column.
	for _, column := range table.Columns {
		column.Table = newKey
	}

	// Update every index.
	for _, index := range table.Indexes {
		index.Table = newTable
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

			if constraint.Reference.Table.Key() == oldName {
				constraint.Reference.Table = newTable
			}
		}

	}

	return nil
}

func containsColumn(columns []string, name string) bool {
	for _, column := range columns {
		if column == name {
			return true
		}
	}

	return false
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
