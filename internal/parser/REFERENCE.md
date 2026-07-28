# ValkDB/postgresparser — reference notes

Source: https://github.com/ValkDB/postgresparser (Apache-2.0, pure Go, ANTLR4-based).
Import path: `github.com/valkdb/postgresparser`. Install: `go get github.com/valkdb/postgresparser@latest`.

These notes summarize the parts of the public API the `provider/postgres`
adapter depends on. It's a package of top-level functions, **not** a
stateful parser object — there is no constructor to instantiate.

## Entry points

| Function | Behavior |
|---|---|
| `ParseSQL(sql string) (*ParsedQuery, error)` | Parses only the **first** statement. Kept for backward compatibility. |
| `ParseSQLAll(sql string) (*ParseBatchResult, error)` | Parses **all** statements; one `StatementParseResult` per input statement, in source order. Individual statement failures don't abort the batch. |
| `ParseSQLStrict(sql string) (*ParsedQuery, error)` | Requires exactly one statement; returns `ErrMultipleStatements` otherwise. |
| `ParseSQLWithOptions` / `ParseSQLAllWithOptions` / `ParseSQLStrictWithOptions` | Same as above, plus a `ParseOptions` argument. |

`ParseOptions.IncludeCreateTableFieldComments` enables inline `--`
field-comment extraction on `CREATE TABLE` columns. `COMMENT ON` extraction
is always on regardless of options.

```go
result, err := postgresparser.ParseSQLAll(`
CREATE TABLE public.api_key (id integer NOT NULL);
CREATE TABLE public.sometable (id integer NOT NULL);
`)
// len(result.Statements) == 2
// result.Statements[0].Index, .RawSQL, .Query, .Warnings
// result.HasFailures == true if any Query is nil or any Warnings exist
```

## Batch result shape

```go
type StatementParseResult struct {
    Index    int
    RawSQL   string
    Query    *ParsedQuery   // nil means this statement failed conversion
    Warnings []ParseWarning
}

type ParseBatchResult struct {
    Statements  []StatementParseResult
    HasFailures bool
}
```

## `ParsedQuery` — the core IR (`ir.go`)

```go
type ParsedQuery struct {
    Command        QueryCommand   // SELECT/INSERT/UPDATE/DELETE/MERGE/DDL/UNKNOWN
    RawSQL         string
    Columns        []SelectColumn
    Tables         []TableRef
    ColumnUsage    []ColumnUsage
    SetOperations  []SetOperation
    Subqueries     []SubqueryRef
    CTEs           []CTE
    Where          []string        // raw predicate text, not a parsed tree
    Having         []string
    GroupBy        []string
    OrderBy        []OrderExpression
    Limit          *LimitClause
    JoinConditions []string
    Parameters     []Parameter
    Placeholders   []Placeholder
    InsertColumns  []string
    SetClauses     []string        // raw "col = expr" text
    Returning      []string
    Upsert         *UpsertClause
    Merge          *MergeClause
    DDLActions     []DDLAction
    Correlations   []JoinCorrelation
    DerivedColumns map[string]string
}
```

Important gotcha: `Where`, `SetClauses`, `GroupBy`, `JoinConditions` are
**raw text fragments**, not structured expression trees. Anything that needs
a left/operator/right split (like `ir.Condition`) has to parse that text
itself — see `provider/postgres/expression.go: parseCondition`.

### `TableRef` (FROM/JOIN entries)

```go
type TableRef struct {
    Schema        string
    Name          string
    Alias         string
    Type          TableType // base | cte | function | subquery
    Raw           string
    JoinType      string    // "INNER"/"LEFT"/"RIGHT"/"FULL"/"CROSS"/"NATURAL"/"" for base FROM
    JoinCondition string    // raw ON/USING text, "" for base/CROSS
    Nested        bool      // true if surfaced from inside a CTE/subquery, not this query's own FROM
}
```

