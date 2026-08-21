package schema

import (
	"fmt"
	"hash/fnv"
	"strings"

	"hyperfulcrum/internal/parser/ir"
)

func ensureConstraintName(
	tableName string,
	constraint *ir.Constraint,
) {

	if constraint.Name != "" {
		return
	}

	switch constraint.Type {

	case ir.PrimaryKey:

		constraint.Name = fmt.Sprintf(
			"pk_%s_%s",
			tableName,
			strings.Join(constraint.Columns, "_"),
		)

	case ir.ForeignKey:

		ref := ""

		if constraint.Reference != nil {
			ref = constraint.Reference.Table.Key()
		}

		constraint.Name = fmt.Sprintf(
			"fk_%s_%s_%s",
			tableName,
			strings.Join(constraint.Columns, "_"),
			ref,
		)

	case ir.Unique:

		constraint.Name = fmt.Sprintf(
			"uq_%s_%s",
			tableName,
			strings.Join(constraint.Columns, "_"),
		)

	case ir.Check:
		hash := fnv.New32a()
		if constraint.Expression != nil {
			_, _ = hash.Write([]byte(constraint.Expression.Raw))
		}

		constraint.Name = fmt.Sprintf(
			"chk_%s_%x",
			tableName,
			hash.Sum32(),
		)

	case ir.Default:

		constraint.Name = fmt.Sprintf(
			"default_%s_%s",
			tableName,
			strings.Join(constraint.Columns, "_"),
		)
	}
}
