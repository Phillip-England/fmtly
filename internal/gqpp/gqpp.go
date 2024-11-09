package gqpp

import (
	"fmt"
	"html"
	"os"
	"strings"
	"tagly/internal/parsley"

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

func GetHtmlFromSelection(s *goquery.Selection) (string, error) {
	htmlStr, err := goquery.OuterHtml(s)
	if err != nil {
		return "", err
	}
	htmlStr = html.UnescapeString(htmlStr)
	lastLine := parsley.GetLastLine(htmlStr)
	leadingSpaces := parsley.CountLeadingSpaces(lastLine)
	lines := parsley.MakeLines(htmlStr)
	for i, line := range lines {
		if i == 0 {
			lines[i] = strings.Repeat(" ", leadingSpaces) + line
		}
	}
	return parsley.JoinLines(lines), nil
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
	childHtml, err := GetHtmlFromSelection(child)
	if err != nil {
		return -1, err
	}
	rootHtml, err := GetHtmlFromSelection(root)
	if err != nil {
		return -1, err
	}
	var potErr error
	root.Find(childNodeName).Each(func(i int, search *goquery.Selection) {
		searchHtml, err := GetHtmlFromSelection(search)
		if err != nil {
			potErr = err
			return
		}
		if searchHtml == childHtml {
			ClimbTreeUntil(search, func(parent *goquery.Selection) bool {
				parentHtml, err := GetHtmlFromSelection(parent)
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
	childHtml, err := GetHtmlFromSelection(child)
	if err != nil {
		return -1, err
	}
	found := false
	var potentialErr error
	root.Find(goquery.NodeName(child)).EachWithBreak(func(i int, search *goquery.Selection) bool {
		searchHtml, err := GetHtmlFromSelection(search)
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
