package ir

type Column struct {
	Table string
	Name  string
	Alias string

	DataType DataType

	Nullable bool

	DefaultValue *Expression
}

type Reference struct {
	Table   string
	Columns []string
}

type ConstraintType string

const (
	PrimaryKey ConstraintType = "PRIMARY_KEY"
	ForeignKey ConstraintType = "FOREIGN_KEY"
	Unique     ConstraintType = "UNIQUE"
	Check      ConstraintType = "CHECK"
	Default    ConstraintType = "DEFAULT"
)

type Constraint struct {
	Name       string
	Type       ConstraintType
	Columns    []string
	Reference  *Reference
	Expression *Expression
}

type Index struct {
	Name    string
	Table   string
	Columns []string
	Unique  bool
}

type AlterOperationType string

const (
	AddColumn AlterOperationType = "ADD_COLUMN"

	DropColumn AlterOperationType = "DROP_COLUMN"

	RenameColumn AlterOperationType = "RENAME_COLUMN"

	AddConstraint AlterOperationType = "ADD_CONSTRAINT"

	DropConstraint AlterOperationType = "DROP_CONSTRAINT"

	RenameTable AlterOperationType = "RENAME_TABLE"
)

type AlterOperation struct {
	Type       AlterOperationType
	Column     *Column
	Constraint *Constraint
	OldName    string
	NewName    string
}

type DDLStatement struct {
	Metadata

	Cmd             Command
	Table           Table
	Columns         []Column
	Constraints     []Constraint
	Indexes         []Index
	AlterOperations []AlterOperation
}

func (d *DDLStatement) Kind() StatementKind {
	return DDL
}

func (d *DDLStatement) Command() Command {
	return d.Cmd
}
