package tag

import (
	"tagly/internal/fungi"
	"tagly/internal/gqpp"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	Html  string
	Depth int
}

func NewIfTagFromSelection(root *goquery.Selection, s *goquery.Selection) (IfTag, error) {
	t := &IfTag{}
	err := fungi.Process(
		func() error { return t.setHtml(s) },
		func() error { return t.setDepth(root) },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (t *IfTag) setHtml(s *goquery.Selection) error {
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return err
	}
	t.Html = htmlStr
	return nil
}

func (t *IfTag) setDepth(root *goquery.Selection) error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Html)
	if err != nil {
		return err
	}
	d, err := gqpp.CalculateNodeDepth(root, s)
	if err != nil {
		return err
	}
	t.Depth = d
	return nil
}
