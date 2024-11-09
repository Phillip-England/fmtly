package tag

import (
	"tagly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	Info TagInfo
}

func NewIfTagFromSelection(root *goquery.Selection, tag *goquery.Selection) (IfTag, error) {
	t := &IfTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, tag, "condition", "tag") },
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
