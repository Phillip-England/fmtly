package tag

import (
	"fmt"
	"fmtly/internal/fungi"
	"fmtly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagInfo struct {
	Name      string
	Selection *goquery.Selection
	Html      string
	AttrStr   string
	Scopes    []Tag
}

func NewTagInfoFromSelection(s *goquery.Selection, name string, setAttrs []string) (*TagInfo, error) {
	info := &TagInfo{
		Selection: s,
		Name:      name,
	}
	if err := fungi.Process(
		func() error { return info.setHtml() },
		func() error { return info.setAttrStr(setAttrs...) },
		func() error { return info.setScopes() },
	); err != nil {
		return nil, err
	}
	return info, nil
}

func (info *TagInfo) setHtml() error {
	htmlStr, err := goquery.OuterHtml(info.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	info.Html = flatHtml
	return nil
}

func (info *TagInfo) setAttrStr(filterAttrs ...string) error {
	attrStr := ""
	for i, node := range info.Selection.Nodes {
		if i == 0 {
			for _, attr := range node.Attr {
				if parsley.EqualsOneof(attr.Key, filterAttrs...) {
					continue
				}
				attrStr += fmt.Sprintf("%s=\"%s\" ", attr.Key, attr.Val)
			}
		}
	}
	if len(attrStr) != 0 {
		attrStr = attrStr[:len(attrStr)-1]
	}
	info.AttrStr = attrStr
	return nil
}

func (info *TagInfo) setScopes() error {
	tags, err := NewTagsFromSelection(info.Selection)
	if err != nil {
		return err
	}
	info.Scopes = tags
	return nil
}
