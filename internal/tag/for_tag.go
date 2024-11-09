package tag

import (
	"tagly/internal/fungi"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	Info     TagInfo
	ForProps []string
}

func NewForTagFromSelection(root *goquery.Selection, tag *goquery.Selection) (ForTag, error) {
	t := &ForTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, tag, "in", "as", "tag", "type") },
		func() error { return t.extractForProps() },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (t *ForTag) setTagInfo(root *goquery.Selection, tag *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, tag, attrsToExclude...)
	if err != nil {
		return err
	}
	t.Info = info
	return nil
}

func (t *ForTag) extractForProps() error {
	props := parsley.ScanBetweenSubStrs(t.Info.Html, "{{", "}}")
	t.ForProps = props
	return nil
}
