package tag

import (
	"fmt"
	"fmtly/internal/fungi"
	"fmtly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagInfo struct {
	Selection *goquery.Selection
	Html      string
	AttrStr   string
}

func NewTagInfoFromSelection(s *goquery.Selection, specialAttrs ...string) (*TagInfo, error) {
	in := &TagInfo{
		Selection: s,
	}
	err := fungi.ProcessErrFuncs(
		in.setHtml,
	)
	if err != nil {
		return nil, err
	}
	if err := in.setAttrStr(specialAttrs...); err != nil {
		return nil, err
	}
	return in, nil
}

func (in *TagInfo) setHtml() error {
	htmlStr, err := goquery.OuterHtml(in.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	in.Html = flatHtml
	return nil
}

func (in *TagInfo) setAttrStr(filterAttrs ...string) error {
	attrStr := ""
	for i, node := range in.Selection.Nodes {
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
	in.AttrStr = attrStr
	return nil
}
