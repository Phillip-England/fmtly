package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

// ##==================================================================
type Element interface {
	SetData(data string)
	GetData() string
	AddChild(child Element)
	GetChildren() []Element
	GetParent() Element
	DeleteSelf()
	DeleteChildren()
	Print()
}

func NewElement(str string, parent Element) (Element, error) {
	str = purse.Flatten(str)
	sel, err := gqpp.NewSelectionFromStr(str)
	if err != nil {
		return nil, err
	}
	elementType := gqpp.GetFirstMatchingAttr(sel, "_component", "_for")
	if elementType == "_component" {
		return NewComponentElement(str, parent)
	}
	if elementType == "_for" {
		return NewForElement(str, parent)
	}
	return nil, fmt.Errorf("provided string is not a valid gtml element: %s", str)
}

func SetElementChildren(elm Element) error {
	sel, err := gqpp.NewSelectionFromStr(elm.GetData())
	if err != nil {
		return err
	}
	children := make([]Element, 0)
	var potErr error
	sel.Find("*[_for]").Each(func(i int, inner *goquery.Selection) {
		if !gqpp.HasParentWithAttrs(inner, sel, "_for") {
			htmlStr, err := gqpp.NewHtmlFromSelection(inner)
			if err != nil {
				potErr = err
				return
			}
			child, err := NewElement(htmlStr, elm)
			if err != nil {
				potErr = err
				return
			}
			children = append(children, child)
		}
	})
	if potErr != nil {
		return potErr
	}
	for _, child := range children {
		elm.AddChild(child)
	}
	return nil
}

func DeleteElement(elm Element) {
	parent := elm.GetParent()
	parentChildren := parent.GetChildren()
	parent.DeleteChildren()
	for _, child := range parentChildren {
		if child == elm {
			continue
		}
		parent.AddChild(child)
	}
}

// ##==================================================================
type ComponentElement struct {
	Data        string
	Children    []Element
	ElementType string
	Parent      Element
}

func NewComponentElement(str string, parent Element) (*ComponentElement, error) {
	elm := &ComponentElement{
		Data:        str,
		ElementType: "component",
		Parent:      parent,
	}
	err := SetElementChildren(elm)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ComponentElement) SetData(data string)    { elm.Data = data }
func (elm *ComponentElement) GetData() string        { return elm.Data }
func (elm *ComponentElement) AddChild(child Element) { elm.Children = append(elm.Children, child) }
func (elm *ComponentElement) GetChildren() []Element { return elm.Children }
func (elm *ComponentElement) GetParent() Element     { return elm.Parent }
func (elm *ComponentElement) DeleteSelf()            { DeleteElement(elm) }
func (elm *ComponentElement) DeleteChildren()        { elm.Children = make([]Element, 0) }
func (elm *ComponentElement) Print()                 { fmt.Println(elm.Data) }
func (elm *ComponentElement) Type() string           { return elm.ElementType }

// ##==================================================================
type ForElement struct {
	Data        string
	Children    []Element
	ElementType string
	Parent      Element
}

func NewForElement(str string, parent Element) (*ForElement, error) {
	elm := &ForElement{
		Data:        str,
		ElementType: "for",
		Parent:      parent,
	}
	err := SetElementChildren(elm)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ForElement) SetData(data string)    { elm.Data = data }
func (elm *ForElement) GetData() string        { return elm.Data }
func (elm *ForElement) GetChildren() []Element { return elm.Children }
func (elm *ForElement) AddChild(child Element) { elm.Children = append(elm.Children, child) }
func (elm *ForElement) GetParent() Element     { return elm.Parent }
func (elm *ForElement) DeleteSelf()            { DeleteElement(elm) }
func (elm *ForElement) DeleteChildren()        { elm.Children = make([]Element, 0) }
func (elm *ForElement) Print()                 { fmt.Println(elm.Data) }
func (elm *ForElement) Type() string           { return elm.ElementType }
