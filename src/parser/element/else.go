package element

import (
	"fmt"
	"gtml/src/parser/attr"
	"gtml/src/parser/gtmlrune"
	"gtml/src/parser/param"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementElse struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	CompNames []string
	Attrs     []attr.Attr
	Runes     []gtmlrune.GtmlRune
}

func NewElse(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementElse, error) {
	elm := &ElementElse{
		CompNames: compNames,
	}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initAttrs() },
		func() error { return elm.initName() },
		func() error { return elm.initRunes() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementElse) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementElse) GetParams() ([]param.Param, error) {
	params := make([]param.Param, 0)
	param, err := param.NewParam(elm.Attr, "bool")
	if err != nil {
		return nil, err
	}
	params = append(params, param)
	return params, nil
}
func (elm *ElementElse) GetHtml() string        { return elm.Html }
func (elm *ElementElse) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementElse) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementElse) GetType() string        { return elm.Type }
func (elm *ElementElse) GetAttr() string        { return elm.Attr }
func (elm *ElementElse) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementElse) GetName() string        { return elm.Name }
func (elm *ElementElse) GetCompNames() []string { return elm.CompNames }
func (elm *ElementElse) GetAttrs() []attr.Attr  { return elm.Attrs }
func (elm *ElementElse) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementElse) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementElse) initType() error {
	elm.Type = KeyElementElse
	return nil
}

func (elm *ElementElse) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementElse) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementElse)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementElse, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementElse) initAttrs() error {
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

func (elm *ElementElse) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementElse) initRunes() error {
	r, err := GetElementRunes(elm)
	if err != nil {
		return err
	}
	elm.Runes = r
	return nil
}
