package schema

type LogicalSchema struct {
	ProjectID string
	Tables    map[string]*Table
}

type Table struct {
	Columns map[string]*Column
	Fks     map[Fkkey]*Fk
}

type Column struct {
	Name     string
	DataType string
	Nullable bool
	IsPk     bool
}

type Fkkey struct {
	ChildColumn  string
	ParentTable  string
	ParentColumn string
}

type Fk struct {
	ChildTable   string
	ParentTable  string
	ChildColumn  string
	ParentColumn string
}

func NewLogicalSchema() *LogicalSchema {
	return &LogicalSchema{
		Tables: make(map[string]*Table),
	}
}
