package tags

import (
	"fmtly/internal/parsley"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Fmt struct {
	Selection  *goquery.Selection
	Html       string
	HtmlOutput string
	Fors       []*For
}

func NewFmtFromStr(s string) (*Fmt, error) {
	t := &Fmt{}
	if err := t.setSelection(s); err != nil {
		return nil, err
	}
	if err := t.setHtml(); err != nil {
		return nil, err
	}
	if err := t.setForTags(); err != nil {
		return nil, err
	}
	if err := t.setForTagsHtmlOutput(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Fmt) setSelection(s string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return err
	}
	selection := doc.Find("fmt")
	t.Selection = selection
	return nil
}

func (t *Fmt) setHtml() error {
	htmlStr, err := goquery.OuterHtml(t.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	t.Html = flatHtml
	t.HtmlOutput = flatHtml
	return nil
}

func (t *Fmt) setForTags() error {
	var potErr error
	potErr = nil
	t.Selection.Find("for").Each(func(i int, s *goquery.Selection) {
		forTag, err := NewForFromSelection(s)
		if err != nil {
			potErr = err
			return
		}
		t.Fors = append(t.Fors, forTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (t *Fmt) setForTagsHtmlOutput() error {
	for _, ft := range t.Fors {
		if strings.Contains(t.HtmlOutput, ft.Html) {
			t.HtmlOutput = strings.Replace(t.HtmlOutput, ft.Html, ft.HtmlOutput, 1)
		}
	}
	parsley.Log(t.HtmlOutput)
	return nil
}
