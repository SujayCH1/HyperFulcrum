package schema

import "hyperfulcrum/internal/parser/ir"

type LogicalSchema struct {
	ProjectID string
	Tables    map[string]*LogicalTable
}

type LogicalTable struct {
	Table ir.Table

	Columns map[string]*ir.Column

	// Key = constraint.Name
	Constraints map[string]*ir.Constraint

	Indexes map[string]*ir.Index
}

func NewLogicalSchema() *LogicalSchema {
	return &LogicalSchema{
		Tables: make(map[string]*LogicalTable),
	}
}

func NewLogicalTable(table ir.Table) *LogicalTable {
	return &LogicalTable{
		Table:       table,
		Columns:     make(map[string]*ir.Column),
		Constraints: make(map[string]*ir.Constraint),
		Indexes:     make(map[string]*ir.Index),
	}
}

func (s *LogicalSchema) EnsureTable(name string) *LogicalTable {
	table, ok := s.Tables[name]
	if ok {
		return table
	}

	table = NewLogicalTable(ir.NewTable(name))
	s.Tables[name] = table

	return table
}

func (s *LogicalSchema) Clone() *LogicalSchema {
	clone := NewLogicalSchema()
	clone.ProjectID = s.ProjectID

	for key, table := range s.Tables {
		clonedTable := NewLogicalTable(table.Table)

		for name, column := range table.Columns {
			value := *column
			value.DataType.Modifiers = append([]string(nil), column.DataType.Modifiers...)
			if column.DefaultValue != nil {
				defaultValue := *column.DefaultValue
				value.DefaultValue = &defaultValue
			}
			if column.Generated != nil {
				generated := *column.Generated
				value.Generated = &generated
			}
			clonedTable.Columns[name] = &value
		}

		for name, constraint := range table.Constraints {
			value := *constraint
			value.Columns = append([]string(nil), constraint.Columns...)
			if constraint.Reference != nil {
				reference := *constraint.Reference
				reference.Columns = append([]string(nil), constraint.Reference.Columns...)
				value.Reference = &reference
			}
			if constraint.Expression != nil {
				expression := *constraint.Expression
				value.Expression = &expression
			}
			clonedTable.Constraints[name] = &value
		}

		for name, index := range table.Indexes {
			value := *index
			value.Columns = append([]string(nil), index.Columns...)
			value.Expressions = append([]ir.Expression(nil), index.Expressions...)
			value.Include = append([]string(nil), index.Include...)
			if index.Predicate != nil {
				predicate := *index.Predicate
				value.Predicate = &predicate
			}
			clonedTable.Indexes[name] = &value
		}

		clone.Tables[key] = clonedTable
	}

	return clone
}
