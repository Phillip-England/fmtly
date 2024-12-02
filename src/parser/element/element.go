package element

import (
	"fmt"
	"gtml/src/parser/attr"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/gqpp"
)

type Element interface {
	GetSelection() *goquery.Selection
	SetHtml(htmlStr string)
	GetHtml() string
	Print()
	GetType() string
	GetAttr() string
	GetAttrParts() []string
	GetCompNames() []string
	GetAttrs() []attr.Attr
	GetName() string
	GetId() string
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
		elm, err := NewComponent(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementFor:
		elm, err := NewFor(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementIf:
		elm, err := NewIf(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementElse:
		elm, err := NewElse(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementPlaceholder:
		elm, err := NewPlaceholder(htmlStr, sel, compNames)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementSlot:
		elm, err := NewSlot(htmlStr, sel, compNames)
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

func CollectElementDirectChildren(sel *goquery.Selection, ogChildren []Element, compNames []string) ([]Element, error) {
	var potErr error
	sel.Children().Each(func(i int, childSel *goquery.Selection) {
		childSelIsElement := gqpp.HasAttr(childSel, GetChildElementList()...)
		if childSelIsElement {
			childHtml, err := gqpp.NewHtmlFromSelection(childSel)
			if err != nil {
				potErr = err
				return
			}
			childElm, err := NewElement(childHtml, compNames)
			if err != nil {
				potErr = err
				return
			}
			ogChildren = append(ogChildren, childElm)
		} else {
			children, err := CollectElementDirectChildren(childSel, ogChildren, compNames)
			if err != nil {
				potErr = err
				return
			}
			ogChildren = children
		}
	})
	if potErr != nil {
		return ogChildren, potErr
	}
	return ogChildren, nil
}

func WalkElementDirectChildren(elm Element, fn func(child Element) error) error {
	childElms, err := CollectElementDirectChildren(elm.GetSelection(), make([]Element, 0), elm.GetCompNames())
	if err != nil {
		return err
	}
	for _, childElm := range childElms {
		err = fn(childElm)
		if err != nil {
			return err
		}
	}

	// var potErr error
	// elm.GetSelection().Children().Each(func(i int, childSel *goquery.Selection) {
	// 	if gqpp.HasAttr(childSel, GetChildElementList()...) {
	// 		childHtml, err := gqpp.NewHtmlFromSelection(childSel)
	// 		if err != nil {
	// 			potErr = err
	// 			return
	// 		}
	// 		childElm, err := NewElement(childHtml, elm.GetCompNames())
	// 		if err != nil {
	// 			potErr = err
	// 			return
	// 		}
	// 		err = fn(childElm)
	// 		if err != nil {
	// 			potErr = err
	// 			return
	// 		}
	// 	}
	// })
	// if potErr != nil {
	// 	return potErr
	// }
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

// think func will need testing and improvement
func ExtractComponentStringsFromFile(fStr string) ([]string, error) {
	compStrs := make([]string, 0)
	clay := fStr
	parts := make([]string, 0)
	for {
		index := strings.Index(clay, "_component=")
		if index == -1 {
			break
		}
		part := clay[:index+len("_component=")]
		parts = append(parts, part)
		clay = clay[index+len("_component="):]
	}
	parts = append(parts, clay)
	filtered := make([]string, 0)
	for i, p := range parts {
		if i%2 == 1 {
			index := strings.LastIndex(p, "<")
			first := p[:index]
			second := p[index:]
			filtered = append(filtered, first)
			filtered = append(filtered, second)
			continue
		}
		filtered = append(filtered, p)
	}
	for i, p := range filtered {
		if i%2 == 1 {
			lastPart := filtered[i-1]
			compStr := lastPart + p
			compStrs = append(compStrs, compStr)
		}
	}
	return compStrs, nil
}

func ReadComponentSelectionsFromFile(path string) ([]*goquery.Selection, error) {
	selections := make([]*goquery.Selection, 0)
	f, err := os.ReadFile(path)
	if err != nil {
		return selections, err
	}
	fStr := string(f)
	compStrs, err := ExtractComponentStringsFromFile(fStr)
	if err != nil {
		return selections, err
	}
	for _, compStr := range compStrs {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(compStr))
		if err != nil {
			return selections, err
		}
		doc.Find("*").Each(func(i int, sel *goquery.Selection) {
			_, exists := sel.Attr(KeyElementComponent)
			if exists {
				selections = append(selections, sel)
			}
		})
	}
	return selections, nil
}

func ConvertSelectionsIntoElements(selections []*goquery.Selection, compNames []string) ([]Element, error) {
	elms := make([]Element, 0)
	for _, sel := range selections {
		htmlStr, err := gqpp.NewHtmlFromSelection(sel)
		if err != nil {
			return elms, err
		}
		elm, err := NewElement(htmlStr, compNames)
		if err != nil {
			return elms, err
		}
		elms = append(elms, elm)
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

func MarkSelectionPlaceholders(sel *goquery.Selection, compNames []string) error {
	ogSelHtml, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return err
	}
	err = MarkSelectionAsPlaceholder(sel, compNames, ogSelHtml)
	if err != nil {
		return err
	}
	var potErr error
	sel.Find("*").Each(func(i int, inner *goquery.Selection) {
		if potErr != nil {
			return // Exit early if there's already an error
		}
		potErr = MarkSelectionAsPlaceholder(inner, compNames, ogSelHtml)
	})
	return potErr
}

func MarkSelectionAsPlaceholder(inner *goquery.Selection, compNames []string, ogSelHtml string) error {
	innerNodeName := goquery.NodeName(inner)
	for _, compName := range compNames {
		if strings.ToLower(compName) == innerNodeName {
			inner.SetAttr("_placeholder", compName)
			// var potErr error
			// inner.Children().Each(func(i int, childSel *goquery.Selection) {
			// 	_, hasSlot := childSel.Attr("_slot")
			// 	if !hasSlot {
			// 		potErr = fmt.Errorf("_placeholder element has children which are not wrapped in an element with a _slot='slotName' attribute: %s", ogSelHtml)
			// 		return
			// 	}
			// })
			// if potErr != nil {
			// 	return potErr
			// }
		}
	}
	return nil
}

func MarkElementPlaceholders(elm Element) (Element, error) {
	clay := elm.GetHtml()
	err := WalkAllElementNodesIncludingRoot(elm, func(sel *goquery.Selection) error {
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

func MarkSelectionAsUnique(sel *goquery.Selection) {
	id := 0
	sel.SetAttr("_id", strconv.Itoa(id))
	id++
	sel.Find("*").Each(func(i int, inner *goquery.Selection) {
		match := gqpp.GetFirstMatchingAttr(inner, GetChildElementList()...)
		if match == "" {
			return // skip elements which don't have a valid _attribute
		}
		idStr := strconv.Itoa(id)
		inner.SetAttr("_id", idStr)
		id++
	})
}

func MarkSelectionsAsUnique(selections []*goquery.Selection) {
	for _, sel := range selections {
		MarkSelectionAsUnique(sel)
	}
}
