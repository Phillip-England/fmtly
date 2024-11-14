package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type IfElement struct {
	Value    *goquery.Selection
	Children []GtmlElement
	Id       string
	Html     string
}

func (elm IfElement) GetChildren() []GtmlElement       { return elm.Children }
func (elm IfElement) GetHtml() string                  { return elm.Html }
func (elm IfElement) GetSelection() *goquery.Selection { return elm.Value }
func (elm IfElement) GetId() string                    { return elm.Id }
func (elm IfElement) HasChildren() bool                { return len(elm.Children) > 0 }
func (elm IfElement) Print()                           { fmt.Println(elm.Html) }

func NewIfElementFromSelection(sel *goquery.Selection) (IfElement, error) {
	elm := &IfElement{
		Value: sel,
		Id:    "if",
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
