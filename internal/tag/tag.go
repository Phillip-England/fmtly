package tag

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Tag interface {
	Info() *TagInfo
	MakeGoOutput() (string, error)
	Out() string
}

func NewTagFromSelection(s *goquery.Selection) (Tag, error) {
	tagName := goquery.NodeName(s)
	switch tagName {
	case "for":
		tag, err := NewForTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return tag, nil
	case "else":
		tag, err := NewElseTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return tag, nil
	case "if":
		tag, err := NewIfTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return tag, nil
	case "fmt":
		tag, err := NewFmtTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return tag, nil
	default:
		return nil, fmt.Errorf("tag name: '%s' is not valid when creating a Tag interface{}", tagName)
	}
}

func NewTagsFromSelection(selection *goquery.Selection) ([]Tag, error) {
	tagNames := []string{"for", "if", "fmt"}
	tags := make([]Tag, 0)
	for _, name := range tagNames {
		var potErr error
		selection.Find(name).Each(func(i int, s *goquery.Selection) {
			tag, err := NewTagFromSelection(s)
			if err != nil {
				potErr = err
				return
			}
			tags = append(tags, tag)
		})
		if potErr != nil {
			return nil, potErr
		}
	}
	return tags, nil
}

func NewTagsFromHtmlStr(s string) ([]Tag, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	children := doc.Find("body").Children()
	tags, err := NewTagsFromSelection(children)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
