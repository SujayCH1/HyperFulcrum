package schema_test

import (
	"testing"

	"hyperfulcrum/internal/parser"
	"hyperfulcrum/internal/schema"
)

func TestApplyParsedSchema(t *testing.T) {
	batch, err := parser.ParseBatch(`
		CREATE TABLE customers (id bigint PRIMARY KEY);
		CREATE TABLE orders (
			id bigint PRIMARY KEY,
			customer_id bigint,
			CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id) REFERENCES customers (id)
		);
		CREATE INDEX orders_customer_idx ON orders (customer_id);
		ALTER TABLE orders ADD COLUMN total numeric(10, 2) DEFAULT 0;
		ALTER TABLE orders ALTER COLUMN total SET NOT NULL;
	`)
	if err != nil {
		t.Fatal(err)
	}

	logicalSchema := schema.NewLogicalSchema()
	if err := schema.ApplyBatch(logicalSchema, batch); err != nil {
		t.Fatal(err)
	}

	orders := logicalSchema.Tables["orders"]
	if orders == nil {
		t.Fatal("orders table was not created")
	}
	if orders.Columns["total"].DataType.String() != "numeric(10,2)" {
		t.Fatalf("unexpected data type: %s", orders.Columns["total"].DataType.String())
	}
	if orders.Columns["total"].Nullable {
		t.Fatal("total should not be nullable")
	}
	if _, exists := orders.Indexes["orders_customer_idx"]; !exists {
		t.Fatal("index was not added")
	}
}

func TestRenameAndDropUpdateRelationships(t *testing.T) {
	batch, err := parser.ParseBatch(`
		CREATE TABLE customers (id bigint PRIMARY KEY);
		CREATE TABLE orders (
			id bigint PRIMARY KEY,
			customer_id bigint,
			CONSTRAINT orders_customer_fk FOREIGN KEY (customer_id) REFERENCES customers (id)
		);
		ALTER TABLE customers RENAME COLUMN id TO customer_id;
	`)
	if err != nil {
		t.Fatal(err)
	}

	logicalSchema := schema.NewLogicalSchema()
	if err := schema.ApplyBatch(logicalSchema, batch); err != nil {
		t.Fatal(err)
	}

	reference := logicalSchema.Tables["orders"].Constraints["orders_customer_fk"].Reference
	if reference.Columns[0] != "customer_id" {
		t.Fatalf("foreign key was not renamed: %#v", reference)
	}

	drop, err := parser.ParseDDL("DROP TABLE customers")
	if err != nil {
		t.Fatal(err)
	}
	if err := schema.ApplyStatement(logicalSchema, drop); err != nil {
		t.Fatal(err)
	}
	if len(logicalSchema.Tables["orders"].Constraints) != 1 {
		t.Fatalf("inbound foreign key was not removed: %#v", logicalSchema.Tables["orders"].Constraints)
	}
}

func TestApplyBatchIsAtomic(t *testing.T) {
	logicalSchema := schema.NewLogicalSchema()
	batch, err := parser.ParseBatch(`
		CREATE TABLE users (id bigint);
		ALTER TABLE missing ADD COLUMN id bigint;
	`)
	if err != nil {
		t.Fatal(err)
	}

	if err := schema.ApplyBatch(logicalSchema, batch); err == nil {
		t.Fatal("expected schema application to fail")
	}
	if len(logicalSchema.Tables) != 0 {
		t.Fatalf("schema was partially changed: %#v", logicalSchema.Tables)
	}
}
