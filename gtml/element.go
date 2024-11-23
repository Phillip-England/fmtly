package gtml

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyElementComponent   = "_component"
	KeyElementFor         = "_for"
	KeyElementIf          = "_if"
	KeyElementElse        = "_else"
	KeyElementPlaceholder = "_placeholder"
	KeyElementSlot        = "_slot"
)

// ##==================================================================
type Element interface {
	GetSelection() *goquery.Selection
	GetParams() ([]Param, error)
	SetHtml(htmlStr string)
	GetHtml() string
	Print()
	GetType() string
	GetAttr() string
	GetAttrParts() []string
	GetProps() []Prop
	GetCompNames() []string
	GetAttrs() []Attr
}

func GetFullElementList() []string {
	childElements := GetChildElementList()
	full := append(childElements, KeyElementComponent)
	return full
}

func GetChildElementList() []string {
	// KeyElementSlot must go last
	// other elements take priority over KeyElementSlot
	return []string{KeyElementFor, KeyElementIf, KeyElementElse, KeyElementPlaceholder, KeyElementSlot}
}

func NewElement(htmlStr string, compNames []string) (Element, error) {
	sel, err := gqpp.NewSelectionFromStr(htmlStr)
	if err != nil {
		return nil, err
	}
	match := gqpp.GetFirstMatchingAttr(sel, GetFullElementList()...)
	switch match {
	case KeyElementComponent:
		elm, err := NewElementComponent(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementFor:
		elm, err := NewElementFor(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementIf:
		elm, err := NewElementIf(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementElse:
		elm, err := NewElementElse(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementPlaceholder:
		elm, err := NewElementPlaceholder(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementSlot:
		elm, err := NewElementSlot(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}

	return nil, fmt.Errorf("provided selection is not a valid element: %s", htmlStr)
}

func WalkElementChildren(elm Element, fn func(child Element) error) error {
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, inner *goquery.Selection) {
		htmlStr, err := gqpp.NewHtmlFromSelection(inner)
		child, err := NewElement(htmlStr, elm.GetCompNames())
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

func WalkElementChildrenIncludingRoot(elm Element, fn func(child Element) error) error {
	err := fn(elm)
	if err != nil {
		return err
	}
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, inner *goquery.Selection) {
		htmlStr, err := gqpp.NewHtmlFromSelection(inner)
		if err != nil {
			potErr = err
			return
		}
		child, err := NewElement(htmlStr, elm.GetCompNames())
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

func GetElementParams(elm Element) ([]Param, error) {
	params := make([]Param, 0)
	elementSpecificParams := make([]Param, 0)
	err := WalkElementChildren(elm, func(child Element) error {
		params, err := child.GetParams()
		if err != nil {
			return err
		}
		for _, param := range params {
			if !slices.Contains(elementSpecificParams, param) && param != nil {
				elementSpecificParams = append(elementSpecificParams, param)
			}
		}
		return nil
	})
	if err != nil {
		return params, err
	}
	strParams := make([]Param, 0)
	for _, prop := range elm.GetProps() {
		param, err := NewParam(prop.GetValue(), "string")
		if err != nil {
			return params, err
		}
		if !slices.Contains(strParams, param) && purse.MustEqualOneOf(prop.GetType(), KeyPropStr, KeyPropSlot) {
			strParams = append(strParams, param)
		}
	}
	params = append(strParams, elementSpecificParams...)
	return params, nil
}

func WalkElementDirectChildren(elm Element, fn func(child Element) error) error {
	var potErr error
	elm.GetSelection().Children().Each(func(i int, childSel *goquery.Selection) {
		if gqpp.HasAttr(childSel, GetChildElementList()...) {
			childHtml, err := gqpp.NewHtmlFromSelection(childSel)
			if err != nil {
				potErr = err
				return
			}
			childElm, err := NewElement(childHtml, elm.GetCompNames())
			if err != nil {
				potErr = err
				return
			}
			err = fn(childElm)
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

func GetElementHtmlWithoutChildren(elm Element) (string, error) {
	elmHtml := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		elmHtml = strings.Replace(elmHtml, childHtml, "", 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	return elmHtml, nil
}

func GetElementProps(elm Element) ([]Prop, error) {
	props := make([]Prop, 0)
	elmHtml := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		elmHtml = strings.Replace(elmHtml, childHtml, "", 1)
		return nil
	})
	if err != nil {
		return props, err
	}
	strProps := purse.ScanBetweenSubStrs(elmHtml, "{{", "}}")
	for _, strProp := range strProps {
		prop, err := NewProp(strProp, elm.GetCompNames())
		if err != nil {
			return props, err
		}
		props = append(props, prop)
	}
	return props, nil
}

func GetElementVars(elm Element) ([]Var, error) {
	vars := make([]Var, 0)
	err := WalkElementDirectChildren(elm, func(child Element) error {
		innerVar, err := NewVar(child)
		if err != nil {
			return nil
		}
		vars = append(vars, innerVar)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	clay := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		newVar, err := NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, KeyVarGoElse, KeyVarGoFor, KeyVarGoIf, KeyVarGoPlaceholder, KeyVarGoSlot) {
			if varType == KeyVarGoPlaceholder {
				call := fmt.Sprintf("%s.WriteString(%s())", builderName, newVar.GetVarName())
				clay = strings.Replace(clay, childHtml, call, 1)
			}
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, newVar.GetVarName())
			clay = strings.Replace(clay, childHtml, call, 1)
		}
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
		singleCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		return singleCall, nil
	}
	series := ""
	for {
		builderIndex := strings.Index(clay, builderName)
		if builderIndex == -1 {
			break
		}
		htmlPart := clay[:builderIndex]
		if htmlPart != "" {
			htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, htmlPart)
			series += htmlCall + "\n"
			clay = strings.Replace(clay, htmlPart, "", 1)
		}
		endBuilderIndex := strings.Index(clay, ")")
		loopCount := 0
		for {
			loopCount++
			nextChar := string(clay[endBuilderIndex+loopCount])
			if nextChar == ")" {
				endBuilderIndex = endBuilderIndex + loopCount
				continue
			}
			break
		}
		builderPart := clay[:endBuilderIndex+1]
		series += builderPart + "\n"
		clay = strings.Replace(clay, builderPart, "", 1)
	}
	if len(clay) > 0 {
		htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		series += htmlCall + "\n"
	}
	return series, nil
}

func WalkAllElementNodes(elm Element, fn func(sel *goquery.Selection) error) error {
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, s *goquery.Selection) {
		err := fn(s)
		if err != nil {
			potErr = err
			return
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func WalkAllElementNodesIncludingRoot(elm Element, fn func(sel *goquery.Selection) error) error {
	err := fn(elm.GetSelection())
	if err != nil {
		return nil
	}
	err = WalkAllElementNodes(elm, func(sel *goquery.Selection) error {
		err := fn(sel)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WalkAllElementNodesWithoutChildren(elm Element, fn func(sel *goquery.Selection) error) error {
	htmlNoChildren, err := GetElementHtmlWithoutChildren(elm)
	if err != nil {
		return err
	}
	sel, err := gqpp.NewSelectionFromStr(htmlNoChildren)
	if err != nil {
		return err
	}
	var potErr error
	sel.Find("*").Each(func(i int, s *goquery.Selection) {
		err := fn(s)
		if err != nil {
			potErr = err
			return
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func ReadComponentElementsFromFile(path string, compNames []string) ([]Element, error) {
	elms := make([]Element, 0)
	f, err := os.ReadFile(path)
	if err != nil {
		return elms, err
	}
	fStr := string(f)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
	if err != nil {
		return elms, err
	}
	var potErr error
	doc.Find("*").Each(func(i int, sel *goquery.Selection) {
		_, exists := sel.Attr(KeyElementComponent)
		if exists {
			htmlStr, err := gqpp.NewHtmlFromSelection(sel)
			if err != nil {
				potErr = err
				return
			}
			elm, err := NewElement(htmlStr, compNames)
			if err != nil {
				potErr = err
				return
			}
			elms = append(elms, elm)
		}
	})
	if potErr != nil {
		return elms, potErr
	}
	return elms, nil
}

func ReadComponentElementNamesFromFile(path string) ([]string, error) {
	names := make([]string, 0)
	f, err := os.ReadFile(path)
	if err != nil {
		return names, err
	}
	fStr := string(f)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
	if err != nil {
		return names, err
	}
	doc.Find("*").Each(func(i int, sel *goquery.Selection) {
		compAttr, exists := sel.Attr(KeyElementComponent)
		if exists {
			names = append(names, compAttr)
		}
	})
	return names, nil
}

func MarkElementPlaceholders(elm Element) (Element, error) {
	clay := elm.GetHtml()
	err := WalkAllElementNodes(elm, func(sel *goquery.Selection) error {
		nodeName := goquery.NodeName(sel)
		ogSelHtml, err := gqpp.NewHtmlFromSelection(sel)
		if err != nil {
			return err
		}
		for _, name := range elm.GetCompNames() {
			if strings.ToLower(name) == nodeName {
				sel.SetAttr("_placeholder", name)
				selHtml, err := gqpp.NewHtmlFromSelection(sel)
				if err != nil {
					return err
				}
				var potErr error
				sel.Children().Each(func(i int, childSel *goquery.Selection) {
					_, hasSlot := childSel.Attr("_slot")
					if !hasSlot {
						potErr = fmt.Errorf("placeholder element has children which are not wrapped in an element with a _slot='slotName' attribute: %s", ogSelHtml)
						return
					}
				})
				if potErr != nil {
					return potErr
				}
				clay = strings.Replace(clay, ogSelHtml, selHtml, 1)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	newElm, err := NewElement(clay, elm.GetCompNames())
	if err != nil {
		return nil, err
	}
	return newElm, nil
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
	CompNames []string
	Attrs     []Attr
}

func NewElementComponent(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementComponent, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParams() ([]Param, error) {
	params := make([]Param, 0)
	return params, nil
}
func (elm *ElementComponent) GetHtml() string        { return elm.Html }
func (elm *ElementComponent) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementComponent) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementComponent) GetType() string        { return elm.Type }
func (elm *ElementComponent) GetAttr() string        { return elm.Attr }
func (elm *ElementComponent) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementComponent) GetName() string        { return elm.Name }
func (elm *ElementComponent) GetProps() []Prop       { return elm.Props }
func (elm *ElementComponent) GetCompNames() []string { return elm.CompNames }
func (elm *ElementComponent) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementComponent) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
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
	CompNames []string
	Attrs     []Attr
}

func NewElementFor(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementFor, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParams() ([]Param, error) {
	params := make([]Param, 0)
	parts := elm.GetAttrParts()
	iterItems := parts[2]
	if strings.Contains(iterItems, ".") {
		return nil, nil
	}
	iterType := parts[3]
	param, err := NewParam(iterItems, iterType)
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
func (elm *ElementFor) GetProps() []Prop       { return elm.Props }
func (elm *ElementFor) GetCompNames() []string { return elm.CompNames }
func (elm *ElementFor) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementFor) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementIf struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
	CompNames []string
	Attrs     []Attr
}

func NewElementIf(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementIf, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementIf) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementIf) GetParams() ([]Param, error) {
	params := make([]Param, 0)
	param, err := NewParam(elm.Attr, "bool")
	if err != nil {
		return nil, err
	}
	params = append(params, param)
	return params, nil
}
func (elm *ElementIf) GetHtml() string        { return elm.Html }
func (elm *ElementIf) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementIf) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementIf) GetType() string        { return elm.Type }
func (elm *ElementIf) GetAttr() string        { return elm.Attr }
func (elm *ElementIf) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementIf) GetName() string        { return elm.Name }
func (elm *ElementIf) GetProps() []Prop       { return elm.Props }
func (elm *ElementIf) GetCompNames() []string { return elm.CompNames }
func (elm *ElementIf) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementIf) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementElse struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
	CompNames []string
	Attrs     []Attr
}

func NewElementElse(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementElse, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementElse) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementElse) GetParams() ([]Param, error) {
	params := make([]Param, 0)
	param, err := NewParam(elm.Attr, "bool")
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
func (elm *ElementElse) GetProps() []Prop       { return elm.Props }
func (elm *ElementElse) GetCompNames() []string { return elm.CompNames }
func (elm *ElementElse) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementElse) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementPlaceholder struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
	CompNames []string
	Attrs     []Attr
}

func NewElementPlaceholder(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementPlaceholder, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementPlaceholder) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementPlaceholder) GetParams() ([]Param, error) {
	// working on building out attribute params and making
	// sure they are being pulled into the () of the component func
	params := make([]Param, 0)
	for _, attr := range elm.Attrs {
		if attr.GetType() == KeyAttrInitParam {
			param, err := NewParam(attr.GetValue(), "string")
			if err != nil {
				return params, err
			}
			params = append(params, param)
		}
	}
	return params, nil
}
func (elm *ElementPlaceholder) GetHtml() string        { return elm.Html }
func (elm *ElementPlaceholder) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementPlaceholder) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementPlaceholder) GetType() string        { return elm.Type }
func (elm *ElementPlaceholder) GetAttr() string        { return elm.Attr }
func (elm *ElementPlaceholder) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementPlaceholder) GetName() string        { return elm.Name }
func (elm *ElementPlaceholder) GetProps() []Prop       { return elm.Props }
func (elm *ElementPlaceholder) GetCompNames() []string { return elm.CompNames }
func (elm *ElementPlaceholder) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementPlaceholder) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementSlot struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
	CompNames []string
	Attrs     []Attr
}

func NewElementSlot(htmlStr string, sel *goquery.Selection, compNames []string) (*ElementSlot, error) {
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
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementSlot) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementSlot) GetParams() ([]Param, error) {
	params := make([]Param, 0)
	return params, nil
}
func (elm *ElementSlot) GetHtml() string        { return elm.Html }
func (elm *ElementSlot) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementSlot) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementSlot) GetType() string        { return elm.Type }
func (elm *ElementSlot) GetAttr() string        { return elm.Attr }
func (elm *ElementSlot) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementSlot) GetName() string        { return elm.Name }
func (elm *ElementSlot) GetProps() []Prop       { return elm.Props }
func (elm *ElementSlot) GetCompNames() []string { return elm.CompNames }
func (elm *ElementSlot) GetAttrs() []Attr       { return elm.Attrs }

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
		attr, err := NewAttr(a.Key, a.Val)
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

func (elm *ElementSlot) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================

// ##==================================================================
