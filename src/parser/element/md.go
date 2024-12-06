package element

import (
	"fmt"
	"gtml/src/parser/attr"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementMd struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	CompNames []string
	Attrs     []attr.Attr
}

func NewMd(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementMd, error) {
	elm := &ElementMd{
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

func (elm *ElementMd) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementMd) GetHtml() string                  { return elm.Html }
func (elm *ElementMd) SetHtml(htmlStr string)           { elm.Html = htmlStr }
func (elm *ElementMd) Print()                           { fmt.Println(elm.Html) }
func (elm *ElementMd) GetType() string                  { return elm.Type }
func (elm *ElementMd) GetAttr() string                  { return elm.Attr }
func (elm *ElementMd) GetAttrParts() []string           { return elm.AttrParts }
func (elm *ElementMd) GetName() string                  { return elm.Name }
func (elm *ElementMd) GetCompNames() []string           { return elm.CompNames }
func (elm *ElementMd) GetAttrs() []attr.Attr            { return elm.Attrs }
func (elm *ElementMd) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementMd) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementMd) initType() error {
	elm.Type = KeyElementMd
	return nil
}

func (elm *ElementMd) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementMd) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementMd)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementMd, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementMd) initAttrs() error {
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

func (elm *ElementMd) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}
