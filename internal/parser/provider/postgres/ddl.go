package postgres

import (
	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

// convertDDL converts a ParsedQuery's DDLActions into an ir.DDLStatement.
// One SQL statement can produce several DDLActions (e.g. ALTER TABLE with
// several sub-clauses). The first action establishes Cmd/Table identity;
// CREATE_TABLE-only actions populate the top-level Columns/Constraints,
// while ALTER-style actions fold into AlterOperations instead.
func convertDDL(pq *postgresparser.ParsedQuery) (*ir.DDLStatement, error) {
	if len(pq.DDLActions) == 0 {
		return &ir.DDLStatement{}, nil
	}

	first := pq.DDLActions[0]

	stmt := &ir.DDLStatement{
		Cmd:   mapDDLCommand(first.Type),
		Table: newTable(first.Schema, first.ObjectName, ""),
	}

	for _, action := range pq.DDLActions {
		switch action.Type {
		case postgresparser.DDLCreateTable:
			for _, col := range action.ColumnDetails {
				stmt.Columns = append(stmt.Columns, ir.Column{
					Name:         col.Name,
					DataType:     parseDataType(col.Type),
					Nullable:     col.Nullable,
					DefaultValue: nonEmptyExpr(col.Default),
				})
			}
			if action.Constraints != nil {
				stmt.Constraints = append(stmt.Constraints, convertConstraints(action.Constraints)...)
			}

		case postgresparser.DDLCreateIndex:
			stmt.Indexes = append(stmt.Indexes, ir.Index{
				Name: action.ObjectName,
				// NOTE: verify against your own test output — DDLAction
				// doesn't have an explicit "owning table" field for CREATE
				// INDEX; Target's doc comment says "generic fully-qualified
				// target path for comment-like actions", so this may need
				// adjusting once you run a real CREATE INDEX through it.
				Table:   action.Target,
				Columns: action.Columns,
				Unique:  containsFlag(action.Flags, "UNIQUE"),
			})

		case postgresparser.DDLAlterTable, postgresparser.DDLDropColumn:
			stmt.AlterOperations = append(stmt.AlterOperations, convertAlterOperation(action))
		}
	}

	return stmt, nil
}

func convertConstraints(c *postgresparser.DDLConstraints) []ir.Constraint {
	var out []ir.Constraint

	if c.PrimaryKey != nil {
		out = append(out, ir.Constraint{
			Name:    c.PrimaryKey.ConstraintName,
			Type:    ir.PrimaryKey,
			Columns: c.PrimaryKey.Columns,
		})
	}

	for _, fk := range c.ForeignKeys {
		out = append(out, ir.Constraint{
			Name:    fk.ConstraintName,
			Type:    ir.ForeignKey,
			Columns: fk.Columns,
			Reference: &ir.Reference{
				Table:   fk.ReferencesTable,
				Columns: fk.ReferencesColumns,
			},
		})
	}

	for _, uq := range c.UniqueKeys {
		out = append(out, ir.Constraint{
			Name:    uq.ConstraintName,
			Type:    ir.Unique,
			Columns: uq.Columns,
		})
	}

	for _, chk := range c.CheckConstraints {
		out = append(out, ir.Constraint{
			Name:       chk.ConstraintName,
			Type:       ir.Check,
			Expression: &ir.Expression{Raw: chk.Expression},
		})
	}

	return out
}

// convertAlterOperation is a best-effort mapping. ValkDB's DDLActionType
// enum currently has no RENAME_COLUMN/RENAME_TABLE variant, so
// ir.RenameColumn/ir.RenameTable can't be populated from it yet — that's a
// gap in the upstream library, not something fixable here.
func convertAlterOperation(action postgresparser.DDLAction) ir.AlterOperation {
	switch action.Type {
	case postgresparser.DDLDropColumn:
		return ir.AlterOperation{
			Type:    ir.DropColumn,
			OldName: firstOrEmpty(action.Columns),
		}

	case postgresparser.DDLAlterTable:
		if len(action.ColumnDetails) > 0 {
			col := action.ColumnDetails[0]
			return ir.AlterOperation{
				Type: ir.AddColumn,
				Column: &ir.Column{
					Name:         col.Name,
					DataType:     parseDataType(col.Type),
					Nullable:     col.Nullable,
					DefaultValue: nonEmptyExpr(col.Default),
				},
			}
		}
		if action.Constraints != nil {
			cs := convertConstraints(action.Constraints)
			if len(cs) > 0 {
				return ir.AlterOperation{
					Type:       ir.AddConstraint,
					Constraint: &cs[0],
				}
			}
		}
	}

	return ir.AlterOperation{}
}
