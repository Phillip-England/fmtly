package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type ComponentElement struct {
	Value         *goquery.Selection
	Children      []GtmlElement
	Id            string
	Html          string
	Parent        GtmlElement
	IsRootElement bool
}

func NewComponentElementFromSelection(sel *goquery.Selection, parent GtmlElement) (ComponentElement, error) {
	elm := &ComponentElement{
		Value:  sel,
		Id:     "component:" + purse.RandStr(16),
		Parent: nil, // will never have a parent
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return *elm, err
	}
	elm.Html = purse.Flatten(htmlStr)
	children, err := GetGtmlElementChildren(elm)
	if err != nil {
		return *elm, err
	}
	elm.Children = children
	if elm.Parent == nil {
		elm.IsRootElement = true
	}
	return *elm, nil
}

func (elm ComponentElement) GetChildren() []GtmlElement       { return elm.Children }
func (elm ComponentElement) GetHtml() string                  { return elm.Html }
func (elm ComponentElement) GetSelection() *goquery.Selection { return elm.Value }
func (elm ComponentElement) GetId() string                    { return elm.Id }
func (elm ComponentElement) HasChildren() bool                { return len(elm.Children) > 0 }
func (elm ComponentElement) Print()                           { fmt.Println(elm.Html) }
func (elm ComponentElement) GetWriteStringCall() (string, bool) {
	return "", false
}
func (elm ComponentElement) GetParent() GtmlElement { return elm.Parent }
func (elm ComponentElement) IsRoot() bool           { return elm.IsRootElement }
func (elm ComponentElement) Delete()                { DeleteGtmlElementNode(elm) }
func (elm ComponentElement) SetChildren(children []GtmlElement) {
	elm.Children = children
	return
}
