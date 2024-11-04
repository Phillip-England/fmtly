package tag

import (
	"fmt"
	"fmtly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	Info          *TagInfo
	ConditionAttr string
	TagAttr       string
}

func NewIfTagFromSelection(s *goquery.Selection) (*IfTag, error) {
	info, err := NewTagInfoFromSelection(s, "if", []string{"condition", "tag"})
	if err != nil {
		return nil, err
	}
	t := &IfTag{
		Info: info,
	}
	err = fungi.ProcessErrFuncs(
		t.setAttrs,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *IfTag) setAttrs() error {
	conditionAttr, exists := t.Info.Selection.Attr("condition")
	if !exists {
		return fmt.Errorf("<if> is missing 'condition' attribute:\n\n%s", t.Info.Html)
	}
	tagAttr, exists := t.Info.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<if> is missing 'tag' attribute:\n\n%s", t.Info.Html)
	}
	t.ConditionAttr = conditionAttr
	t.TagAttr = tagAttr
	return nil
}

func (t *IfTag) Html() string          { return t.Info.Html }
func (t *IfTag) Name() string          { return t.Info.Name }
func (t *IfTag) Scopes() []Tag         { return t.Info.Scopes }
func (t *IfTag) ParentTagName() string { return goquery.NodeName(t.Info.Selection.Parent()) }
