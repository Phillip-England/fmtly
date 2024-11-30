package element

import (
	"fmt"
	"gtml/src/parser/attr"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementSlot struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	CompNames []string
	Attrs     []attr.Attr
}

func NewSlot(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementSlot, error) {
	elm := &ElementSlot{
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

func (elm *ElementSlot) GetSelection() *goquery.Selection { return elm.Selection }

//	func (elm *ElementSlot) GetParams() ([]param.Param, error) {
//		params := make([]param.Param, 0)
//		return params, nil
//	}
func (elm *ElementSlot) GetHtml() string        { return elm.Html }
func (elm *ElementSlot) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementSlot) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementSlot) GetType() string        { return elm.Type }
func (elm *ElementSlot) GetAttr() string        { return elm.Attr }
func (elm *ElementSlot) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementSlot) GetName() string        { return elm.Name }
func (elm *ElementSlot) GetCompNames() []string { return elm.CompNames }
func (elm *ElementSlot) GetAttrs() []attr.Attr  { return elm.Attrs }
func (elm *ElementSlot) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementSlot) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementSlot) initType() error {
	elm.Type = KeyElementSlot
	return nil
}

func (elm *ElementSlot) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementSlot) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementSlot)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementSlot, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementSlot) initAttrs() error {
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

func (elm *ElementSlot) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}
