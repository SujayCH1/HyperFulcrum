package ir

// Expression is a raw SQL expression fragment. Providers are responsible for
// producing (or best-effort reconstructing) this text from their own AST/IR.
type Expression struct {
	Raw string
}

type JoinType string

const (
	InnerJoin JoinType = "INNER"
	LeftJoin  JoinType = "LEFT"
	RightJoin JoinType = "RIGHT"
	FullJoin  JoinType = "FULL"
	CrossJoin JoinType = "CROSS"
)

// Condition models a single predicate, e.g. from a WHERE clause or JOIN ON.
type Condition struct {
	Left     Expression
	Operator string
	Right    Expression
}

// Join models a single join between two tables.
type Join struct {
	Type      JoinType
	Left      Table
	Right     Table
	Condition Expression
}
