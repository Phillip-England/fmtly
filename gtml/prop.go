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
	if strings.Count(val, ".") == 1 && len(strings.Split(val, ".")) == 2 {
		prop, err := NewPropForType(str, val)
		if err != nil {
			return nil, err
		}
		return prop, nil
	}
	return nil, nil
}

// ##==================================================================
type PropForType struct {
	Raw   string
	Value string
}

func NewPropForType(raw string, val string) (*PropForType, error) {
	return nil, nil
}

// ##==================================================================
type PropForStr struct {
	Raw   string
	Value string
}

func NewPropForStr(raw string, val string) (*PropForType, error) {
	return nil, nil
}

// ##==================================================================
type PropStr struct {
	Raw   string
	Value string
}

func NewPropStr(raw string, val string) (*PropForType, error) {
	return nil, nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
