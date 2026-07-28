package ir

import "errors"

// Sentinel errors shared by parser implementations. They live in ir rather
// than in the top-level parser package specifically so provider packages
// (e.g. provider/postgres) can return them without importing parser —
// parser imports providers via factory.go, so a provider importing parser
// back would create an import cycle.
var (
	ErrNilStatement   = errors.New("nil statement")
	ErrUnsupportedSQL = errors.New("unsupported sql statement")
	ErrEmptySQL       = errors.New("empty sql")
	ErrInvalidAST     = errors.New("invalid parser ast")
)
