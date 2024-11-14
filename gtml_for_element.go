package gtml

import (
	"fmt"
	"gtml/internal/gqpp"
	"gtml/internal/purse"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ForElement struct {
	Value         *goquery.Selection
	Children      []GtmlElement
	Id            string
	Html          string
	ForAttr       string
	ItemName      string
	ItemType      string
	ItemsName     string
	ItemsType     string
	Props         []GtmlProp
	BuilderName   string
	LoopName      string
	Parent        GtmlElement
	IsRootElement bool
}

func NewForElementFromSelection(sel *goquery.Selection, parent GtmlElement) (ForElement, error) {
	elm := &ForElement{
		Value:  sel,
		Id:     "for:" + purse.RandStr(16),
		Parent: parent,
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
	forAttr, exists := sel.Attr("_for")
	if !exists {
		return *elm, fmt.Errorf("_for required on a for element: %s", elm.Html)
	}
	elm.ForAttr = forAttr
	parts := strings.Split(elm.ForAttr, " ")
	if len(parts) != 4 {
		return *elm, fmt.Errorf("_for attr value must follow this schema: ITEM of ITEMS []TYPE")
	}
	firstPart := parts[0]
	elm.ItemName = firstPart
	thirdPart := parts[2]
	elm.ItemsName = thirdPart
	fourthPart := parts[3]
	elm.ItemsType = fourthPart
	fourthPart = purse.RemoveAllSubStr(fourthPart, "[]")
	elm.ItemType = fourthPart
	props := NewPropsFromStr(elm.Html)
	elm.Props = props
	elm.BuilderName = fmt.Sprintf("%sBuilder", elm.ItemName)
	elm.LoopName = fmt.Sprintf("%sLoop", elm.ItemName)
	if elm.Parent == nil {
		elm.IsRootElement = true
	}
	return *elm, nil
}

func (elm ForElement) GetChildren() []GtmlElement       { return elm.Children }
func (elm ForElement) GetHtml() string                  { return elm.Html }
func (elm ForElement) GetSelection() *goquery.Selection { return elm.Value }
func (elm ForElement) GetId() string                    { return elm.Id }
func (elm ForElement) HasChildren() bool                { return len(elm.Children) > 0 }
func (elm ForElement) Print()                           { fmt.Println(elm.Html) }
func (elm ForElement) GetWriteStringCall() (string, bool) {
	call := fmt.Sprintf("%s.WriteString(%s)", elm.BuilderName, elm.LoopName)
	return call, true
}
func (elm ForElement) GetParent() GtmlElement { return elm.Parent }
func (elm ForElement) IsRoot() bool           { return elm.IsRootElement }
func (elm ForElement) Delete()                { DeleteGtmlElementNode(elm) }

func (elm ForElement) SetChildren(children []GtmlElement) { elm.Children = children }