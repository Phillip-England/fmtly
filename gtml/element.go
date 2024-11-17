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
type Element interface {
	GetSelection() *goquery.Selection
	GetParam() (string, error)
}

func GetChildElementList() []string {
	return []string{"_for"}
}

func NewElement(sel *goquery.Selection) (Element, error) {
	match := gqpp.GetFirstMatchingAttr(sel, "_component", "_for")
	switch match {
	case "_component":
		return NewElementComponent(sel), nil
	case "_for":
		return NewElementFor(sel), nil
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
	if match == "" {
		return ""
	}
	return strings.Replace(match, "_", "", 1)
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
	err = WalkElementStrProps(elm, func(prop, val string) error {
		strProp := val + " " + "string"
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

func WalkElementStrProps(elm Element, fn func(prop string, val string) error) error {
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
	err = WalkElementStrProps(elm, func(prop, val string) error {
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, val)
		htmlStr = strings.Replace(htmlStr, prop, call, 1)
		return nil
	})
	fmt.Println(htmlStr)
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
}

func NewElementComponent(sel *goquery.Selection) *ElementComponent {
	elm := &ElementComponent{
		Selection: sel,
	}
	return elm
}

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParam() (string, error)        { return "", nil }

// ##==================================================================
type ElementFor struct {
	Selection *goquery.Selection
}

func NewElementFor(sel *goquery.Selection) *ElementFor {
	elm := &ElementFor{
		Selection: sel,
	}
	return elm
}

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParam() (string, error) {
	parts, err := ForceElementAttrParts(elm, "_for", 4)
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
