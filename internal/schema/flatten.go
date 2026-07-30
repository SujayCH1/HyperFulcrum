package schema

import (
	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/repository"
)

func FlattenLogicalSchema(
	schema *LogicalSchema,
) (
	[]repository.Column,
	[]repository.FkEdges,
) {

	return flattenColumns(schema), flattenFKEdges(schema)
}

func flattenColumns(
	schema *LogicalSchema,
) []repository.Column {

	var columns []repository.Column

	for tableName, table := range schema.Tables {

		primaryKeys := make(map[string]struct{})

		for _, constraint := range table.Constraints {

			if constraint.Type != ir.PrimaryKey {
				continue
			}

			for _, column := range constraint.Columns {
				primaryKeys[column] = struct{}{}
			}
		}

		for _, column := range table.Columns {

			_, isPrimaryKey := primaryKeys[column.Name]

			columns = append(columns, repository.Column{
				ProjectID:  schema.ProjectID,
				TableName:  tableName,
				ColumnName: column.Name,
				DataType:   column.DataType.Name,
				IsNullable: column.Nullable,
				IsPk:       isPrimaryKey,
			})
		}
	}

	return columns
}

func flattenFKEdges(
	schema *LogicalSchema,
) []repository.FkEdges {

	var edges []repository.FkEdges

	for tableName, table := range schema.Tables {

		for _, constraint := range table.Constraints {

			if constraint.Type != ir.ForeignKey {
				continue
			}

			if constraint.Reference == nil {
				continue
			}

			for i, childColumn := range constraint.Columns {

				if i >= len(constraint.Reference.Columns) {
					break
				}

				edges = append(edges, repository.FkEdges{
					ProjectId:    schema.ProjectID,
					ParentTable:  constraint.Reference.Table,
					ParentColumn: constraint.Reference.Columns[i],
					ChildTable:   tableName,
					ChildColumn:  childColumn,
				})
			}
		}
	}

	return edges
}
