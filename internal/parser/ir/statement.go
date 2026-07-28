package ir

type Statement interface {
	Kind() StatementKind

	Command() Command
}
