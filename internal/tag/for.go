package tag

import (
	"fmt"
	"fmtly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	Info     *TagInfo
	InAttr   string
	AsAttr   string
	TypeAttr string
	TagAttr  string
}

func NewForTagFromSelection(selection *goquery.Selection) (*ForTag, error) {
	info, err := NewTagInfoFromSelection(selection, "in", "as", "tag", "")
	if err != nil {
		return nil, err
	}
	t := &ForTag{
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

func (t *ForTag) setAttrs() error {
	inAttr, exists := t.Info.Selection.Attr("in")
	if !exists {
		return fmt.Errorf("<for> is missing 'in' attribute:\n\n%s", t.Info.Html)
	}
	asAttr, exists := t.Info.Selection.Attr("as")
	if !exists {
		return fmt.Errorf("<for> is missing 'as' attribute:\n\n%s", t.Info.Html)
	}
	typeAttr, exists := t.Info.Selection.Attr("type")
	if !exists {
		typeAttr = "any"
	}
	tagAttr, exists := t.Info.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<for> is missing 'tag' attribute:\n\n%s", t.Info.Html)
	}
	t.InAttr = inAttr
	t.AsAttr = asAttr
	t.TypeAttr = typeAttr
	t.TagAttr = tagAttr
	return nil
}
