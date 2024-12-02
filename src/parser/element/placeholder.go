package element

import (
	"fmt"
	"gtml/src/parser/attr"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementPlaceholder struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	CompNames []string
	Attrs     []attr.Attr
}

func NewPlaceholder(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementPlaceholder, error) {
	elm := &ElementPlaceholder{
		CompNames: compNames,
	}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initAttrs() },
		func() error { return elm.initName() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementPlaceholder) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementPlaceholder) GetHtml() string                  { return elm.Html }
func (elm *ElementPlaceholder) SetHtml(htmlStr string)           { elm.Html = htmlStr }
func (elm *ElementPlaceholder) Print()                           { fmt.Println(elm.Html) }
func (elm *ElementPlaceholder) GetType() string                  { return elm.Type }
func (elm *ElementPlaceholder) GetAttr() string                  { return elm.Attr }
func (elm *ElementPlaceholder) GetAttrParts() []string           { return elm.AttrParts }
func (elm *ElementPlaceholder) GetName() string                  { return elm.Name }
func (elm *ElementPlaceholder) GetCompNames() []string           { return elm.CompNames }
func (elm *ElementPlaceholder) GetAttrs() []attr.Attr            { return elm.Attrs }
func (elm *ElementPlaceholder) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementPlaceholder) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementPlaceholder) initType() error {
	elm.Type = KeyElementPlaceholder
	return nil
}

func (elm *ElementPlaceholder) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementPlaceholder) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementPlaceholder)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementPlaceholder, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementPlaceholder) initAttrs() error {
	for _, a := range elm.GetSelection().Get(0).Attr {
		if purse.MustEqualOneOf(a.Key, GetChildElementList()...) {
			continue
		}
		attr, err := attr.NewAttr(a.Key, a.Val)
		if err != nil {
			return err
		}
		elm.Attrs = append(elm.Attrs, attr)
	}
	return nil
}

func (elm *ElementPlaceholder) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}
