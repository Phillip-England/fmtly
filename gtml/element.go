package gtml

import (
	"fmt"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
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

func WalkElementProps(elm Element, fn func(prop Prop) error) error {
	allProps := make([]Prop, 0)
	for _, prop := range elm.GetProps() {
		allProps = append(allProps, prop)
	}
	err := WalkElementChildren(elm, func(child Element) error {
		for _, prop := range child.GetProps() {
			allProps = append(allProps, prop)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, prop := range allProps {
		err := fn(prop)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	clay := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		loopVar, err := NewGoVar(child)
		if err != nil {
			return err
		}
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, loopVar.GetVarName())
		clay = strings.Replace(clay, childHtml, call, 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	err = WalkElementProps(elm, func(prop Prop) error {
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, prop.GetValue())
		clay = strings.Replace(clay, prop.GetRaw(), call, 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	if strings.Index(clay, builderName) == -1 {
		singleCall := fmt.Sprintf("%s.WriteString(%s)", builderName, clay)
		return singleCall, nil
	}
	series := ""
	for {
		builderIndex := strings.Index(clay, builderName)
		if builderIndex == -1 {
			break
		}
		htmlPart := clay[:builderIndex]
		htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, htmlPart)
		series += htmlCall + "\n"
		clay = strings.Replace(clay, htmlPart, "", 1)
		endBuilderIndex := strings.Index(clay, ")")
		builderPart := clay[:endBuilderIndex+1]
		builderCall := fmt.Sprintf("%s.WriteString(%s)", builderName, builderPart)
		series += builderCall + "\n"
		clay = strings.Replace(clay, builderPart, "", 1)
	}
	if len(clay) > 0 {
		htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		series += htmlCall + "\n"
	}
	return series, nil
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
	elm := &ElementComponent{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
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

func (elm *ElementComponent) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementComponent) initProps() error {
	elmHtml := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		elmHtml = strings.Replace(elmHtml, childHtml, "", 1)
		return nil
	})
	if err != nil {
		return err
	}
	strProps := purse.ScanBetweenSubStrs(elmHtml, "{{", "}}")
	for _, strProp := range strProps {
		prop, err := NewProp(strProp)
		if err != nil {
			return err
		}
		elm.Props = append(elm.Props, prop)
	}
	return nil
}

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
	elm := &ElementFor{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
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

func (elm *ElementFor) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementFor) initProps() error {
	elmHtml := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		elmHtml = strings.Replace(elmHtml, childHtml, "", 1)
		return nil
	})
	if err != nil {
		return err
	}
	strProps := purse.ScanBetweenSubStrs(elmHtml, "{{", "}}")
	for _, strProp := range strProps {
		prop, err := NewProp(strProp)
		if err != nil {
			return err
		}
		elm.Props = append(elm.Props, prop)
	}
	return nil
}
