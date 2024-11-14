package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type ElseElement struct {
	Value         *goquery.Selection
	Children      []GtmlElement
	Id            string
	Html          string
	ElseAttr      string
	BuilderName   string
	Parent        GtmlElement
	IsRootElement bool
}

func NewElseElementFromSelection(sel *goquery.Selection, parent GtmlElement) (ElseElement, error) {
	elm := &ElseElement{
		Value:  sel,
		Id:     "else:" + purse.RandStr(16),
		Parent: parent,
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return *elm, err
	}
	elseAttr, exists := sel.Attr("_else")
	if !exists {
		return *elm, fmt.Errorf("_else element requires an _else attribute: %s", elm.Html)
	}
	elm.ElseAttr = elseAttr
	elm.Html = purse.Flatten(htmlStr)
	children, err := GetGtmlElementChildren(elm)
	if err != nil {
		return *elm, err
	}
	elm.Children = children
	elm.BuilderName = fmt.Sprintf("%sBuilderElse", elm.ElseAttr)
	if elm.Parent == nil {
		elm.IsRootElement = true
	}
	return *elm, nil
}

func (elm ElseElement) GetChildren() []GtmlElement       { return elm.Children }
func (elm ElseElement) GetHtml() string                  { return elm.Html }
func (elm ElseElement) GetSelection() *goquery.Selection { return elm.Value }
func (elm ElseElement) GetId() string                    { return elm.Id }
func (elm ElseElement) HasChildren() bool                { return len(elm.Children) > 0 }
func (elm ElseElement) Print()                           { fmt.Println(elm.Html) }
func (elm ElseElement) GetWriteStringCall() (string, bool) {
	call := fmt.Sprintf("%s.WriteString(%s)", elm.BuilderName, elm.ElseAttr)
	return call, true
}
func (elm ElseElement) GetParent() GtmlElement { return elm.Parent }
func (elm ElseElement) IsRoot() bool           { return elm.IsRootElement }
func (elm ElseElement) Delete()                { DeleteGtmlElementNode(elm) }

func (elm ElseElement) SetChildren(children []GtmlElement) { elm.Children = children }
