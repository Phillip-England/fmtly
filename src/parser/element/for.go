package element

import (
	"fmt"
	"gtml/src/parser/attr"
	"gtml/src/parser/gtmlrune"
	"gtml/src/parser/param"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type ElementFor struct {
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

func NewFor(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementFor, error) {
	elm := &ElementFor{
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

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParams() ([]param.Param, error) {
	params := make([]param.Param, 0)
	parts := elm.GetAttrParts()
	iterItems := parts[2]
	if strings.Contains(iterItems, ".") {
		return nil, nil
	}
	iterType := parts[3]
	param, err := param.NewParam(iterItems, iterType)
	if err != nil {
		return nil, err
	}
	params = append(params, param)
	return params, nil
}
func (elm *ElementFor) GetHtml() string        { return elm.Html }
func (elm *ElementFor) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementFor) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementFor) GetType() string        { return elm.Type }
func (elm *ElementFor) GetAttr() string        { return elm.Attr }
func (elm *ElementFor) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementFor) GetName() string        { return elm.Name }
func (elm *ElementFor) GetCompNames() []string { return elm.CompNames }
func (elm *ElementFor) GetAttrs() []attr.Attr  { return elm.Attrs }
func (elm *ElementFor) GetId() string {
	salt, _ := elm.GetSelection().Attr("_id")
	return salt
}

func (elm *ElementFor) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementFor) initType() error {
	elm.Type = KeyElementFor
	return nil
}

func (elm *ElementFor) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementFor) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementFor)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementFor, 4)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementFor) initAttrs() error {
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

func (elm *ElementFor) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementFor) initRunes() error {
	r, err := GetElementRunes(elm)
	if err != nil {
		return err
	}
	elm.Runes = r
	return nil
}
