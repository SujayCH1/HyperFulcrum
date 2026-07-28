package postgres

import (
	"github.com/valkdb/postgresparser"

	"hyperfulcrum/internal/parser/ir"
)

func mapJoinType(kind string) ir.JoinType {
	switch kind {
	case "LEFT":
		return ir.LeftJoin
	case "RIGHT":
		return ir.RightJoin
	case "FULL":
		return ir.FullJoin
	case "CROSS":
		return ir.CrossJoin
	default: // "INNER", "NATURAL", "LATERAL", ""
		return ir.InnerJoin
	}
}

func mapDMLCommand(cmd postgresparser.QueryCommand) ir.Command {
	switch cmd {
	case postgresparser.QueryCommandSelect:
		return ir.Select
	case postgresparser.QueryCommandInsert:
		return ir.Insert
	case postgresparser.QueryCommandUpdate:
		return ir.Update
	case postgresparser.QueryCommandDelete:
		return ir.Delete
	default:
		return ir.Select
	}
}

func mapDDLCommand(t postgresparser.DDLActionType) ir.Command {
	switch t {
	case postgresparser.DDLCreateTable:
		return ir.CreateTable
	case postgresparser.DDLAlterTable, postgresparser.DDLDropColumn:
		return ir.AlterTable
	case postgresparser.DDLDropTable:
		return ir.DropTable
	case postgresparser.DDLCreateIndex:
		return ir.CreateIndex
	case postgresparser.DDLDropIndex:
		return ir.DropIndex
	case postgresparser.DDLTruncate:
		return ir.Truncate
	case postgresparser.DDLComment:
		return ir.Comment
	default:
		return ir.AlterTable
	}
}
