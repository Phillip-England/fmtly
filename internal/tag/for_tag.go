package tag

import (
	"tagly/internal/fungi"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	TagInfo TagInfo
}

func NewForTagFromSelection(root *goquery.Selection, tag *goquery.Selection) (ForTag, error) {
	t := &ForTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, tag) },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (t *ForTag) setTagInfo(root *goquery.Selection, tag *goquery.Selection) error {
	info, err := NewTagInfoFromSelection(root, tag)
	if err != nil {
		return err
	}
	t.TagInfo = info
	return nil
}
