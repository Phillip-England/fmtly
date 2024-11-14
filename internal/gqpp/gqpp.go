package gqpp

import (
	"fmt"
	"gtml/internal/purse"
	"html"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func NewSelectionFromFilePath(path string) (*goquery.Selection, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fStr := string(f)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
	if err != nil {
		return nil, err
	}
	body := doc.Find("body")
	return body, nil
}

func NewSelectionFromHtmlStr(htmlStr string) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}
	out := doc.Find("body").Children()
	return out, nil
}

func ChangeSelectionTagName(s *goquery.Selection, tagName string) (*goquery.Selection, error) {
	attrStr := GetAttrStr(s)
	htmlStr, err := s.Html()
	if err != nil {
		return nil, err
	}
	out := ""
	if len(attrStr) == 0 {
		out = fmt.Sprintf("<%s>%s</%s>", tagName, htmlStr, tagName)
	} else {
		out = fmt.Sprintf("<%s %s>%s</%s>", tagName, attrStr, htmlStr, tagName)
	}
	newSel, err := NewSelectionFromHtmlStr(out)
	if err != nil {
		return nil, err
	}
	return newSel, nil
}

func GetAttrStr(selection *goquery.Selection, filter ...string) string {
	var attrs []string
	filterMap := make(map[string]struct{})
	for _, f := range filter {
		filterMap[f] = struct{}{}
	}
	selection.Each(func(i int, sel *goquery.Selection) {
		for _, attr := range sel.Nodes[0].Attr {
			if _, found := filterMap[attr.Key]; !found {
				attrs = append(attrs, fmt.Sprintf(`%s="%s"`, attr.Key, attr.Val))
			}
		}
	})
	return strings.Join(attrs, " ")
}

func NewHtmlFromSelection(s *goquery.Selection) (string, error) {
	htmlStr, err := goquery.OuterHtml(s)
	if err != nil {
		return "", err
	}
	htmlStr = html.UnescapeString(htmlStr)
	lastLine := purse.GetLastLine(htmlStr)
	leadingSpaces := purse.CountLeadingSpaces(lastLine)
	lines := purse.MakeLines(htmlStr)
	for i, line := range lines {
		if i == 0 {
			lines[i] = strings.Repeat(" ", leadingSpaces) + line
		}
	}
	return purse.JoinLines(lines), nil
}

func ClimbTreeUntil(s *goquery.Selection, cond func(parent *goquery.Selection) bool) error {
	parent := s.Parent()
	if cond(parent) {
		return nil
	}
	return ClimbTreeUntil(parent, cond)
}

func AttrFromStr(str string, attrName string) (string, bool, error) {
	s, err := NewSelectionFromHtmlStr(str)
	if err != nil {
		return "", false, err
	}
	out, exists := s.Attr(attrName)
	return out, exists, nil
}

func CalculateNodeDepth(root *goquery.Selection, child *goquery.Selection) (int, error) {
	depth := 0
	childNodeName := goquery.NodeName(child)
	childHtml, err := NewHtmlFromSelection(child)
	if err != nil {
		return -1, err
	}
	rootHtml, err := NewHtmlFromSelection(root)
	if err != nil {
		return -1, err
	}
	var potErr error
	root.Find(childNodeName).Each(func(i int, search *goquery.Selection) {
		searchHtml, err := NewHtmlFromSelection(search)
		if err != nil {
			potErr = err
			return
		}
		if searchHtml == childHtml {
			ClimbTreeUntil(search, func(parent *goquery.Selection) bool {
				if parent.Length() == 0 {
					potErr = fmt.Errorf("child node: %s not found within parent node: %s", childHtml[0:30], rootHtml[0:30])
				}
				parentHtml, err := NewHtmlFromSelection(parent)
				if err != nil {
					potErr = err
					return true
				}
				if parentHtml == rootHtml {
					return true
				}
				depth++
				return false
			})
		}
	})
	if potErr != nil {
		return -1, potErr
	}
	return depth, nil
}

