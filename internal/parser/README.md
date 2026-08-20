# Parser

The parser converts PostgreSQL 16 SQL into project-owned representations.

```go
ddl, err := parser.ParseDDL("CREATE TABLE users (id bigint PRIMARY KEY)")
route, err := parser.ParseRoute("SELECT * FROM users WHERE id = $1")
batch, err := parser.ParseBatch("CREATE TABLE users (id bigint); ALTER TABLE users ADD COLUMN name text;")
```

DDL statements produce `ir.DDLStatement` values for the schema module. DML
statements produce compact `ir.RouteStatement` values for query routing.

Parsing is strict. `Parse` rejects multiple statements and `ParseBatch` fails
the complete batch when any statement is invalid or unsupported.

Supported schema statements include `CREATE TABLE`, `ALTER TABLE`, table and
column renames, `DROP TABLE`, `CREATE INDEX`, `DROP INDEX`, and `TRUNCATE`.

Routing extraction supports `SELECT`, `INSERT`, `UPDATE`, `DELETE`, and
`MERGE`. `RoutingComplete` is false when a query contains a construct that the
router must handle conservatively, such as an `OR`, CTE, subquery, set
operation, or computed routing value.
