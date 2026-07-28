package postgres

import (
	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

func convertStatement(pq *postgresparser.ParsedQuery) (ir.Statement, error) {
	if pq == nil {
		return nil, ir.ErrNilStatement
	}

	switch pq.Command {
	case postgresparser.QueryCommandSelect,
		postgresparser.QueryCommandInsert,
		postgresparser.QueryCommandUpdate,
		postgresparser.QueryCommandDelete:
		return convertDML(pq)

	case postgresparser.QueryCommandDDL:
		return convertDDL(pq)

	case postgresparser.QueryCommandMerge:
		return nil, ir.ErrUnsupportedSQL

	default:
		return nil, ir.ErrUnsupportedSQL
	}
}
