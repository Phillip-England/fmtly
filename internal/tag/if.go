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
	info, err := NewTagInfoFromSelection(s, "condition", "tag")
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
