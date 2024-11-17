package gtml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/gqpp"
)

// ##==================================================================
const (
	KeyElementComponent = "_component"
	KeyElementFor       = "_for"
)

// ##==================================================================
type Element interface {
	GetSelection() *goquery.Selection
	GetParam() (string, error)
	GetHtml() string
	Print()
	GetType() string
	GetAttr() string
	GetAttrParts() []string
	GetProps() []Prop
}

func GetFullElementList() []string {
	childElements := GetChildElementList()
	full := append(childElements, KeyElementComponent)
	return full
}

func GetChildElementList() []string {
	return []string{KeyElementFor}
}
func NewElement(sel *goquery.Selection) (Element, error) {
	match := gqpp.GetFirstMatchingAttr(sel, GetFullElementList()...)
	switch match {
	case KeyElementComponent:
		elm, err := NewElementComponent(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementFor:
		elm, err := NewElementFor(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("provided selection is not a valid element: %s", htmlStr)
}

func WalkElementChildren(elm Element, fn func(child Element) error) error {
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, inner *goquery.Selection) {
		child, err := NewElement(inner)
		if err != nil {
			// skip elements which are not a valid Element
		} else {
			err = fn(child)
			if err != nil {
				potErr = err
				return
			}
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func GetElementParams(elm Element) (string, error) {
	elementSpecificParams := make([]string, 0)
	err := WalkElementChildren(elm, func(child Element) error {
		param, err := child.GetParam()
		if err != nil {
			return err
		}
		if !slices.Contains(elementSpecificParams, param) && param != "" {
			elementSpecificParams = append(elementSpecificParams, param)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	strParams := make([]string, 0)
	for _, prop := range elm.GetProps() {
		strProp := prop.GetValue() + " " + "string"
		if !slices.Contains(strParams, strProp) && strings.Count(strProp, ".") == 0 {
			strParams = append(strParams, strProp)
		}
	}
	paramSlice := append(strParams, elementSpecificParams...)
	params := strings.Join(paramSlice, ", ")
	return params, nil
}

func WalkElementDirectChildren(elm Element, fn func(child Element) error) error {
	err := WalkElementChildren(elm, func(child Element) error {
		if !gqpp.HasParentWithAttrs(child.GetSelection(), elm.GetSelection(), GetChildElementList()...) {
			err := fn(child)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	htmlStr := elm.GetHtml()
	calls := ""
	err := WalkElementChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		if childHtml == "" {
			return nil
		}
		goVar, err := NewGoVar(child)
		if err != nil {
			return err
		}
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, goVar.GetVarName())
		calls += call + "\n"
		htmlStr = strings.Replace(htmlStr, childHtml, call, 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	for _, prop := range elm.GetProps() {
		call := PropAsWriteString(prop, builderName)
		htmlStr = strings.Replace(htmlStr, prop.GetRaw(), call, 1)
	}
	finalCalls := ""
	for {
		index := strings.Index(htmlStr, builderName)
		if index == -1 {
			break
		}
		part := htmlStr[:index]
		if part != "" && part != " " {
			finalCalls += fmt.Sprintf("%s.WriteString(`%s`)\n", builderName, part)
			htmlStr = strings.Replace(htmlStr, part, "", 1)
		}
		index = strings.Index(htmlStr, ")")
		if index == -1 {
			break
		}
		part = htmlStr[:index+1]
		if part != "" && part != " " {
			finalCalls += part + "\n"
			htmlStr = strings.Replace(htmlStr, part, "", 1)
		}
	}
	finalCalls += fmt.Sprintf("%s.WriteString(`%s`)\n", builderName, htmlStr)
	return finalCalls, nil
}

// ##==================================================================
type ElementComponent struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
}

func NewElementComponent(sel *goquery.Selection) (*ElementComponent, error) {
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	attr, err := gqpp.ForceElementAttr(sel, KeyElementComponent)
	if err != nil {
		return nil, err
	}
	parts, err := gqpp.ForceElementAttrParts(sel, KeyElementComponent, 1)
	if err != nil {
		return nil, err
	}
	elm := &ElementComponent{
		Selection: sel,
		Html:      htmlStr,
		Type:      KeyElementComponent,
		Attr:      attr,
		AttrParts: parts,
	}
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	props, err := NewProps(elm)
	if err != nil {
		return nil, err
	}
	elm.Props = props
	return elm, nil
}

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParam() (string, error)        { return "", nil }
func (elm *ElementComponent) GetHtml() string                  { return elm.Html }
func (elm *ElementComponent) Print()                           { fmt.Println(elm.Html) }
func (elm *ElementComponent) GetType() string                  { return elm.Type }
func (elm *ElementComponent) GetAttr() string                  { return elm.Attr }
func (elm *ElementComponent) GetAttrParts() []string           { return elm.AttrParts }
func (elm *ElementComponent) GetName() string                  { return elm.Name }
func (elm *ElementComponent) GetProps() []Prop                 { return elm.Props }

// ##==================================================================
type ElementFor struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
}

func NewElementFor(sel *goquery.Selection) (*ElementFor, error) {
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	attr, err := gqpp.ForceElementAttr(sel, KeyElementFor)
	if err != nil {
		return nil, err
	}
	parts, err := gqpp.ForceElementAttrParts(sel, KeyElementFor, 4)
	if err != nil {
		return nil, err
	}
	elm := &ElementFor{
		Selection: sel,
		Html:      htmlStr,
		Type:      KeyElementFor,
		Attr:      attr,
		AttrParts: parts,
	}
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	props, err := NewProps(elm)
	if err != nil {
		return nil, err
	}
	elm.Props = props
	return elm, nil
}

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParam() (string, error) {
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementFor, 4)
	if err != nil {
		return "", err
	}
	iterItems := parts[2]
	if strings.Contains(iterItems, ".") {
		return "", nil
	}
	iterType := parts[3]
	return iterItems + " " + iterType, nil
}
func (elm *ElementFor) GetHtml() string        { return elm.Html }
func (elm *ElementFor) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementFor) GetType() string        { return elm.Type }
func (elm *ElementFor) GetAttr() string        { return elm.Attr }
func (elm *ElementFor) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementFor) GetName() string        { return elm.Name }
func (elm *ElementFor) GetProps() []Prop       { return elm.Props }
