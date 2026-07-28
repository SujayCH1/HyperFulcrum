package postgres

import (
	"errors"

	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

func translateErr(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, postgresparser.ErrNoStatements):
		return ir.ErrEmptySQL
	case errors.Is(err, postgresparser.ErrMultipleStatements):
		return ir.ErrUnsupportedSQL
	case errors.Is(err, postgresparser.ErrNilContext):
		return ir.ErrInvalidAST
	}

	var parseErrs *postgresparser.ParseErrors
	if errors.As(err, &parseErrs) {
		return errors.Join(ir.ErrInvalidAST, err)
	}

	return err
}
