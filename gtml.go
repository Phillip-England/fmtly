package gtml

import (
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type GtmlRoot struct {
	Value ComponentElement
	Nodes []GtmlNode
}

func NewGtmlRootFromStr(str string) (GtmlRoot, error) {
	str = purse.Flatten(str)
	root := &GtmlRoot{}
	comp, err := NewComponentElementFromStr(str)
	if err != nil {
		return *root, err
	}
	root.Value = comp
	return *root, nil
}

type GtmlNode struct {
	Value  GtmlElement
	Parent *GtmlNode
	Nodes  []GtmlNode
}

type GtmlElement interface{}

type ComponentElement struct {
	Selection *goquery.Selection
	Html      string
}

func NewComponentElementFromStr(str string) (ComponentElement, error) {
	str = purse.Flatten(str)
	elm := &ComponentElement{
		Html: str,
	}
	sel, err := gqpp.NewSelectionFromHtmlStr(str)
	if err != nil {
		return *elm, err
	}
	elm.Selection = sel
	return *elm, nil
}

type ForElement struct {
	Selection *goquery.Selection
	Html      string
}

func NewForElementFromStr(str string) (ForElement, error) {
	str = purse.Flatten(str)
	elm := &ForElement{
		Html: str,
	}
	sel, err := gqpp.NewSelectionFromHtmlStr(str)
	if err != nil {
		return *elm, err
	}
	elm.Selection = sel
	return *elm, nil
}

// you can filter all of the first generation children in a selection out if they do not have sepecific attrs
