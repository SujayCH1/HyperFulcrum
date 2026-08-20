package ir

import (
	"strconv"
	"strings"
)

type DataType struct {
	Name            string
	Modifiers       []string
	ArrayDimensions int
}

func (d DataType) String() string {
	value := d.Name

	if len(d.Modifiers) > 0 {
		value += "(" + strings.Join(d.Modifiers, ",") + ")"
	}

	for range d.ArrayDimensions {
		value += "[]"
	}

	return value
}

func NewDataType(value string) DataType {
	dataType := DataType{}
	value = strings.TrimSpace(value)

	for strings.HasSuffix(value, "[]") {
		dataType.ArrayDimensions++
		value = strings.TrimSpace(strings.TrimSuffix(value, "[]"))
	}

	open := strings.IndexByte(value, '(')
	if open == -1 || !strings.HasSuffix(value, ")") {
		dataType.Name = value
		return dataType
	}

	dataType.Name = strings.TrimSpace(value[:open])
	for modifier := range strings.SplitSeq(value[open+1:len(value)-1], ",") {
		dataType.Modifiers = append(dataType.Modifiers, strings.TrimSpace(modifier))
	}

	return dataType
}

func DataTypeModifier(value int64) string {
	return strconv.FormatInt(value, 10)
}
