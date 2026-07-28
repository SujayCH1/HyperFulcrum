# internal/parser

Dialect-agnostic SQL parsing for HyperFulcrum. Callers depend only on
`parser.Parser` and the `ir` package; concrete SQL dialects live under
`provider/<dialect>` and are swapped in by `factory.go`.

```
internal/parser/
├── interfaces.go        // Parser interface (Parse, ParseBatch)
├── factory.go            // New() -> wires up the active provider
├── errors.go              // Re-exports of ir.Err* for callers
│
├── ir/                    // Dialect-neutral intermediate representation
│   ├── common.go          // StatementKind, Command enums
│   ├── identifier.go      // Identifier, Table
│   ├── datatype.go        // DataType
│   ├── statement.go       // Statement interface
│   ├── schema.go          // Column, Constraint, Index, DDLStatement, AlterOperation
│   ├── query.go           // DMLStatement
│   ├── expression.go       // Expression, Condition, Join, JoinType
│   ├── batch.go           // Batch
│   ├── metadata.go        // Metadata (optional, non-essential context)
│   └── errors.go          // Canonical sentinel errors
│
└── provider/
    └── postgres/          // ValkDB postgresparser adapter
        ├── parser.go       // PostgreSQLParser: Parse / ParseBatch
        ├── converter.go    // Dispatches ParsedQuery -> DML/DDL converter
        ├── ddl.go          // DDLAction(s) -> ir.DDLStatement
        ├── dml.go          // ParsedQuery -> ir.DMLStatement
        ├── expression.go   // Raw predicate text -> ir.Condition
        ├── datatype.go     // Raw column type text -> ir.DataType
        ├── mapper.go       // postgresparser enums -> ir enums
        ├── helpers.go      // small constructors/parsers
        └── errors.go       // postgresparser errors -> ir.Err*
```

## Setup

The Postgres provider depends on the external ValkDB library. If `go build`
or your editor reports the provider package as unresolved/"missing metadata",
check this first:

```bash
go get github.com/valkdb/postgresparser@latest
go mod tidy
```

Also confirm your `go.mod`'s `module` line matches the import paths used
throughout this package exactly (case-sensitive) — e.g. `module hyperfulcrum`
if that's what every file imports as `hyperfulcrum/internal/parser/...`.

## Usage

```go
p := parser.New()

stmt, err := p.Parse("SELECT id, name FROM users WHERE active = true")
if err != nil {
    // handle
}

switch s := stmt.(type) {
case *ir.DMLStatement:
    // s.Tables, s.Conditions, s.Columns, ...
case *ir.DDLStatement:
    // s.Table, s.Columns, s.Constraints, ...
}

batch, err := p.ParseBatch("CREATE TABLE a (id int); CREATE TABLE b (id int);")
```

`Parse` uses the first statement only; `ParseBatch` parses every statement in
the input and skips (rather than fails on) any individual statement that
couldn't be converted — check `len(batch.Statements)` against the number of
statements you sent in if you need to detect partial failures.

## Why errors live in `ir`, not `parser`

`factory.go` (package `parser`) imports `provider/postgres` to build the
concrete implementation. If `provider/postgres` imported `parser` back — even
just for the four sentinel errors — you get `parser → postgres → parser`,
an import cycle Go refuses to compile.

The fix: `ir` has zero internal dependencies, and both `parser` and every
provider already import it for `ir.Statement`, `ir.Batch`, etc. So the
sentinel errors (`ErrNilStatement`, `ErrUnsupportedSQL`, `ErrEmptySQL`,
`ErrInvalidAST`) are defined once in `ir/errors.go`. `parser/errors.go`
re-exports them as `var ErrEmptySQL = ir.ErrEmptySQL` purely so existing call
sites can keep writing `parser.ErrEmptySQL`.

This also means `provider/postgres.New()` returns the **concrete** type
`*postgres.PostgreSQLParser`, not `parser.Parser`. Go interfaces are
satisfied structurally, so the provider package never needs to import the
interface's package — the conversion to `parser.Parser` happens implicitly
in `factory.go`, which is the one place that legitimately imports both sides.

## Adding another dialect (e.g. MySQL)

1. New directory `provider/mysql/`, package `mysql`.
2. It may import `ir` and its own upstream SQL library. It must **not**
   import `hyperfulcrum/internal/parser` — return sentinel errors from `ir`
   directly, and return the concrete `*MySQLParser` type from `mysql.New()`.
3. Wire it into `factory.go`, or extend `New()` to pick a provider based on a
   dialect argument if you need more than one active at once.

## Known limitations in the current Postgres provider

These are conscious simplifications, not oversights — flagged so they don't
surprise you later:

- **Join reconstruction** (`dml.go: convertJoins`) pairs each joined table
  with the table immediately preceding it in source order. This matches
  simple chained joins but won't correctly reconstruct a join tree that
  branches (e.g. two tables both joining back to the same earlier alias).
- **`CREATE INDEX` table attribution** (`ddl.go: convertDDL`) uses
  `DDLAction.Target`, whose doc comment describes it as generic and
  primarily used for `COMMENT ON`. Verify against real `CREATE INDEX` output
  before relying on it.
- **No rename support**: ValkDB's `DDLActionType` has no `RENAME_COLUMN` /
  `RENAME_TABLE` equivalent yet, so `ir.RenameColumn` / `ir.RenameTable` are
  defined but currently unreachable from this provider.
- **`MERGE` is rejected outright** (`converter.go`) rather than force-fit
  onto `ir.DMLStatement`, since ValkDB's `MergeClause` (target/source/
  WHEN-MATCHED actions) has no equivalent shape in the current IR.
- **UPDATE's new value** is stored in `ir.Column.DefaultValue` when
  converting `SET col = expr` clauses, since `ir.Column` has no dedicated
  "assigned value" field. Fine functionally, but worth a clearer field if
  this trips someone up during review.