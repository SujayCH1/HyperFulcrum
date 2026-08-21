package parser

import (
	"fmt"

	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/parser/provider/postgres"
)

var defaultParser = postgres.New()

func New() Parser {
	return defaultParser
}

func Parse(sql string) (ir.Statement, error) {
	return defaultParser.Parse(sql)
}

func ParseBatch(sql string) (*ir.Batch, error) {
	return defaultParser.ParseBatch(sql)
}

func ParseDDLBatch(sql string) (*ir.Batch, error) {
	batch, err := defaultParser.ParseBatch(sql)
	if err != nil {
		return nil, err
	}

	for _, statement := range batch.Statements {
		if statement.Kind() != ir.DDL {
			return nil, fmt.Errorf("expected DDL statement, got %s", statement.Kind())
		}
	}

	return batch, nil
}

func ParseDDL(sql string) (*ir.DDLStatement, error) {
	statement, err := defaultParser.Parse(sql)
	if err != nil {
		return nil, err
	}

	ddl, ok := statement.(*ir.DDLStatement)
	if !ok {
		return nil, fmt.Errorf("expected DDL statement, got %s", statement.Kind())
	}

	return ddl, nil
}

func ParseRoute(sql string) (*ir.RouteStatement, error) {
	statement, err := defaultParser.Parse(sql)
	if err != nil {
		return nil, err
	}

	route, ok := statement.(*ir.RouteStatement)
	if !ok {
		return nil, fmt.Errorf("expected DML statement, got %s", statement.Kind())
	}

	return route, nil
}
