package parser

import "hyperfulcrum/internal/parser/ir"

// Canonical definitions live in ir/errors.go (see the comment there for why).
// These are re-exported so existing call sites can keep writing
// parser.ErrEmptySQL, parser.ErrNilStatement, etc.
var (
	ErrNilStatement   = ir.ErrNilStatement
	ErrUnsupportedSQL = ir.ErrUnsupportedSQL
	ErrEmptySQL       = ir.ErrEmptySQL
	ErrInvalidAST     = ir.ErrInvalidAST
)
