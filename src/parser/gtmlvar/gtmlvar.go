package gtmlvar

import (
	"fmt"
	"gtml/src/parser/element"

	"github.com/phillip-england/purse"
)

type Var interface {
	GetData() string
	GetVarName() string
	GetType() string
	Print()
	GetBuilderName() string
	GetElement() element.Element
}

func NewVar(elm element.Element) (Var, error) {
	switch elm.GetType() {
	case element.KeyElementFor:
		v, err := NewGoFor(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementIf:
		v, err := NewGoIf(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementElse:
		v, err := NewGoElse(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementPlaceholder:
		v, err := NewGoPlaceholder(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementSlot:
		v, err := NewGoSlot(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("element does not corrospond to a valid Var: %s", elm.GetHtml())
}

func NewVarsFromElement(elm element.Element) ([]Var, error) {
	vars := make([]Var, 0)
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		innerVar, err := NewVar(child)
		if err != nil {
			return nil
		}
		vars = append(vars, innerVar)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func GetVarsAsWriteStringCalls(elm element.Element, builderName string) ([]string, error) {
	calls := make([]string, 0)
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		newVar, err := NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, GetFullVarList()...) {
			if varType == KeyVarGoPlaceholder {
				call := fmt.Sprintf("%s.WriteString(%s())", builderName, newVar.GetVarName())
				calls = append(calls, call)
				return nil
			}
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, newVar.GetVarName())
			calls = append(calls, call)
		}
		return nil
	})
	if err != nil {
		return calls, err
	}
	return calls, nil
}
