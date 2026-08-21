package schema

import (
	"hyperfulcrum/internal/parser/ir"
	"hyperfulcrum/internal/repository"
)

func BuildLogicalSchemaFromMetadata(
	projectID string,
	columns []repository.Column,
	fkEdges []repository.FkEdges,
) (*LogicalSchema, error) {

	schema := NewLogicalSchema()
	schema.ProjectID = projectID

	if err := addColumns(schema, columns); err != nil {
		return nil, err
	}

	if err := addForeignKeys(schema, fkEdges); err != nil {
		return nil, err
	}

	return schema, nil
}

func addColumns(
	schema *LogicalSchema,
	columns []repository.Column,
) error {

	for _, column := range columns {

		table := schema.EnsureTable(column.TableName)

		if err := addColumn(table, ir.Column{
			Name:     column.ColumnName,
			DataType: ir.NewDataType(column.DataType),
			Nullable: column.IsNullable,
		}); err != nil {
			return err
		}

		if column.IsPk {

			if err := addConstraint(table, ir.Constraint{
				Type:    ir.PrimaryKey,
				Columns: []string{column.ColumnName},
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func addForeignKeys(
	schema *LogicalSchema,
	fkEdges []repository.FkEdges,
) error {

	for _, edge := range fkEdges {

		table := schema.EnsureTable(edge.ChildTable)

		if err := addConstraint(table, ir.Constraint{
			Type:    ir.ForeignKey,
			Columns: []string{edge.ChildColumn},
			Reference: &ir.Reference{
				Table: ir.NewTable(edge.ParentTable),
				Columns: []string{
					edge.ParentColumn,
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
