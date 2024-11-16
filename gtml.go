package gtml

import (
	"fmt"
	"go/format"
	"gtml/internal/fungi"
	"gtml/internal/gqpp"
	"gtml/internal/purse"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ##==================================================================
type Element interface {
	GetSelection() *goquery.Selection
}

func GetChildElementList() []string {
	return []string{"_for"}
}

func NewElement(sel *goquery.Selection) (Element, error) {
	match := gqpp.GetFirstMatchingAttr(sel, "_component", "_for")
	switch match {
	case "_component":
		return NewComponentElement(sel), nil
	case "_for":
		return NewForElement(sel), nil
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("provided selection is not a valid element: %s", htmlStr)
}

func PrintElement(elm Element) {
	htmlStr, _ := GetElementHtml(elm)
	fmt.Println(purse.Flatten(htmlStr))
}

func GetElementHtml(elm Element) (string, error) {
	htmlStr, err := goquery.OuterHtml(elm.GetSelection())
	if err != nil {
		return "", err
	}
	return purse.Flatten(htmlStr), nil
}

func GetElementType(elm Element) string {
	match := gqpp.GetFirstMatchingAttr(elm.GetSelection(), "_component", "_for")
	switch match {
	case "_component":
		return "component"
	case "_for":
		return "for"
	}
	return ""
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
	htmlStr, err := GetElementHtml(elm)
	if err != nil {
		return nil, err
	}
	props := purse.ScanBetweenSubStrs(htmlStr, "{{", "}}")
	return props, nil
}

func GetPropValue(prop string) string {
	return purse.Squeeze(purse.RemoveAllSubStr(prop, "{{", "}}"))
}

func WalkElementProps(elm Element, fn func(prop string, val string) error) error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	for _, prop := range props {
		val := GetPropValue(prop)
		err := fn(prop, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func ForceElementAttr(elm Element, attrToCheck string) (string, error) {
	attr, exists := elm.GetSelection().Attr(attrToCheck)
	if !exists {
		htmlStr, err := GetElementHtml(elm)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("element is required to have the '%s' attribute: %s", attrToCheck, htmlStr)
	}
	return attr, nil
}

func ForceElementAttrParts(elm Element, attrToCheck string, partsExpected int) ([]string, error) {
	attr, err := ForceElementAttr(elm, attrToCheck)
	if err != nil {
		return make([]string, 0), nil
	}
	parts := strings.Split(attr, " ")
	if len(parts) != partsExpected {
		htmlStr, err := GetElementHtml(elm)
		if err != nil {
			return make([]string, 0), err
		}
		return make([]string, 0), fmt.Errorf("attribute '%s' expects %d distinct parts in element: %s", attrToCheck, partsExpected, htmlStr)
	}
	return parts, nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	htmlStr, err := GetElementHtml(elm)
	if err != nil {
		return "", err
	}
	err = WalkElementChildren(elm, func(child Element) error {
		childHtml, err := GetElementHtml(child)
		if err != nil {
			return err
		}
		goVar, err := NewGoVar(child)
		if err != nil {
			return err
		}
		varName, err := GetGoVarName(goVar)
		if err != nil {
			return err
		}
		htmlStr = strings.Replace(htmlStr, childHtml, fmt.Sprintf("%s.WriteString(%s)", builderName, varName), 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	err = WalkElementProps(elm, func(prop, val string) error {
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, val)
		htmlStr = strings.Replace(htmlStr, prop, call, 1)
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
type ComponentElement struct {
	Selection *goquery.Selection
}

func NewComponentElement(sel *goquery.Selection) *ComponentElement {
	elm := &ComponentElement{
		Selection: sel,
	}
	return elm
}

func (elm *ComponentElement) GetSelection() *goquery.Selection { return elm.Selection }

// ##==================================================================
type ForElement struct {
	Selection *goquery.Selection
}

func NewForElement(sel *goquery.Selection) *ForElement {
	elm := &ForElement{
		Selection: sel,
	}
	return elm
}

func (elm *ForElement) GetSelection() *goquery.Selection { return elm.Selection }

// ##==================================================================
type GoFunc interface {
	GetData() string
	SetData(str string)
	GetVars() []GoVar
}

func NewGoFunc(elm Element) (GoFunc, error) {
	if GetElementType(elm) == "component" {
		fn, err := NewGoComponentFunc(elm)
		if err != nil {
			return nil, err
		}
		return fn, nil
	}
	htmlStr, err := GetElementHtml(elm)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("provided element does not corrospond to a valid GoFunc: %s", htmlStr)
}

func PrintGoFunc(fn GoFunc) {
	fmt.Println(fn.GetData())
}

func GetGoFuncSkeleton() string {
	return purse.RemoveFirstLine(`
func NAME(PARAMS) RETURNTYPE {
BODY
RETURNVAL
}`)
}

func writeGoFuncSkeleton(fn GoFunc, indicator string, leaveIndicator bool, write string) {
	skeleton := fn.GetData()
	if leaveIndicator {
		skeleton = strings.Replace(skeleton, indicator, write+"\n"+indicator, 1)
	} else {
		skeleton = strings.Replace(skeleton, indicator, write, 1)
	}
	fn.SetData(skeleton)
}

func WriteGoFuncBody(fn GoFunc, leaveIndicator bool, write string) {
	writeGoFuncSkeleton(fn, "BODY", leaveIndicator, write)
}

func WriteGoFuncReturnVal(fn GoFunc, write string) {
	writeGoFuncSkeleton(fn, "RETURNVAL", false, write)
}

func WriteGoFuncReturnType(fn GoFunc, write string) {
	writeGoFuncSkeleton(fn, "RETURNTYPE", false, write)
}

func WriteGoFuncName(fn GoFunc, write string) {
	writeGoFuncSkeleton(fn, "NAME", false, write)
}

func WriteGoFuncParams(fn GoFunc, write string) {
	writeGoFuncSkeleton(fn, "PARAMS", false, write)
}

// ##==================================================================
type GoComponentFunc struct {
	Vars []GoVar
	Data string
}

func NewGoComponentFunc(elm Element) (*GoComponentFunc, error) {
	fn := &GoComponentFunc{}
	fn.Data = GetGoFuncSkeleton()
	WriteGoFuncBody(fn, true, "\tvar builder strings.Builder")
	WriteGoFuncReturnType(fn, "string")
	compAttr, err := ForceElementAttr(elm, "_component")
	if err != nil {
		return nil, err
	}
	WriteGoFuncName(fn, compAttr)
	err = WalkElementDirectChildren(elm, func(child Element) error {
		goVar, err := NewGoVar(child)
		if err != nil {
			return err
		}
		fn.Vars = append(fn.Vars, goVar)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func (fn *GoComponentFunc) GetData() string    { return fn.Data }
func (fn *GoComponentFunc) SetData(str string) { fn.Data = str }
func (fn *GoComponentFunc) GetVars() []GoVar   { return fn.Vars }

// ##==================================================================
type GoVar interface {
	GetData() string
}

func NewGoVar(elm Element) (GoVar, error) {
	match := GetElementType(elm)
	switch match {
	case "for":
		v, err := NewGoLoopVar(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	htmlStr, err := GetElementHtml(elm)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("element does not corrospond to a valid GoToken: %s", htmlStr)
}

func PrintGoVar(v GoVar) {
	fmt.Println(v.GetData())
}

func GetGoVarName(v GoVar) (string, error) {
	parts := strings.Split(v.GetData(), ":=")
	if len(parts) == 0 {
		return "", fmt.Errorf("GoVar does not contain a := symbol, so we cannot parse the name: %s", v.GetData())
	}
	firstPart := purse.Squeeze(parts[0])
	return firstPart, nil
}

// ##==================================================================
type GoLoopVar struct {
	Element     Element
	VarName     string
	BuilderName string
	Vars        []GoVar
	WriteVarsAs string
	Data        string
	IterItems   string
	IterItem    string
	IterType    string
}

func NewGoLoopVar(elm Element) (*GoLoopVar, error) {
	v := &GoLoopVar{
		Element: elm,
	}
	err := fungi.Process(
		func() error { return v.initBasicInfo() },
		func() error { return v.initVars() },
		func() error { return v.initWriteVarsAs() },
		func() error { return v.initData() },
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *GoLoopVar) GetData() string { return v.Data }

func (v *GoLoopVar) initBasicInfo() error {
	attrParts, err := ForceElementAttrParts(v.Element, "_for", 4)
	if err != nil {
		return err
	}
	v.VarName = attrParts[0] + "Loop"
	v.BuilderName = attrParts[0] + "Builder"
	v.IterItems = attrParts[2]
	v.IterItem = attrParts[0]
	v.IterType = purse.RemoveAllSubStr(attrParts[3], "[]")
	return nil
}

func (v *GoLoopVar) initVars() error {
	err := WalkElementDirectChildren(v.Element, func(child Element) error {
		innerVar, err := NewGoVar(child)
		if err != nil {
			return nil
		}
		v.Vars = append(v.Vars, innerVar)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *GoLoopVar) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += "\t" + inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoLoopVar) initData() error {
	htmlStr, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := collect(%s, func(i int, %s %s) string {
	var %s strings.Builder
%s
%s
	return %s.String()
})`, v.VarName, v.IterItems, v.IterItem, v.IterType, v.BuilderName, v.WriteVarsAs, htmlStr, v.BuilderName))
	code, err := format.Source([]byte(v.Data))
	if err != nil {
		return err
	}
	v.Data = string(code)
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
