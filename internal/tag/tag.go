package tag

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Tag interface {
	Html() string
	Name() string
	Scopes() []Tag
	ParentTagName() string
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
	tagNames := []string{"for", "if", "else", "fmt"}
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

func NewFmtTagsFromDir(dir string) ([]Tag, error) {
	var fmtTags []Tag
	err := filepath.Walk("./components", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fStr := string(f)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
		if err != nil {
			return err
		}
		var potErr error
		potErr = nil
		doc.Find("fmt").Each(func(i int, s *goquery.Selection) {
			ft, err := NewTagFromSelection(s)
			if err != nil {
				potErr = err
				return
			}
			fmtTags = append(fmtTags, ft)
		})
		if potErr != nil {
			return potErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fmtTags, nil
}
