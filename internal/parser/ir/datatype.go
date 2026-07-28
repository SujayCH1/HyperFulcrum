package ir

type DataType struct {
	Name string

	// VARCHAR(255)
	Length *int

	// NUMERIC(10,2)
	Precision *int
	Scale     *int

	// TEXT[], BIGINT[]
	Array bool
}
