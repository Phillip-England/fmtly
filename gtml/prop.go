package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/purse"
)

// ##==================================================================
type Prop interface{}

func NewProp(str string) (Prop, error) {
	if !strings.HasPrefix(str, "{{") && !strings.HasSuffix(str, "}}") {
		return nil, fmt.Errorf("provided string is not a valid Prop: %s", str)
	}
	val := purse.RemoveAllSubStr(str, "{{", "}}")
	val = purse.Squeeze(val)
	if val == "" {
		return nil, fmt.Errorf("empty prop tag provided: %s", str)
	}
	if strings.Count(val, ".") == 1 {
		// forTypeProp,

	}
	return nil, nil
}

// ##==================================================================
type ForTypeProp struct {
	Raw   string
	Value string
}

// ##==================================================================

// ##==================================================================
