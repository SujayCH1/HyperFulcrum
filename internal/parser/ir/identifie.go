package ir

// Identifier represents a qualified SQL name such as schema.name, optionally
// aliased (e.g. column references, constraint targets).
type Identifier struct {
	Schema string
	Name   string
	Alias  string
}

// Table represents a table/relation reference — used in DML FROM/JOIN
// clauses as well as DDL statement targets (CREATE/ALTER/DROP TABLE).
type Table struct {
	Schema string
	Name   string
	Alias  string
}
