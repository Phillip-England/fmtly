package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
)

// ##==================================================================
type Placeholder interface {
	Print()
	GetFoundAs() string
	GetPointingTo() Element
}

func NewPlaceholder(foundAsHtml string, pointingTo Element) (Placeholder, error) {
	if pointingTo.GetType() != KeyElementComponent {
		return nil, fmt.Errorf("a placeholder must point to a valid _component element: %s", pointingTo.GetHtml())
	}
	place, err := NewPlaceholderComponent(foundAsHtml, pointingTo)
	if err != nil {
		return nil, err
	}
	return place, nil
}

// ##==================================================================
type PlaceholderComponent struct {
	Name              string
	Html              string
	PointingTo        Element
	Params            []Param
	Attrs             []Attr
	FuncParamSlice    []string
	ComponentFuncCall string
}

func NewPlaceholderComponent(foundAsHtml string, pointingTo Element) (*PlaceholderComponent, error) {
	place := &PlaceholderComponent{
		Html:       foundAsHtml,
		PointingTo: pointingTo,
	}
	err := fungi.Process(
		func() error { return place.initName() },
		func() error { return place.initParamNames() },
		func() error { return place.initAttrs() },
		func() error { return place.initFuncParamSlice() },
		func() error { return place.initComponentFuncCall() },
	)
	if err != nil {
		return nil, err
	}
	place.Print()
	return place, nil
}

func (place *PlaceholderComponent) initName() error {
	nameAttr := place.PointingTo.GetAttr()
	place.Name = nameAttr
	return nil
}

func (place *PlaceholderComponent) initParamNames() error {
	params, err := GetElementParams(place.PointingTo)
	if err != nil {
		return err
	}
	for _, param := range params {
		place.Params = append(place.Params, param)
	}
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

	place.FuncParamSlice = funcParamSlice
	return nil
}

func (place *PlaceholderComponent) initComponentFuncCall() error {
	paramStr := strings.Join(place.FuncParamSlice, ", ")
	call := fmt.Sprintf("%s(%s)", place.Name, paramStr)
	place.ComponentFuncCall = call
	return nil
}

func (place *PlaceholderComponent) Print() {
	fmt.Println("Name: " + place.Name)
	fmt.Println("Html: " + place.Html)
	fmt.Print("PointingTo: ")
	place.PointingTo.Print()
	fmt.Println("ComponentFuncCall: " + place.ComponentFuncCall)
}
func (place *PlaceholderComponent) GetFoundAs() string     { return place.Html }
func (place *PlaceholderComponent) GetPointingTo() Element { return place.PointingTo }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
