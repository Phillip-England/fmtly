package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type IfElement struct {
	Value         *goquery.Selection
	Children      []GtmlElement
	Id            string
	Html          string
	BuilderName   string
	IfAttr        string
	Parent        GtmlElement
	IsRootElement bool
}

func NewIfElementFromSelection(sel *goquery.Selection, parent GtmlElement) (IfElement, error) {
	elm := &IfElement{
		Value:  sel,
		Id:     "if:" + purse.RandStr(16),
		Parent: parent,
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return *elm, err
	}
	elm.Html = purse.Flatten(htmlStr)
	ifAttr, exists := sel.Attr("_if")
	if !exists {
		return *elm, fmt.Errorf("_if element requires an _if attribute: %s", elm.Html)
	}
	elm.IfAttr = ifAttr
	children, err := GetGtmlElementChildren(elm)
	if err != nil {
		return *elm, err
	}
	elm.Children = children
	elm.BuilderName = fmt.Sprintf("%sBuilderIf", elm.IfAttr)
	if elm.Parent == nil {
		elm.IsRootElement = true
	}
	return *elm, nil
}

func (elm IfElement) GetChildren() []GtmlElement       { return elm.Children }
func (elm IfElement) GetHtml() string                  { return elm.Html }
func (elm IfElement) GetSelection() *goquery.Selection { return elm.Value }
func (elm IfElement) GetId() string                    { return elm.Id }
func (elm IfElement) HasChildren() bool                { return len(elm.Children) > 0 }
func (elm IfElement) Print()                           { fmt.Println(elm.Html) }
func (elm IfElement) GetWriteStringCall() (string, bool) {
	call := fmt.Sprintf("%s.WriteString(%s)", elm.BuilderName, elm.IfAttr)
	return call, true
}
func (elm IfElement) GetParent() GtmlElement { return elm.Parent }
func (elm IfElement) IsRoot() bool           { return elm.IsRootElement }
func (elm IfElement) Delete()                { DeleteGtmlElementNode(elm) }

func (elm IfElement) SetChildren(children []GtmlElement) { elm.Children = children }