func CountMatchingParentTags(root, child *goquery.Selection, tagNames ...string) (int, error) {
	count := 0
	tagSet := make(map[string]struct{})
	for _, tag := range tagNames {
		tagSet[tag] = struct{}{}
	}
	childHtml, err := NewHtmlFromSelection(child)
	if err != nil {
		return -1, err
	}
	found := false
	var potentialErr error
	root.Find(goquery.NodeName(child)).EachWithBreak(func(i int, search *goquery.Selection) bool {
		searchHtml, err := NewHtmlFromSelection(search)
		if err != nil {
			potentialErr = err
			return false
		}
		if searchHtml == childHtml {
			found = true
			current := search.Parent()
			for current.Length() > 0 {
				nodeName := goquery.NodeName(current)
				if _, exists := tagSet[nodeName]; exists {
					count++
				}
				current = current.Parent()
			}
			return false
		}
		return true
	})
	if potentialErr != nil {
		return -1, potentialErr
	}
	if !found {
		return -1, fmt.Errorf("child node not found within the root")
	}
	return count, nil
}

func NewHtmlFromSelectionWithNewTag(s *goquery.Selection, newTagName string, newTagAttrStr string) (string, error) {
	htmlStr, err := s.Html()
	if err != nil {
		return "", err
	}
	openTag := ""
	if newTagAttrStr == "" {
		openTag = fmt.Sprintf("<%s>", newTagName)
	} else {
		openTag = fmt.Sprintf("<%s %s>", newTagName, newTagAttrStr)
	}
	closeTag := fmt.Sprintf("</%s>", newTagName)
	out := fmt.Sprintf("%s%s%s", openTag, htmlStr, closeTag)
	newSel, err := NewSelectionFromHtmlStr(out)
	if err != nil {
		return "", err
	}
	finalOut, err := NewHtmlFromSelection(newSel)
	if err != nil {
		return "", err
	}
	return finalOut, nil
}

func FindDeepestMatchingSelection(selection *goquery.Selection, selectors ...string) (*goquery.Selection, bool) {
	var deepestSelection *goquery.Selection
	maxDepth := -1

	for _, selector := range selectors {
		found := selection.Find(selector)
		if found.Length() > 0 {
			for i := 0; i < found.Length(); i++ {
				node := found.Eq(i)

				// Calculate the depth of this node within the original selection
				depth, err := CalculateNodeDepth(selection, node)
				if err != nil {
					return nil, false
				}

				// Update the deepest node found if this one is deeper
				if depth > maxDepth {
					maxDepth = depth
					deepestSelection = node
				}
			}
		}
	}

	if deepestSelection == nil {
		return nil, false
	}
	return deepestSelection, true
}

func HasMatchingElements(selection *goquery.Selection, selectors ...string) bool {
	for _, selector := range selectors {
		if selection.Find(selector).Length() > 0 {
			return true
		}
	}
	return false
}

func GetFirstMatchingAttr(selection *goquery.Selection, attrs ...string) string {
	for _, attr := range attrs {
		if _, exists := selection.Attr(attr); exists {
			return attr
		}
	}
	return ""
}

func GetAttrPart(selection *goquery.Selection, attrName string, part int) (string, error) {
	attr, exists := selection.Attr(attrName)
	if !exists {
		return "", fmt.Errorf("attr: '%s' does not exist when it should", attrName)
	}
	parts := strings.Split(attr, " ")
	if len(parts) < part {
		return "", fmt.Errorf("attr: '%s' must contain %d parts", attrName, part)
	}
	return parts[part], nil
}

func GetAttr(selection *goquery.Selection, attrName string) (string, error) {
	attr, exists := selection.Attr(attrName)
	if !exists {
		return "", fmt.Errorf("attr: '%s' does not exist", attrName)
	}
	return attr, nil
}

func HasAttr(selection *goquery.Selection, attrs ...string) bool {
	for _, attr := range attrs {
		if _, exists := selection.Attr(attr); exists {
			return true
		}
	}
	return false
}

func HasParentWithAttrs(sel *goquery.Selection, stopAt *goquery.Selection, attrs ...string) bool {
	// Create a set of the attribute names for quick lookup
	attrSet := make(map[string]struct{})
	for _, attr := range attrs {
		attrSet[attr] = struct{}{}
	}

	// Traverse up the parent hierarchy
	current := sel.Parent()
	for current.Length() > 0 {
		// Stop if we reach the specified stopAt selection
		if current.IsSelection(stopAt) {
			break
		}

		for _, node := range current.Nodes {
			for _, attr := range node.Attr {
				if _, found := attrSet[attr.Key]; found {
					return true
				}
			}
		}
		current = current.Parent()
	}

	return false
}
