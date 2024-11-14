package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type ComponentElement struct {
	Value    *goquery.Selection
	Children []GtmlElement
	Id       string
	Html     string
}

func NewComponentElementFromSelection(sel *goquery.Selection) (ComponentElement, error) {
	elm := &ComponentElement{
		Value: sel,
		Id:    "component",
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
