package tag

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type Tag interface{}

func NewTagFromSelection(s *goquery.Selection) (Tag, error) {
	tagName := goquery.NodeName(s)
	switch tagName {
	case "for":
		forTag, err := NewForTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return forTag, nil
	case "else":
		elseTag, err := NewElseTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return elseTag, err
	case "if":
		ifTag, err := NewIfTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return ifTag, err
	case "fmt":
		fmtTag, err := NewFmtTagFromSelection(s)
		if err != nil {
			return nil, err
		}
		return fmtTag, err
	default:
		return nil, fmt.Errorf("tag name '%s' is not valid.\nValid tag names are: for, if, else, fmt", tagName)
	}
}

func NewTagsFromSelection(s *goquery.Selection) ([]*Tag, error) {
	var tags []*Tag
	tagNames := []string{"for", "if", "else", "fmt"}
	for _, tagName := range tagNames {
		search := s.Find(tagName)
		tag, err := NewTagFromSelection(search)
		if err != nil {
			return nil, err
		}
		tags = append(tags, &tag)
	}
	return tags, nil
}
