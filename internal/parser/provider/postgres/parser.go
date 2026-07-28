package postgres

import (
	"strings"

	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

type PostgreSQLParser struct {
	opts postgresparser.ParseOptions
}

type Option func(*PostgreSQLParser)

func WithFieldComments() Option {
	return func(p *PostgreSQLParser) {
		p.opts.IncludeCreateTableFieldComments = true
	}
}

// New returns *PostgreSQLParser (a concrete type, not parser.Parser) — this
// package must never import hyperfulcrum/internal/parser, or you get the
// cycle you just hit. *PostgreSQLParser still satisfies parser.Parser
// structurally; factory.go (which imports both packages) does the
// interface conversion in `return postgres.New()`.
func New(opts ...Option) *PostgreSQLParser {
	p := &PostgreSQLParser{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *PostgreSQLParser) Parse(sql string) (ir.Statement, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, ir.ErrEmptySQL
	}

	result, err := postgresparser.ParseSQLWithOptions(sql, p.opts)
	if err != nil {
		return nil, translateErr(err)
	}

	return convertStatement(result)
}

func (p *PostgreSQLParser) ParseBatch(sql string) (*ir.Batch, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, ir.ErrEmptySQL
	}

	result, err := postgresparser.ParseSQLAllWithOptions(sql, p.opts)
	if err != nil {
		return nil, translateErr(err)
	}

	batch := &ir.Batch{}

	for _, stmtResult := range result.Statements {
		if stmtResult.Query == nil {
			continue
		}

		s, err := convertStatement(stmtResult.Query)
		if err != nil {
			return nil, err
		}

		batch.Statements = append(batch.Statements, s)
	}

	return batch, nil
}
