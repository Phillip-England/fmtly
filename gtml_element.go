package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type GtmlElement interface {
	GetChildren() []GtmlElement
	GetHtml() string
	GetSelection() *goquery.Selection
	GetId() string
	HasChildren() bool
	Print()
	GetWriteStringCall() (string, bool)
	GetParent() GtmlElement
	IsRoot() bool
	Delete()
	SetChildren(children []GtmlElement)
}

func NewGtmlElementFromStr(str string) (GtmlElement, error) {
	sel, err := gqpp.NewSelectionFromHtmlStr(str)
	if err != nil {
		return nil, err
	}
	_, exists := sel.Attr("_component")
	if exists {
		elm, err := NewComponentElementFromSelection(sel, nil)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	_, exists = sel.Attr("_for")
	if exists {
		elm, err := NewForElementFromSelection(sel, nil)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	_, exists = sel.Attr("_if")
	if exists {
		elm, err := NewIfElementFromSelection(sel, nil)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	_, exists = sel.Attr("_else")
	if exists {
		elm, err := NewElseElementFromSelection(sel, nil)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	return nil, fmt.Errorf("Provided html is not a valid gtml element: %s", str)
}

func GetGtmlElementChildren(parent GtmlElement) ([]GtmlElement, error) {
	children := make([]GtmlElement, 0)
	sel := parent.GetSelection()
	var potErr error
	sel.Find("*[_for]").Each(func(i int, forSel *goquery.Selection) {
		if !gqpp.HasParentWithAttrs(forSel, sel, "_for", "_else", "_if") {
			forElm, err := NewForElementFromSelection(forSel, parent)
			if err != nil {
				potErr = err
				return
			}
			children = append(children, forElm)
		}
	})
	if potErr != nil {
		return nil, potErr
	}
	sel.Find("*[_if]").Each(func(i int, forSel *goquery.Selection) {
		if !gqpp.HasParentWithAttrs(forSel, sel, "_for", "_else", "_if") {
			ifElm, err := NewIfElementFromSelection(forSel, parent)
			if err != nil {
				potErr = err
				return
			}
			children = append(children, ifElm)
		}
	})
	if potErr != nil {
		return nil, potErr
	}
	sel.Find("*[_else]").Each(func(i int, forSel *goquery.Selection) {
		if !gqpp.HasParentWithAttrs(forSel, sel, "_for", "_else", "_if") {
			elseElm, err := NewElseElementFromSelection(forSel, parent)
			if err != nil {
				potErr = err
				return
			}
			children = append(children, elseElm)
		}
	})
	if potErr != nil {
		return nil, potErr
	}
	return children, nil
}

func WalkGtmlChildren(elm GtmlElement, fn func(child GtmlElement) error) error {
	for _, child := range elm.GetChildren() {
		err := fn(child)
		if err != nil {
			return err
		}
		if child.HasChildren() {
			err := WalkGtmlChildren(child, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func WalkUpGtmlBranches(elm GtmlElement, fn func(child GtmlElement) error) error {
	collect := make([]GtmlElement, 0)
	final := make([]GtmlElement, 0)
	WalkGtmlChildren(elm, func(child GtmlElement) error {
		collect = append(collect, child)
		if !child.HasChildren() {
			collect = purse.ReverseSlice[GtmlElement](collect)
			final = append(final, collect...)
			collect = make([]GtmlElement, 0)
		}
		return nil
	})
	for _, inner := range final {
		err := fn(inner)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteGtmlElementNode(elm GtmlElement) {
	if elm.IsRoot() {
		return
	}
	parent := elm.GetParent()
	filtered := make([]GtmlElement, 0)
	for _, child := range parent.GetChildren() {
		if child.GetId() == elm.GetId() {
			continue
		}
		filtered = append(filtered, child)
	}
	parent.SetChildren(filtered)
	parent.Print()
	// for _, child := range parent.GetChildren() {
	// 	child.Print()
	// }
}
