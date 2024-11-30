package element

import (
	"fmt"
	"gtml/src/parser/attr"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementIf struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	CompNames []string
	Attrs     []attr.Attr
}

func NewIf(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementIf, error) {
	elm := &ElementIf{
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

func (elm *ElementIf) GetSelection() *goquery.Selection { return elm.Selection }

//	func (elm *ElementIf) GetParams() ([]param.Param, error) {
//		params := make([]param.Param, 0)
//		param, err := param.NewParam(elm.Attr, "bool")
//		if err != nil {
//			return nil, err
//		}
//		params = append(params, param)
//		return params, nil
//	}
func (elm *ElementIf) GetHtml() string        { return elm.Html }
func (elm *ElementIf) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementIf) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementIf) GetType() string        { return elm.Type }
func (elm *ElementIf) GetAttr() string        { return elm.Attr }
func (elm *ElementIf) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementIf) GetName() string        { return elm.Name }
func (elm *ElementIf) GetCompNames() []string { return elm.CompNames }
func (elm *ElementIf) GetAttrs() []attr.Attr  { return elm.Attrs }
func (elm *ElementIf) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementIf) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementIf) initType() error {
	elm.Type = KeyElementIf
	return nil
}

func (elm *ElementIf) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementIf) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementIf)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementIf, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementIf) initAttrs() error {
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

func (elm *ElementIf) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}
