package parser

import "hyperfulcrum/internal/parser/provider/postgres"

func New() Parser {
	return postgres.New()
}
