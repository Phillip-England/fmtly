package gtmlfunc

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlvar"
	"gtml/src/parser/param"
)

type Func interface {
	GetData() string
	SetData(str string)
	GetVars() []gtmlvar.Var
	GetParams() []param.Param
	Print()
}

func NewFunc(elm element.Element, siblings []element.Element) (Func, error) {
	filtered := make([]element.Element, 0)
	for _, sibling := range siblings {
		if sibling.GetName() == elm.GetName() {
			continue
		}
		filtered = append(filtered, sibling)
	}
	if elm.GetType() == element.KeyElementComponent || elm.GetType() == element.KeyElementPlaceholder {
		fn, err := NewGoComponentFunc(elm, filtered)
		if err != nil {
			return nil, err
		}
		return fn, nil
	}
	return nil, fmt.Errorf("provided element does not corrospond to a valid GoFunc: %s", elm.GetHtml())
}
