package tag

import (
	"tagly/internal/fungi"
	"tagly/internal/gqpp"

	"github.com/PuerkitoBio/goquery"
)

type FmtTag struct {
	Info    TagInfo
	ForTags []ForTag
	IfTags  []IfTag
}

func NewFmtTagsFromFilePath(path string) ([]FmtTag, error) {
	s, err := gqpp.NewSelectionFromFilePath(path)
	if err != nil {
		return nil, err
	}
	out := make([]FmtTag, 0)
	var potErr error
	potErr = nil
	s.Find("fmt").Each(func(i int, fmtSel *goquery.Selection) {
		fmtTag, err := NewFmtTagFromSelection(fmtSel)
		if err != nil {
			potErr = err
			return
		}
		out = append(out, fmtTag)
	})
	if potErr != nil {
		return nil, err
	}
	return out, nil
}

func NewFmtTagFromSelection(s *goquery.Selection) (FmtTag, error) {
	t := &FmtTag{}
	err := fungi.Process(
		func() error { return t.setTagInfo(s, s, "name", "tag") },
		func() error { return t.extractForTags() },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (t *FmtTag) setTagInfo(root *goquery.Selection, tag *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, tag, attrsToExclude...)
	if err != nil {
		return err
	}
	t.Info = info
	return nil
}

func (t *FmtTag) extractForTags() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	var potErr error
	potErr = nil
	s.Find("for").Each(func(i int, forSel *goquery.Selection) {
		forTag, err := NewForTagFromSelection(s, forSel)
		if err != nil {
			potErr = err
		}
		t.ForTags = append(t.ForTags, forTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (t *FmtTag) extractIfTags() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	var potErr error
	potErr = nil
	s.Find("if").Each(func(i int, forSel *goquery.Selection) {
		ifTag, err := NewIfTagFromSelection(s, forSel)
		if err != nil {
			potErr = err
		}
		t.IfTags = append(t.IfTags, ifTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
