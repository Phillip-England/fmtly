package tag

import (
	"fmt"
	"fmtly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type FmtTag struct {
	Info     *TagInfo
	NameAttr string
	TagAttr  string
}

func NewFmtTagFromSelection(selection *goquery.Selection) (*FmtTag, error) {
	info, err := NewTagInfoFromSelection(selection, "fmt", []string{"name", "tag"})
	if err != nil {
		return nil, err
	}
	t := &FmtTag{
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

func (t *FmtTag) setAttrs() error {
	nameAttr, exists := t.Info.Selection.Attr("name")
	if !exists {
		return fmt.Errorf("<fmt> is missing 'name' attribute:\n\n%s", t.Info.Html)
	}
	tagAttr, exists := t.Info.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<fmt> is missing 'tag' attribute:\n\n%s", t.Info.Html)
	}
	t.NameAttr = nameAttr
	t.TagAttr = tagAttr
	return nil
}

func (t *FmtTag) Html() string          { return t.Info.Html }
func (t *FmtTag) Name() string          { return t.Info.Name }
func (t *FmtTag) Scopes() []Tag         { return t.Info.Scopes }
func (t *FmtTag) ParentTagName() string { return goquery.NodeName(t.Info.Selection.Parent()) }
