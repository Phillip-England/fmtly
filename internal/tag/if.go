package tag

import (
	"fmt"
	"fmtly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	Selection       *goquery.Selection
	Html            string
	AttrStr         string
	ParamOutput     string
	ConditionAttr   string
	TagAttr         string
	ElseTag         *ElseTag
	TrueHtmlOutput  string
	FalseHtmlOutput string
	HtmlOutput      string
}

type ElseTag struct {
	Selection *goquery.Selection
	Html      string
}

func NewElseTagFromSelection(s *goquery.Selection) (*ElseTag, error) {
	t := &ElseTag{
		Selection: s,
	}
	if err := t.setHtml(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *ElseTag) setHtml() error {
	htmlStr, err := t.Selection.Html()
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	t.Html = flatHtml
	return nil
}

func NewIfTagFromSelection(s *goquery.Selection) (*IfTag, error) {
	t := &IfTag{
		Selection: s,
	}
	if err := t.setHtml(); err != nil {
		return nil, err
	}
	if err := t.setAttrStr(); err != nil {
		return nil, err
	}
	if err := t.setAttrs(); err != nil {
		return nil, err
	}
	if err := t.setElseTag(); err != nil {
		return nil, err
	}
	if err := t.setTrueHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setFalseHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setHtmlOutput(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *IfTag) setHtml() error {
	htmlStr, err := goquery.OuterHtml(t.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	t.Html = flatHtml
	return nil
}

func (t *IfTag) setAttrStr() error {
	attrStr := ""
	for i, node := range t.Selection.Nodes {
		if i == 0 {
			for _, attr := range node.Attr {
				if parsley.EqualsOneof(attr.Key, "condition", "tag") {
					continue
				}
				attrStr += fmt.Sprintf("%s=\"%s\" ", attr.Key, attr.Val)
			}
		}
	}
	if len(attrStr) != 0 {
		attrStr = attrStr[:len(attrStr)-1]
	}
	t.AttrStr = attrStr
	return nil
}

func (t *IfTag) setAttrs() error {
	conditionAttr, exists := t.Selection.Attr("condition")
	if !exists {
		return fmt.Errorf("<if> is missing 'condition' attribute:\n\n%s", t.Html)
	}
	tagAttr, exists := t.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<if> is missing 'tag' attribute:\n\n%s", t.Html)
	}
	t.ConditionAttr = conditionAttr
	t.TagAttr = tagAttr
	return nil
}

func (t *IfTag) setElseTag() error {
	s := t.Selection.Find("else")
	secondElse := s.Find("else")
	if secondElse == nil {
		return fmt.Errorf("<if> tag has more than 1 <else>:\n\n%s", t.Html)
	}
	elseTag, err := NewElseTagFromSelection(s)
	if err != nil {
		return err
	}
	t.ElseTag = elseTag
	return nil
}

func (t *IfTag) setTrueHtmlOutput() error {
	sCopy := t.Selection
	sCopy.Find("else").Remove()
	htmlStr, err := sCopy.Html()
	if err != nil {
		return err
	}
	t.TrueHtmlOutput = fmt.Sprintf(`<%s %s>%s</%s>`, t.TagAttr, t.AttrStr, parsley.FlattenStr(htmlStr), t.TagAttr)
	return nil
}

func (t *IfTag) setFalseHtmlOutput() error {
	t.FalseHtmlOutput = fmt.Sprintf(`<%s %s>%s</%s>`, t.TagAttr, t.AttrStr, t.ElseTag.Html, t.TagAttr)
	return nil
}

func (t *IfTag) setHtmlOutput() error {
	t.HtmlOutput = fmt.Sprintf(`ifElse(%s, "%s", "%s")`, t.ConditionAttr, t.TrueHtmlOutput, t.FalseHtmlOutput)
	return nil
}
