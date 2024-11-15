package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ##==================================================================
type Element interface {
	SetData(data string)
	GetData() string
	AddChild(child Element)
	GetChildren() []Element
	HasChildren() bool
	GetParent() Element
	DeleteSelf() error
	DeleteChildren()
	GetAttrData() string
	Type() string
	IsRoot() bool
	Print()
	GetProps() []string
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
	elm.DeleteChildren()
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

func DeleteElement(elm Element) error {
	parent := elm.GetParent()
	parentChildren := parent.GetChildren()
	if len(parentChildren) == 0 {
		return fmt.Errorf("attempted to delete element whos parent has no childre (which should not occur): %s", elm.GetData())
	}
	parent.DeleteChildren()
	for _, child := range parentChildren {
		if child == elm {
			continue
		}
		parent.AddChild(child)
	}
	return nil
}

func CloneElement(elm Element) (Element, error) {
	var clone Element
	if elm.Type() == "component" {
		if elm.IsRoot() {
			clone, err := NewComponentElement(elm.GetData(), nil)
			if err != nil {
				return nil, err
			}
			return clone, nil
		} else {
			parent, err := CloneElement(elm.GetParent())
			if err != nil {
				return nil, err
			}
			clone, err := NewComponentElement(elm.GetData(), parent)
			if err != nil {
				return nil, err
			}
			return clone, nil
		}
	}
	if elm.Type() == "for" {
		if elm.IsRoot() {
			clone, err := NewForElement(elm.GetData(), nil)
			if err != nil {
				return nil, err
			}
			return clone, nil
		} else {
			parent, err := CloneElement(elm.GetParent())
			if err != nil {
				return nil, err
			}
			clone, err := NewForElement(elm.GetData(), parent)
			if err != nil {
				return nil, err
			}
			return clone, nil
		}
	}
	return clone, nil
}

func WalkChildElements(root Element, fn func(next Element) error) error {
	for _, child := range root.GetChildren() {
		err := fn(child)
		if err != nil {
			return err
		}
		if child.HasChildren() {
			err := WalkChildElements(child, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WalkToRoot(elm Element, fn func(next Element) error) error {
	err := fn(elm)
	if err != nil {
		return err
	}
	if elm.IsRoot() {
		return nil
	}
	return WalkToRoot(elm.GetParent(), fn)
}

func SetElementData(elm Element, data string) error {
	originalData := elm.GetData()
	newData := purse.Flatten(data)
	err := WalkToRoot(elm, func(next Element) error {
		overwrite := strings.Replace(next.GetData(), originalData, newData, 1)
		next.SetData(overwrite)
		return nil
	})
	if err != nil {
		return err
	}
	if !elm.HasChildren() {
		return nil
	}
	err = SetElementChildren(elm)
	if err != nil {
		return err
	}
	return nil
}

func WalkUpElementBranches(elm Element, fn func(next Element) error) error {
	out := make([]Element, 0)
	finalOut := make([]Element, 0)
	err := WalkChildElements(elm, func(next Element) error {
		out = append(out, next)
		if !next.HasChildren() {
			out = purse.ReverseSlice[Element](out)
			finalOut = append(finalOut, out...)
			out = make([]Element, 0)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, next := range finalOut {
		err := fn(next)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetElementProps(elm Element) []string {
	return purse.ScanBetweenSubStrs(elm.GetData(), "{{", "}}")
}

// ##==================================================================
type ComponentElement struct {
	Data          string
	Children      []Element
	ElementType   string
	Parent        Element
	IsRootElement bool
	AttrData      string
	Props         []string
}

func NewComponentElement(str string, parent Element) (*ComponentElement, error) {
	elm := &ComponentElement{
		Data:          str,
		ElementType:   "component",
		Parent:        parent,
		IsRootElement: parent == nil,
	}
	err := SetElementChildren(elm)
	if err != nil {
		return nil, err
	}
	sel, err := gqpp.NewSelectionFromStr(elm.GetData())
	if err != nil {
		return nil, err
	}
	attr, _ := sel.Attr("_component")
	elm.AttrData = attr
	elm.Props = GetElementProps(elm)
	return elm, nil
}

func (elm *ComponentElement) SetData(data string)    { elm.Data = data }
func (elm *ComponentElement) GetData() string        { return elm.Data }
func (elm *ComponentElement) AddChild(child Element) { elm.Children = append(elm.Children, child) }
func (elm *ComponentElement) GetChildren() []Element { return elm.Children }
func (elm *ComponentElement) HasChildren() bool      { return len(elm.Children) > 0 }
func (elm *ComponentElement) GetParent() Element     { return elm.Parent }
func (elm *ComponentElement) DeleteSelf() error      { return DeleteElement(elm) }
func (elm *ComponentElement) DeleteChildren()        { elm.Children = make([]Element, 0) }
func (elm *ComponentElement) GetAttrData() string    { return elm.AttrData }
func (elm *ComponentElement) Type() string           { return elm.ElementType }
func (elm *ComponentElement) IsRoot() bool           { return elm.IsRootElement }
func (elm *ComponentElement) Print()                 { fmt.Println(elm.Data) }
func (elm *ComponentElement) GetProps() []string     { return elm.Props }

// ##==================================================================
type ForElement struct {
	Data          string
	Children      []Element
	ElementType   string
	Parent        Element
	IsRootElement bool
	AttrParts     []string
	AttrData      string
	Props         []string
}

func NewForElement(str string, parent Element) (*ForElement, error) {
	elm := &ForElement{
		Data:          str,
		ElementType:   "for",
		Parent:        parent,
		IsRootElement: parent == nil,
	}
	err := SetElementChildren(elm)
	if err != nil {
		return nil, err
	}
	sel, err := gqpp.NewSelectionFromStr(elm.GetData())
	if err != nil {
		return nil, err
	}
	forAttr, _ := sel.Attr("_for")
	elm.AttrData = forAttr
	parts := strings.Split(forAttr, " ")
	if len(parts) != 4 {
		return nil, fmt.Errorf("_for element requires attributes with the following schema: ITEM of ITEMS []TYPE: %s", str)
	}
	elm.AttrParts = parts
	elm.Props = GetElementProps(elm)
	return elm, nil
}

func (elm *ForElement) SetData(data string)    { elm.Data = data }
func (elm *ForElement) GetData() string        { return elm.Data }
func (elm *ForElement) GetChildren() []Element { return elm.Children }
func (elm *ForElement) HasChildren() bool      { return len(elm.Children) > 0 }
func (elm *ForElement) AddChild(child Element) { elm.Children = append(elm.Children, child) }
func (elm *ForElement) GetParent() Element     { return elm.Parent }
func (elm *ForElement) DeleteSelf() error      { return DeleteElement(elm) }
func (elm *ForElement) DeleteChildren()        { elm.Children = make([]Element, 0) }
func (elm *ForElement) GetAttrData() string    { return elm.AttrData }
func (elm *ForElement) Type() string           { return elm.ElementType }
func (elm *ForElement) IsRoot() bool           { return elm.IsRootElement }
func (elm *ForElement) Print()                 { fmt.Println(elm.Data) }
func (elm *ForElement) GetProps() []string     { return elm.Props }

// ##==================================================================

// ##==================================================================
