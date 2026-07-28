package ir

// DMLStatement is the IR for SELECT/INSERT/UPDATE/DELETE statements.
type DMLStatement struct {
	Metadata

	Cmd        Command
	Tables     []Table
	Columns    []Column
	Conditions []Condition
	Joins      []Join
	GroupBy    []Expression
	OrderBy    []Expression
	Limit      *int
	Offset     *int
}

func (d *DMLStatement) Kind() StatementKind {
	return DML
}

func (d *DMLStatement) Command() Command {
	return d.Cmd
}
