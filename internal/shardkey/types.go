package shardkey

type ColumnRef struct {
	Table string
	Column string
}
type CandidateSet map[string][]ColumnRef

type FanoutStats struct {
	IncomingFKs       int
	ReferencingTables int
}

type RankedCandidate struct {
	Column  ColumnRef
	Score   int
	Reasons []string
}

type ShardKeyDecision struct {
	Table   string
	Column  ColumnRef
	Score   int
	Reasons []string
}

type InferenceResult struct {
	ProjectID string
	Decisions []ShardKeyDecision
}

