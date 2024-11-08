package gqpp

import (
	"fmt"
	"fmtly/internal/parsley"
	"html"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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
