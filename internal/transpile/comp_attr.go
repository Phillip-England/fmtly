package transpile

import (
	"fmt"
	"fmtly/internal/parsley"
	"strings"
)

type CompAttr struct {
	Name  string
	Value string
}

func NewCompAttr(rawAttr string) (*CompAttr, error) {
	rawAttr = strings.Replace(rawAttr, "=", " ", 1)
	rawAttr = strings.Replace(rawAttr, "\"", "", 1)
	rawAttr = strings.Replace(rawAttr, "'", "", 1)
	rawAttr = rawAttr[:len(rawAttr)-1]
	parts := parsley.SplitSpaces(rawAttr)
	if len(parts) != 2 {
		return nil, fmt.Errorf("attr is malformed: %s", rawAttr)
	}
	compAttr := &CompAttr{
		Name:  parts[0],
		Value: parts[1],
	}
	return compAttr, nil
}
