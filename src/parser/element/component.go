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

type ElementComponent struct {
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

func NewComponent(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementComponent, error) {
	elm := &ElementComponent{
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

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParams() ([]param.Param, error) {
	params := make([]param.Param, 0)
	return params, nil
}
func (elm *ElementComponent) GetHtml() string        { return elm.Html }
func (elm *ElementComponent) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementComponent) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementComponent) GetType() string        { return elm.Type }
func (elm *ElementComponent) GetAttr() string        { return elm.Attr }
func (elm *ElementComponent) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementComponent) GetName() string        { return elm.Name }
func (elm *ElementComponent) GetCompNames() []string { return elm.CompNames }
func (elm *ElementComponent) GetAttrs() []attr.Attr  { return elm.Attrs }
func (elm *ElementComponent) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementComponent) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementComponent) initType() error {
	elm.Type = KeyElementComponent
	return nil
}

func (elm *ElementComponent) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementComponent) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementComponent)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementComponent, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementComponent) initAttrs() error {
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

func (elm *ElementComponent) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementComponent) initRunes() error {
	r, err := GetElementRunes(elm)
	if err != nil {
		return err
	}
	elm.Runes = r
	return nil
}
