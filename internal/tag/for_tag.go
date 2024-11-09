package tag

import (
	"tagly/internal/fungi"
	"tagly/internal/gqpp"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	Info     TagInfo
	ForProps []ForProp
	IsRoot   bool
	AsParam  string
}

func NewForTagFromSelection(root *goquery.Selection, tag *goquery.Selection) (ForTag, error) {
	t := &ForTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, tag, "in", "as", "tag", "type") },
		func() error { return t.extractForProps(tag) },
		func() error { return t.checkIfRoot(root, tag) },
		func() error { return t.setAsParam(tag) },
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

func (t *ForTag) setAsParam(tag *goquery.Selection) error {
	inAttr, _ := tag.Attr("in")
	typeAttr, _ := tag.Attr("type")
	if t.IsRoot {
		t.AsParam = inAttr + " []" + typeAttr
	}
	return nil
}

func (t *ForTag) extractForProps(s *goquery.Selection) error {
	props, err := NewForPropsFromSelection(s)
	if err != nil {
		return err
	}
	t.ForProps = props
	return nil
}

func (t *ForTag) checkIfRoot(root *goquery.Selection, tag *goquery.Selection) error {
	forCount, err := gqpp.CountMatchingParentTags(root, tag, "for")
	if err != nil {
		return err
	}
	if forCount == 0 {
		t.IsRoot = true
	}
	return nil
}
