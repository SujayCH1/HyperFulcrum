package postgres

import (
	"errors"
	"fmt"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v5"

	"hyperfulcrum/internal/parser/ir"
)

type PostgreSQLParser struct{}

func New() *PostgreSQLParser {
	return &PostgreSQLParser{}
}

func (p *PostgreSQLParser) Parse(sql string) (ir.Statement, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, ir.ErrEmptySQL
	}

	result, err := pg_query.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ir.ErrInvalidAST, err)
	}

	if len(result.Stmts) != 1 {
		return nil, ir.ErrMultipleSQL
	}

	return convertStatement(result.Stmts[0].Stmt, sql)
}

func (p *PostgreSQLParser) ParseBatch(sql string) (*ir.Batch, error) {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return nil, ir.ErrEmptySQL
	}

	result, err := pg_query.Parse(sql)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ir.ErrInvalidAST, err)
	}

	if len(result.Stmts) == 0 {
		return nil, ir.ErrEmptySQL
	}

	batch := &ir.Batch{
		Statements: make([]ir.Statement, 0, len(result.Stmts)),
	}

	for i, raw := range result.Stmts {
		statement, err := convertStatement(raw.Stmt, statementSQL(sql, result.Stmts, i))
		if err != nil {
			return nil, fmt.Errorf("statement %d: %w", i+1, err)
		}
		batch.Statements = append(batch.Statements, statement)
	}

	return batch, nil
}

func statementSQL(sql string, statements []*pg_query.RawStmt, index int) string {
	start := int(statements[index].StmtLocation)
	end := len(sql)

	if statements[index].StmtLen > 0 {
		end = start + int(statements[index].StmtLen)
	} else if index+1 < len(statements) {
		end = int(statements[index+1].StmtLocation)
	}

	if start < 0 || start > len(sql) || end < start || end > len(sql) {
		return ""
	}

	return strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(sql[start:end]), ";"))
}

func convertStatement(node *pg_query.Node, rawSQL string) (ir.Statement, error) {
	if node == nil {
		return nil, ir.ErrNilStatement
	}

	metadata := ir.Metadata{
		RawSQL:        rawSQL,
		SourceDialect: "postgres",
	}

	switch {
	case node.GetCreateStmt() != nil:
		return convertCreateTable(node.GetCreateStmt(), metadata)
	case node.GetAlterTableStmt() != nil:
		return convertAlterTable(node.GetAlterTableStmt(), metadata)
	case node.GetDropStmt() != nil:
		return convertDrop(node.GetDropStmt(), metadata)
	case node.GetIndexStmt() != nil:
		return convertCreateIndex(node.GetIndexStmt(), metadata)
	case node.GetRenameStmt() != nil:
		return convertRename(node.GetRenameStmt(), metadata)
	case node.GetTruncateStmt() != nil:
		return convertTruncate(node.GetTruncateStmt(), metadata)
	case node.GetSelectStmt() != nil:
		return convertRoute(node, ir.Select, metadata)
	case node.GetInsertStmt() != nil:
		return convertRoute(node, ir.Insert, metadata)
	case node.GetUpdateStmt() != nil:
		return convertRoute(node, ir.Update, metadata)
	case node.GetDeleteStmt() != nil:
		return convertRoute(node, ir.Delete, metadata)
	case node.GetMergeStmt() != nil:
		return convertRoute(node, ir.Merge, metadata)
	default:
		return nil, ir.ErrUnsupportedSQL
	}
}

func unsupported(value any) error {
	return errors.Join(ir.ErrUnsupportedSQL, fmt.Errorf("%v", value))
}
