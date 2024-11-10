package tag

import (
	"fmt"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	Info           TagInfo
	TrueInnerHtml  string
	FalseInnerHtml string
	AsParam        string
	AttrCondition  string
	AttrTag        string
}

func NewIfTagFromSelection(root *goquery.Selection, tag *goquery.Selection) (IfTag, error) {
	t := &IfTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, tag, "condition", "tag") },
		func() error { return t.setAttrs(tag) },
		func() error { return t.setTrueInnerHtml() },
		func() error { return t.setFalseInnerHtml() },
		func() error { return t.setAsParam() },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (t *IfTag) setTagInfo(root *goquery.Selection, tag *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, tag, attrsToExclude...)
	if err != nil {
		return err
	}
	t.Info = info
	return nil
}

func (t *IfTag) setAttrs(tag *goquery.Selection) error {
	condAttr, _ := tag.Attr("condition")
	tagAttr, _ := tag.Attr("tag")
	t.AttrCondition = condAttr
	t.AttrTag = tagAttr
	return nil
}

func (t *IfTag) setTrueInnerHtml() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	elseTags := s.Find("else")
	if elseTags.Length() > 1 {
		return fmt.Errorf("<for> tags may only have a single <else> within:\n\n%s", t.Info.Html)
	}
	elseTags.Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})
	htmlStr, err := gqpp.GetHtmlFromSelectionWithNewTag(s, t.AttrTag, t.Info.AttrStr)
	if err != nil {
		return err
	}
	htmlStr = parsley.RemoveEmptyLines(htmlStr)
	t.TrueInnerHtml = htmlStr
	return nil
}

func (t *IfTag) setFalseInnerHtml() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	elseTag := s.Find("else")
	if elseTag.Length() > 1 {
		return fmt.Errorf("<for> tags may only have a single <else> within:\n\n%s", t.Info.Html)
	}
	htmlStr, err := gqpp.GetHtmlFromSelectionWithNewTag(elseTag, t.AttrTag, t.Info.AttrStr)
	if err != nil {
		return err
	}
	t.FalseInnerHtml = htmlStr
	return nil
}

func (t *IfTag) setAsParam() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	cond, _ := s.Attr("condition")
	t.AsParam = cond + " bool"
	return nil
}

func (t IfTag) GetInfo() TagInfo { return t.Info }

func (t IfTag) TranspileToGo() (string, error) {
	trueHtml := t.TrueInnerHtml
	trueHtml = parsley.ReplaceFirstLine(trueHtml, parsley.TrimLeadingSpaces(parsley.GetFirstLine(trueHtml)))
	out := parsley.RemoveFirstLine(fmt.Sprintf(`
elseIf(%s, %s%s
%s,%s
%s
%s)    
`, t.AttrCondition, parsley.BackTick(), trueHtml, parsley.BackTick(), parsley.BackTick(), t.FalseInnerHtml, parsley.BackTick()))
	return out, nil
}
