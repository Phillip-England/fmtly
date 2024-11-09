package tag

import (
	"tagly/internal/fungi"
	"tagly/internal/gqpp"

	"github.com/PuerkitoBio/goquery"
)

type TagInfo struct {
	Html    string
	Depth   int
	AttrStr string
}

func NewTagInfoFromSelection(root *goquery.Selection, tag *goquery.Selection, attrsToExclude ...string) (TagInfo, error) {
	tagInfo := &TagInfo{}
	err := fungi.Process(
		func() error { return tagInfo.setHtml(tag) },
		func() error { return tagInfo.setDepth(root) },
		func() error { return tagInfo.setAttrStr(attrsToExclude...) },
	)
	if err != nil {
		return *tagInfo, err
	}
	return *tagInfo, nil
}

func (t *TagInfo) setHtml(s *goquery.Selection) error {
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return err
	}
	t.Html = htmlStr
	return nil
}

func (t *TagInfo) setDepth(root *goquery.Selection) error {
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

func (t *TagInfo) setAttrStr(attrsToExclude ...string) error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Html)
	if err != nil {
		return err
	}
	attrStr := gqpp.GetAttrStr(s, attrsToExclude...)
	t.AttrStr = attrStr
	return nil
}
