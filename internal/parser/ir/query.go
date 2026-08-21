package ir

type RouteValueKind string

const (
	UnknownValue   RouteValueKind = "UNKNOWN"
	LiteralValue   RouteValueKind = "LITERAL"
	ParameterValue RouteValueKind = "PARAMETER"
	ColumnValue    RouteValueKind = "COLUMN"
)

type RouteValue struct {
	Kind      RouteValueKind
	Value     string
	Parameter int
	Table     string
	Column    string
}

type RoutePredicate struct {
	Table    string
	Column   string
	Operator string
	Value    RouteValue
}

type RouteAssignment struct {
	Column string
	Value  RouteValue
}

type RouteStatement struct {
	Metadata

	Cmd             Command
	Tables          []Table
	Predicates      []RoutePredicate
	Assignments     []RouteAssignment
	InsertColumns   []string
	InsertRows      [][]RouteValue
	RoutingComplete bool
}

func (r *RouteStatement) Kind() StatementKind {
	return DML
}

func (r *RouteStatement) Command() Command {
	return r.Cmd
}
