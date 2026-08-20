package ir

type Column struct {
	Table string
	Name  string
	Alias string

	DataType DataType

	Nullable bool

	DefaultValue *Expression
	Identity     string
	Generated    *Expression
}

type Reference struct {
	Table    Table
	Columns  []string
	OnUpdate string
	OnDelete string
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
	Name              string
	Type              ConstraintType
	Columns           []string
	Reference         *Reference
	Expression        *Expression
	Deferrable        bool
	InitiallyDeferred bool
	NullsNotDistinct  bool
}

type Index struct {
	Name        string
	Table       Table
	Columns     []string
	Expressions []Expression
	Include     []string
	Predicate   *Expression
	Unique      bool
}

type AlterOperationType string

const (
	AddColumn AlterOperationType = "ADD_COLUMN"

	DropColumn AlterOperationType = "DROP_COLUMN"

	RenameColumn AlterOperationType = "RENAME_COLUMN"

	AddConstraint AlterOperationType = "ADD_CONSTRAINT"

	DropConstraint AlterOperationType = "DROP_CONSTRAINT"

	RenameTable AlterOperationType = "RENAME_TABLE"

	AlterColumnType AlterOperationType = "ALTER_COLUMN_TYPE"

	SetNotNull AlterOperationType = "SET_NOT_NULL"

	DropNotNull AlterOperationType = "DROP_NOT_NULL"

	SetDefault AlterOperationType = "SET_DEFAULT"

	DropDefault AlterOperationType = "DROP_DEFAULT"
)

type AlterOperation struct {
	Type       AlterOperationType
	Column     *Column
	Constraint *Constraint
	OldName    string
	NewName    string
	Expression *Expression
	DataType   *DataType
}

type DDLStatement struct {
	Metadata

	Cmd             Command
	Table           Table
	Tables          []Table
	Columns         []Column
	Constraints     []Constraint
	Indexes         []Index
	AlterOperations []AlterOperation
	IfExists        bool
	IfNotExists     bool
	Cascade         bool
}

func (d *DDLStatement) Kind() StatementKind {
	return DDL
}

func (d *DDLStatement) Command() Command {
	return d.Cmd
}
