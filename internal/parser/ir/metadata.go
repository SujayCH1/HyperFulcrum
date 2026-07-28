package ir

// Metadata carries optional, non-essential information about how a
// Statement or Batch was produced. It's not required for the IR's
// correctness — callers that don't need it can ignore it entirely.
type Metadata struct {
	// RawSQL preserves the original source text for the parsed statement.
	RawSQL string

	// SourceDialect names the provider that produced this IR (e.g. "postgres").
	SourceDialect string

	// Warnings holds non-fatal notices from the underlying parser, such as
	// syntax-recovery notices emitted during batch parsing.
	Warnings []string
}
