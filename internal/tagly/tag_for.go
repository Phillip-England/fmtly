package tagly

import (
	"fmt"
	"tagly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type TagFor struct {
	Info     TagInfo
	AttrTag  string
	AttrType string
	AttrAs   string
	AttrIn   string
}

func NewTagForFromSelection(root *goquery.Selection, ogSel *goquery.Selection) (TagFor, error) {
	t := &TagFor{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, ogSel, "in", "as", "tag", "type") },
		func() error { return t.setAttrs(ogSel) },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (tag *TagFor) setTagInfo(root *goquery.Selection, ogSel *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, ogSel, attrsToExclude...)
	if err != nil {
		return err
	}
	tag.Info = info
	return nil
}

func (tag *TagFor) setAttrs(ogSel *goquery.Selection) error {
	inAttr, exists := ogSel.Attr("in")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'in' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	typeAttr, exists := ogSel.Attr("type")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'type' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	asAttr, exists := ogSel.Attr("as")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'type' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	tagAttr, exists := ogSel.Attr("tag")
	if !exists {
		tagAttr = "div"
	}
	tag.AttrIn = inAttr
	tag.AttrType = typeAttr
	tag.AttrAs = asAttr
	tag.AttrTag = tagAttr
	return nil
}
