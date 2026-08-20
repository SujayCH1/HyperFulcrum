package ir

import "strings"

type Table struct {
	Schema string
	Name   string
	Alias  string
}

func (t Table) Key() string {
	if t.Schema == "" || t.Schema == "public" {
		return t.Name
	}

	return t.Schema + "." + t.Name
}

func NewTable(value string) Table {
	parts := strings.SplitN(value, ".", 2)
	if len(parts) == 1 {
		return Table{Name: value}
	}

	return Table{Schema: parts[0], Name: parts[1]}
}
