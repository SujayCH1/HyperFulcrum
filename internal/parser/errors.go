package parser

import "hyperfulcrum/internal/parser/ir"

var (
	ErrNilStatement   = ir.ErrNilStatement
	ErrUnsupportedSQL = ir.ErrUnsupportedSQL
	ErrEmptySQL       = ir.ErrEmptySQL
	ErrInvalidAST     = ir.ErrInvalidAST
	ErrMultipleSQL    = ir.ErrMultipleSQL
)
