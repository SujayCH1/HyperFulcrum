package schema

import "hyperfulcrum/internal/parser/ir"

type LogicalSchema struct {
	ProjectID string
	Tables    map[string]*LogicalTable
}

type LogicalTable struct {
	Name string

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

func NewLogicalTable(name string) *LogicalTable {
	return &LogicalTable{
		Name:        name,
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

	table = NewLogicalTable(name)
	s.Tables[name] = table

	return table
}
