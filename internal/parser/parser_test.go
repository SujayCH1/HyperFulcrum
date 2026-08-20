package parser_test

import (
	"errors"
	"testing"

	"hyperfulcrum/internal/parser"
	"hyperfulcrum/internal/parser/ir"
)

func TestParseCreateTable(t *testing.T) {
	statement, err := parser.ParseDDL(`
		CREATE TABLE sales.orders (
			id bigint PRIMARY KEY,
			customer_id bigint NOT NULL,
			total numeric(10, 2) DEFAULT 0,
			tags varchar(40)[],
			CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id)
				REFERENCES public.customers (id) ON UPDATE RESTRICT ON DELETE CASCADE,
			CONSTRAINT orders_total_check CHECK (total >= 0)
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	if statement.Command() != ir.CreateTable {
		t.Fatalf("expected CREATE TABLE, got %s", statement.Command())
	}
	if statement.Table.Key() != "sales.orders" {
		t.Fatalf("expected sales.orders, got %s", statement.Table.Key())
	}
	if len(statement.Columns) != 4 {
		t.Fatalf("expected 4 columns, got %d", len(statement.Columns))
	}
	if statement.Columns[2].DataType.String() != "numeric(10,2)" {
		t.Fatalf("unexpected numeric type: %s", statement.Columns[2].DataType.String())
	}
	if statement.Columns[3].DataType.String() != "character varying(40)[]" {
		t.Fatalf("unexpected array type: %s", statement.Columns[3].DataType.String())
	}
	if statement.Columns[2].DefaultValue == nil || statement.Columns[2].DefaultValue.Raw != "0" {
		t.Fatalf("unexpected default: %#v", statement.Columns[2].DefaultValue)
	}

	var foreignKey *ir.Constraint
	for i := range statement.Constraints {
		if statement.Constraints[i].Type == ir.ForeignKey {
			foreignKey = &statement.Constraints[i]
			break
		}
	}
	if foreignKey == nil || foreignKey.Reference == nil {
		t.Fatal("foreign key was not parsed")
	}
	if foreignKey.Reference.Table.Key() != "customers" {
		t.Fatalf("unexpected reference table: %s", foreignKey.Reference.Table.Key())
	}
	if foreignKey.Reference.OnUpdate != "RESTRICT" || foreignKey.Reference.OnDelete != "CASCADE" {
		t.Fatalf("unexpected foreign key actions: %#v", foreignKey.Reference)
	}
}

func TestParseAlterTable(t *testing.T) {
	tests := []struct {
		sql       string
		operation ir.AlterOperationType
	}{
		{"ALTER TABLE users ADD COLUMN email text NOT NULL", ir.AddColumn},
		{"ALTER TABLE users DROP COLUMN email", ir.DropColumn},
		{"ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE (email)", ir.AddConstraint},
		{"ALTER TABLE users DROP CONSTRAINT users_email_unique", ir.DropConstraint},
		{"ALTER TABLE users ALTER COLUMN email TYPE varchar(120)", ir.AlterColumnType},
		{"ALTER TABLE users ALTER COLUMN email SET NOT NULL", ir.SetNotNull},
		{"ALTER TABLE users ALTER COLUMN email DROP NOT NULL", ir.DropNotNull},
		{"ALTER TABLE users ALTER COLUMN email SET DEFAULT 'unknown'", ir.SetDefault},
		{"ALTER TABLE users ALTER COLUMN email DROP DEFAULT", ir.DropDefault},
		{"ALTER TABLE users RENAME COLUMN email TO address", ir.RenameColumn},
		{"ALTER TABLE users RENAME TO accounts", ir.RenameTable},
	}

	for _, test := range tests {
		t.Run(string(test.operation), func(t *testing.T) {
			statement, err := parser.ParseDDL(test.sql)
			if err != nil {
				t.Fatal(err)
			}
			if len(statement.AlterOperations) == 0 {
				t.Fatal("expected an alter operation")
			}
			if statement.AlterOperations[0].Type != test.operation {
				t.Fatalf("expected %s, got %s", test.operation, statement.AlterOperations[0].Type)
			}
		})
	}
}

func TestParseIndexAndMultipleDrop(t *testing.T) {
	index, err := parser.ParseDDL("CREATE UNIQUE INDEX users_email_idx ON public.users (email)")
	if err != nil {
		t.Fatal(err)
	}
	if index.Table.Key() != "users" || len(index.Indexes) != 1 {
		t.Fatalf("unexpected index statement: %#v", index)
	}
	if index.Indexes[0].Name != "users_email_idx" || !index.Indexes[0].Unique {
		t.Fatalf("unexpected index: %#v", index.Indexes[0])
	}

	drop, err := parser.ParseDDL("DROP TABLE users, sales.orders")
	if err != nil {
		t.Fatal(err)
	}
	if len(drop.Tables) != 2 || drop.Tables[1].Key() != "sales.orders" {
		t.Fatalf("unexpected drop targets: %#v", drop.Tables)
	}
}

func TestParseIdentityAndGeneratedColumns(t *testing.T) {
	statement, err := parser.ParseDDL(`
		CREATE TABLE inventory (
			id bigint GENERATED ALWAYS AS IDENTITY,
			quantity integer,
			price numeric(10, 2),
			total numeric(10, 2) GENERATED ALWAYS AS (quantity * price) STORED,
			matrix integer[][]
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	if statement.Columns[0].Identity != "ALWAYS" {
		t.Fatalf("unexpected identity: %s", statement.Columns[0].Identity)
	}
	if statement.Columns[3].Generated == nil || statement.Columns[3].Generated.Raw == "" {
		t.Fatal("generated expression was not parsed")
	}
	if statement.Columns[4].DataType.ArrayDimensions != 2 {
		t.Fatalf("unexpected array dimensions: %d", statement.Columns[4].DataType.ArrayDimensions)
	}
}

func TestParseIsStrict(t *testing.T) {
	_, err := parser.Parse("SELECT 1; SELECT 2")
	if !errors.Is(err, parser.ErrMultipleSQL) {
		t.Fatalf("expected multiple statement error, got %v", err)
	}

	_, err = parser.ParseBatch("CREATE TABLE one (id bigint); definitely invalid; CREATE TABLE two (id bigint)")
	if !errors.Is(err, parser.ErrInvalidAST) {
		t.Fatalf("expected invalid AST error, got %v", err)
	}

	_, err = parser.ParseDDL("CREATE VIEW active_users AS SELECT * FROM users")
	if !errors.Is(err, parser.ErrUnsupportedSQL) {
		t.Fatalf("expected unsupported SQL error, got %v", err)
	}

	_, err = parser.ParseDDL("ALTER TABLE users ENABLE ROW LEVEL SECURITY")
	if !errors.Is(err, parser.ErrUnsupportedSQL) {
		t.Fatalf("expected unsupported alter error, got %v", err)
	}
}

func TestParseRoute(t *testing.T) {
	statement, err := parser.ParseRoute(`
		SELECT o.id
		FROM orders o
		JOIN customers c ON c.id = o.customer_id
		WHERE o.tenant_id = $1 AND c.region = 'west=1'
	`)
	if err != nil {
		t.Fatal(err)
	}

	if !statement.RoutingComplete {
		t.Fatal("expected routing extraction to be complete")
	}
	if len(statement.Tables) != 2 {
		t.Fatalf("expected 2 tables, got %d", len(statement.Tables))
	}
	if len(statement.Predicates) != 3 {
		t.Fatalf("expected 3 predicates, got %d", len(statement.Predicates))
	}

	predicate := statement.Predicates[1]
	if predicate.Table != "o" || predicate.Column != "tenant_id" {
		t.Fatalf("unexpected shard predicate: %#v", predicate)
	}
	if predicate.Value.Kind != ir.ParameterValue || predicate.Value.Parameter != 1 {
		t.Fatalf("unexpected parameter: %#v", predicate.Value)
	}

	literal := statement.Predicates[2].Value
	if literal.Kind != ir.LiteralValue || literal.Value != "west=1" {
		t.Fatalf("unexpected literal: %#v", literal)
	}
}

func TestParseInsertRoute(t *testing.T) {
	statement, err := parser.ParseRoute(`
		INSERT INTO orders (tenant_id, total)
		VALUES ($1, 10.50), ($2, 20)
	`)
	if err != nil {
		t.Fatal(err)
	}

	if !statement.RoutingComplete {
		t.Fatal("expected routing extraction to be complete")
	}
	if len(statement.InsertColumns) != 2 || len(statement.InsertRows) != 2 {
		t.Fatalf("unexpected insert routing data: %#v", statement)
	}
	if statement.InsertRows[1][0].Parameter != 2 {
		t.Fatalf("unexpected second row: %#v", statement.InsertRows[1])
	}
}

func TestParseUpdateRoute(t *testing.T) {
	statement, err := parser.ParseRoute(`
		UPDATE orders
		SET tenant_id = $2, total = -10
		WHERE tenant_id = $1
	`)
	if err != nil {
		t.Fatal(err)
	}

	if !statement.RoutingComplete {
		t.Fatal("expected routing extraction to be complete")
	}
	if len(statement.Assignments) != 2 {
		t.Fatalf("unexpected assignments: %#v", statement.Assignments)
	}
	if statement.Assignments[0].Value.Parameter != 2 {
		t.Fatalf("unexpected shard assignment: %#v", statement.Assignments[0])
	}
	if statement.Assignments[1].Value.Value != "-10" {
		t.Fatalf("unexpected negative value: %#v", statement.Assignments[1])
	}
}

func TestUnsafeRoutesAreIncomplete(t *testing.T) {
	tests := []string{
		"SELECT * FROM users WHERE tenant_id = $1 OR tenant_id = $2",
		"SELECT * FROM users WHERE tenant_id = lower($1)",
		"WITH selected AS (SELECT * FROM users WHERE tenant_id = $1) SELECT * FROM selected",
		"SELECT * FROM users WHERE tenant_id IN (SELECT tenant_id FROM accounts)",
	}

	for _, sql := range tests {
		statement, err := parser.ParseRoute(sql)
		if err != nil {
			t.Fatal(err)
		}
		if statement.RoutingComplete {
			t.Fatalf("expected incomplete routing for %q", sql)
		}
	}
}

func TestSubqueryTablesArePreserved(t *testing.T) {
	statement, err := parser.ParseRoute(`
		SELECT * FROM users
		WHERE tenant_id IN (SELECT tenant_id FROM accounts)
	`)
	if err != nil {
		t.Fatal(err)
	}

	if statement.RoutingComplete {
		t.Fatal("subquery routing must be conservative")
	}
	if len(statement.Tables) != 2 || statement.Tables[1].Name != "accounts" {
		t.Fatalf("subquery table was not preserved: %#v", statement.Tables)
	}
}

func BenchmarkParseRoute(b *testing.B) {
	for b.Loop() {
		_, err := parser.ParseRoute("SELECT * FROM orders WHERE tenant_id = $1 AND id = $2")
		if err != nil {
			b.Fatal(err)
		}
	}
}
