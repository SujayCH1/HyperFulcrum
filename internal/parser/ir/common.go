package ir

type StatementKind string

const (
	DDL     StatementKind = "DDL"
	DML     StatementKind = "DML"
	TCL     StatementKind = "TCL"
	DCL     StatementKind = "DCL"
	Utility StatementKind = "UTILITY"
)

type Command string

const (
	// DDL
	CreateTable Command = "CREATE_TABLE"
	AlterTable  Command = "ALTER_TABLE"
	DropTable   Command = "DROP_TABLE"

	CreateIndex Command = "CREATE_INDEX"
	DropIndex   Command = "DROP_INDEX"

	CreateSchema Command = "CREATE_SCHEMA"
	DropSchema   Command = "DROP_SCHEMA"

	Truncate Command = "TRUNCATE"
	Comment  Command = "COMMENT"

	// DML
	Select Command = "SELECT"
	Insert Command = "INSERT"
	Update Command = "UPDATE"
	Delete Command = "DELETE"
	Merge  Command = "MERGE"

	// TCL
	Begin    Command = "BEGIN"
	Commit   Command = "COMMIT"
	Rollback Command = "ROLLBACK"
)
