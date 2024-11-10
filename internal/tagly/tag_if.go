package tagly

import (
	"fmt"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagIf struct {
	Info           TagInfo
	TrueInnerHtml  string
	FalseInnerHtml string
	AttrCondition  string
	AttrTag        string
}

func NewTagIfFromSelection(root *goquery.Selection, ogSel *goquery.Selection) (TagIf, error) {
	tag := &TagIf{}
	err := fungi.Process(
		func() error { return tag.setTagInfo(root, ogSel, "condition", "tag") },
		func() error { return tag.setAttrs(ogSel) },
		func() error { return tag.setTrueInnerHtml() },
		func() error { return tag.setFalseInnerHtml() },
	)
	if err != nil {
		return *tag, err
	}
	return *tag, nil
}

func (tag *TagIf) setTagInfo(root *goquery.Selection, ogSel *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, ogSel, attrsToExclude...)
	if err != nil {
		return err
	}
	tag.Info = info
	return nil
}

func (tag *TagIf) setAttrs(ogSel *goquery.Selection) error {
	condAttr, exists := ogSel.Attr("condition")
	if !exists {
		return fmt.Errorf("<if> tag requires a 'condition' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	tagAttr, exists := ogSel.Attr("tag")
	if !exists {
		tagAttr = "div"
	}
	tag.AttrCondition = condAttr
	tag.AttrTag = tagAttr
	return nil
}

func (tag *TagIf) setTrueInnerHtml() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	elseTags := sel.Find("else")
	if elseTags.Length() > 1 {
		return fmt.Errorf("<if> tags may only have a single <else> within:\n\n%s", tag.Info.Html)
	}
	elseTags.Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})
	htmlStr, err := sel.Html()
	if err != nil {
		return err
	}
	htmlStr = parsley.RemoveEmptyLines(htmlStr)
	tag.TrueInnerHtml = parsley.FlattenStr(htmlStr)
	return nil
}

func (tag *TagIf) setFalseInnerHtml() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	elseTag := sel.Find("else")
	if elseTag.Length() > 1 {
		return fmt.Errorf("<for> tags may only have a single <else> within:\n\n%s", tag.Info.Html)
	}
	htmlStr, err := elseTag.Html()
	if err != nil {
		return err
	}
	tag.FalseInnerHtml = parsley.FlattenStr(htmlStr)
	return nil
}