Each joined table carries its *own* join type/condition — there's no
explicit pairwise join list. Reconstructing an `ir.Join{Left, Right}` means
pairing each entry with whichever table preceded it (see the limitation
noted in the provider's own README).

### `SelectColumn`

```go
type SelectColumn struct {
    Expression string // raw projection expression, e.g. "COUNT(o.id)"
    Alias      string
}
```

### DDL (`DDLAction`)

```go
type DDLActionType string
const (
    DDLCreateTable, DDLDropTable, DDLDropColumn, DDLAlterTable,
    DDLCreateIndex, DDLDropIndex, DDLTruncate, DDLComment DDLActionType = ...
)

type DDLAction struct {
    Type           DDLActionType
    ObjectName     string   // unqualified table/index/object name
    ObjectType     string   // TABLE, COLUMN, INDEX, ...
    Schema         string
    Columns        []string
    ColumnDetails  []DDLColumn      // CREATE TABLE only
    Constraints    *DDLConstraints  // CREATE TABLE, or ALTER TABLE ADD CONSTRAINT
    Flags          []string         // IF_EXISTS, CONCURRENTLY, CASCADE, UNIQUE, ...
    IndexType      string           // btree/gin/gist/hash — CREATE INDEX only
    IncludeColumns []string         // CREATE INDEX ... INCLUDE (...)
    Predicate      string           // partial-index WHERE expression, no leading "WHERE"
    Target         string           // generic fully-qualified target — mainly for COMMENT ON
    Comment        string           // COMMENT ON text
}
```

One SQL statement can produce **several** `DDLAction`s (e.g. a multi-clause
`ALTER TABLE`). There's currently no `DDLActionType` for renames
(`RENAME COLUMN`/`RENAME TO`) — if you need that, it isn't modeled upstream
yet.

```go
type DDLColumn struct {
    Name, Type, Default string
    Nullable            bool
    Comment             []string // preceding "--" line comments, if enabled
}

type DDLConstraints struct {
    PrimaryKey       *DDLPrimaryKey
    ForeignKeys      []DDLForeignKey
    UniqueKeys        []DDLUniqueConstraint
    CheckConstraints []DDLCheckConstraint
}
```

`FKAction` values (`ON DELETE`/`ON UPDATE`): `NO ACTION`, `RESTRICT`,
`CASCADE`, `SET NULL`, `SET DEFAULT`.

### `LimitClause`

```go
type LimitClause struct {
    Limit    string // raw text — may be a placeholder like "$1", not always an int literal
    Offset   string
    IsNested bool
}
```

## Errors (`errors.go`)

```go
var ErrNoStatements      = errors.New("no statements found")
var ErrMultipleStatements = errors.New("multiple statements found")
var ErrNilContext        = errors.New("nil context")

type MultipleStatementsError struct{ StatementCount int } // Unwrap() -> ErrMultipleStatements
type ParseErrors struct{ SQL string; Errors []SyntaxError } // aggregated syntax errors, .Error() renders them with carets
type SyntaxError struct{ Line, Column int; Message string; TokenIndex int }
```

Use `errors.Is` / `errors.As` against these — `provider/postgres/errors.go`
does exactly this to translate them into `ir.Err*`.

## Supported statements (per upstream docs)

| Category | Statements | Status |
|---|---|---|
| DML | SELECT, INSERT, UPDATE, DELETE, MERGE | Full IR extraction |
| DDL | CREATE TABLE, ALTER TABLE, DROP TABLE/INDEX, CREATE INDEX, TRUNCATE, COMMENT ON | Full IR extraction |
| Utility | SET, SHOW, RESET | Graceful — returns `UNKNOWN`, no error |
| Other | GRANT, REVOKE, CREATE VIEW/FUNCTION/TRIGGER, COPY, EXPLAIN, VACUUM, BEGIN/COMMIT/ROLLBACK | Not yet supported — may error or return `UNKNOWN` |

Full detail: `docs/supported-statements.md` and `docs/parsed-query.md` in the
upstream repo.

## Not currently used by the HyperFulcrum adapter, but available upstream

- `analysis` subpackage — `AnalyzeSQL`, column-usage roles, WHERE condition
  extraction (optionally schema-aware), schema-aware JOIN relationship
  detection, placeholder role classification. Useful if you later want
  lineage/index-advisor style features beyond plain IR conversion.
- `Placeholder` / `PlaceholderRole` — richer than `Parameter`; tells you
  *where* a `?`/`$N` sits syntactically (WHERE value, LIMIT, function arg,
  etc.). Not currently mapped into `ir`.