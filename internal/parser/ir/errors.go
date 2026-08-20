package ir

import "errors"

var (
	ErrNilStatement   = errors.New("nil statement")
	ErrUnsupportedSQL = errors.New("unsupported sql statement")
	ErrEmptySQL       = errors.New("empty sql")
	ErrInvalidAST     = errors.New("invalid parser ast")
	ErrMultipleSQL    = errors.New("multiple sql statements")
)
