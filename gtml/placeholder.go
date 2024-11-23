package gtml

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
type Placeholder interface {
	Print()
	GetHtml() string
	GetFuncCall() string
}

func NewPlaceholder(foundAsHtml string, name string) (Placeholder, error) {
	place, err := NewPlaceholderComponent(foundAsHtml, name)
	if err != nil {
		return nil, err
	}
	return place, nil
}

// ##==================================================================
type PlaceholderComponent struct {
	Name           string
	NodeName       string
	Html           string
	Attrs          []Attr
	FuncParamSlice []string
	FuncParamStr   string
	FuncCall       string
}

func NewPlaceholderComponent(foundAsHtml string, name string) (*PlaceholderComponent, error) {
	place := &PlaceholderComponent{
		Name: name,
	}
	err := fungi.Process(
		func() error { return place.initNodeName(foundAsHtml) },
		func() error { return place.initHtml(foundAsHtml) },
		func() error { return place.initAttrs() },
		func() error { return place.initFuncParamSlice() },
		func() error { return place.initComponentFuncCall() },
	)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func (place *PlaceholderComponent) initNodeName(foundAsHtml string) error {
	sel, err := gqpp.NewSelectionFromStr(foundAsHtml)
	if err != nil {
		return err
	}
	nodeName := goquery.NodeName(sel)
	place.NodeName = nodeName
	return nil
}

func (place *PlaceholderComponent) initHtml(foundAsHtml string) error {
	place.Html = foundAsHtml
	return nil
}

func (place *PlaceholderComponent) initAttrs() error {
	sel, err := gqpp.NewSelectionFromStr(place.Html)
	if err != nil {
		return err
	}
	for _, node := range sel.Nodes {
		for _, attr := range node.Attr {
			attrType, err := NewAttr(attr.Key, attr.Val)
			if err != nil {
				return err
			}
			place.Attrs = append(place.Attrs, attrType)
		}
	}
	return nil
}

func (place *PlaceholderComponent) initFuncParamSlice() error {
	funcParamSlice := make([]string, 0)
	for _, attr := range place.Attrs {
		if attr.GetKey() == KeyElementComponent {
			continue
		}
		if attr.GetType() == KeyAttrEmpty {
			continue
		}
		if attr.GetType() == KeyAttrStr {
			funcParamSlice = append(funcParamSlice, `"`+attr.GetValue()+`"`)
			continue
		}
		if attr.GetType() == KeyAttrAtParam {
			val := attr.GetValue()[1:]
			funcParamSlice = append(funcParamSlice, val)
			continue
		}
		if attr.GetType() == KeyAttrPlaceholder {
			sqVal := purse.Squeeze(attr.GetValue())
			sqVal = strings.Replace(sqVal, "{{", "", 1)
			sqVal = strings.Replace(sqVal, "}}", "", 1)
			funcParamSlice = append(funcParamSlice, "PARAM."+sqVal)
			continue
		}
		funcParamSlice = append(funcParamSlice, attr.GetValue())
	}
	place.FuncParamSlice = funcParamSlice
	place.FuncParamStr = strings.Join(funcParamSlice, ", ")
	return nil
}

func (place *PlaceholderComponent) initComponentFuncCall() error {
	call := fmt.Sprintf("%s(%s)", place.Name, place.FuncParamStr)
	place.FuncCall = call
	return nil
}
func (place *PlaceholderComponent) Print()              { fmt.Println(place.Html) }
func (place *PlaceholderComponent) GetHtml() string     { return place.Html }
func (place *PlaceholderComponent) GetFuncCall() string { return place.FuncCall }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
