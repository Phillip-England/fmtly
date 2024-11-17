package gtml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyElementComponent       = "COMPONENT"
	IndicatorElementComponent = "_component"
	KeyElementFor             = "FOR"
	IndicatorElementFor       = "_for"
)

// ##==================================================================
type Element interface {
	GetSelection() *goquery.Selection
	GetParam() (string, error)
	GetHtml() string
	Print()
	GetType() string
}

func GetFullElementList() []string {
	childList := GetChildElementList()
	full := append(childList, IndicatorElementComponent)
	return full
}

func GetChildElementList() []string {
	return []string{IndicatorElementFor}
}

func NewElement(sel *goquery.Selection) (Element, error) {
	match := gqpp.GetFirstMatchingAttr(sel, GetFullElementList()...)
	switch match {
	case IndicatorElementComponent:
		comp, err := NewElementComponent(sel)
		if err != nil {
			return nil, err
		}
		return comp, nil
	case IndicatorElementFor:
		comp, err := NewElementFor(sel)
		if err != nil {
			return nil, err
		}
		return comp, nil
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
		if err == nil {
			err := fn(child)
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
	err = WalkElementProps(elm, func(prop Prop) error {
		strProp := prop.GetValue() + " " + "string"
		if !slices.Contains(strParams, strProp) && strings.Count(strProp, ".") == 0 {
			strParams = append(strParams, strProp)
		}
		return nil
	})
	if err != nil {
		return "", err
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

func GetElementProps(elm Element) ([]string, error) {
	props := purse.ScanBetweenSubStrs(elm.GetHtml(), "{{", "}}")
	return props, nil
}

func GetPropValue(prop string) string {
	return purse.Squeeze(purse.RemoveAllSubStr(prop, "{{", "}}"))
}

func WalkElementProps(elm Element, fn func(prop Prop) error) error {
	props, err := NewProps(elm)
	if err != nil {
		return err
	}
	for _, prop := range props {
		err := fn(prop)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	htmlStr := elm.GetHtml()
	calls := ""
	err := WalkElementChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
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
	err = WalkElementProps(elm, func(prop Prop) error {
		call := PropAsWriteString(prop, builderName)
		htmlStr = strings.Replace(htmlStr, prop.GetRaw(), call, 1)
		return nil
	})
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
}

func NewElementComponent(sel *goquery.Selection) (*ElementComponent, error) {
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	elm := &ElementComponent{
		Selection: sel,
		Html:      htmlStr,
		Type:      KeyElementComponent,
	}
	return elm, nil
}

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParam() (string, error)        { return "", nil }
func (elm *ElementComponent) GetHtml() string                  { return elm.Html }
func (elm *ElementComponent) Print()                           { fmt.Println(elm.Html) }
func (elm *ElementComponent) GetType() string                  { return elm.Type }

// ##==================================================================
type ElementFor struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
}

func NewElementFor(sel *goquery.Selection) (*ElementFor, error) {
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	elm := &ElementFor{
		Selection: sel,
		Html:      htmlStr,
		Type:      KeyElementFor,
	}
	return elm, nil
}

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParam() (string, error) {
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), IndicatorElementFor, 4)
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
func (elm *ElementFor) GetHtml() string { return elm.Html }
func (elm *ElementFor) Print()          { fmt.Println(elm.Html) }
func (elm *ElementFor) GetType() string { return elm.Type }
