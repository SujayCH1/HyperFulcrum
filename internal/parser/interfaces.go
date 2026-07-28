package parser

import "hyperfulcrum/internal/parser/ir"

type Parser interface {
	Parse(sql string) (ir.Statement, error)
	ParseBatch(sql string) (*ir.Batch, error)
}
