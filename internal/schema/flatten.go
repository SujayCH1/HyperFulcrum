package schema

import (
	"sort"

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
	tableNames := make([]string, 0, len(schema.Tables))
	for tableName := range schema.Tables {
		tableNames = append(tableNames, tableName)
	}
	sort.Strings(tableNames)

	for _, tableName := range tableNames {
		table := schema.Tables[tableName]

		primaryKeys := make(map[string]struct{})

		for _, constraint := range table.Constraints {

			if constraint.Type != ir.PrimaryKey {
				continue
			}

			for _, column := range constraint.Columns {
				primaryKeys[column] = struct{}{}
			}
		}

		columnNames := make([]string, 0, len(table.Columns))
		for columnName := range table.Columns {
			columnNames = append(columnNames, columnName)
		}
		sort.Strings(columnNames)

		for _, columnName := range columnNames {
			column := table.Columns[columnName]

			_, isPrimaryKey := primaryKeys[column.Name]

			columns = append(columns, repository.Column{
				ProjectID:  schema.ProjectID,
				TableName:  tableName,
				ColumnName: column.Name,
				DataType:   column.DataType.String(),
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
	tableNames := make([]string, 0, len(schema.Tables))
	for tableName := range schema.Tables {
		tableNames = append(tableNames, tableName)
	}
	sort.Strings(tableNames)

	for _, tableName := range tableNames {
		table := schema.Tables[tableName]

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
					ParentTable:  constraint.Reference.Table.Key(),
					ParentColumn: constraint.Reference.Columns[i],
					ChildTable:   tableName,
					ChildColumn:  childColumn,
				})
			}
		}
	}

	sort.Slice(edges, func(i, j int) bool {
		left := edges[i]
		right := edges[j]
		if left.ChildTable != right.ChildTable {
			return left.ChildTable < right.ChildTable
		}
		if left.ChildColumn != right.ChildColumn {
			return left.ChildColumn < right.ChildColumn
		}
		if left.ParentTable != right.ParentTable {
			return left.ParentTable < right.ParentTable
		}
		return left.ParentColumn < right.ParentColumn
	})

	return edges
}
