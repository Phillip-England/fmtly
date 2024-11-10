package tag

import (
	"fmt"
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type FmtTag struct {
	Name     string
	Info     TagInfo
	StrProps []StrProp
	ForTags  []ForTag
	IfTags   []IfTag
	Tags     []Tag
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
		func() error { return t.setName() },
		func() error { return t.extractStrProps(s) },
		func() error { return t.extractForTags() },
		func() error { return t.extractIfTags() },
		func() error { return t.combineAllTags() },
		func() error { return t.sortTagsByDepth() },
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

func (t *FmtTag) setName() error {
	s, err := gqpp.NewSelectionFromHtmlStr(t.Info.Html)
	if err != nil {
		return err
	}
	name, exists := s.Attr("name")
	if !exists {
		return fmt.Errorf("<fmt> tags require an attr:\n\n%s", t.Info.Html)
	}
	t.Name = name
	return nil
}

func (t *FmtTag) extractStrProps(s *goquery.Selection) error {
	props, err := NewStrPropsFromSelection(s)
	if err != nil {
		return err
	}
	t.StrProps = props
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

func (t *FmtTag) combineAllTags() error {
	tags := make([]Tag, 0)
	for _, tag := range t.ForTags {
		tags = append(tags, tag)
	}
	for _, tag := range t.IfTags {
		tags = append(tags, tag)
	}
	t.Tags = tags
	return nil
}

func (t *FmtTag) sortTagsByDepth() error {
	maxDepth := 0
	for _, tag := range t.Tags {
		depth := tag.GetInfo().Depth
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	sorted := make([]Tag, 0)
	for {
		if len(sorted) == len(t.Tags) {
			break
		}
		for _, tag := range t.Tags {
			depth := tag.GetInfo().Depth
			if depth == maxDepth {
				sorted = append(sorted, tag)
			}
		}
		maxDepth--
	}
	t.Tags = sorted
	return nil
}

type StrProp struct {
	Raw     string
	Value   string
	AsParam string
}

func NewStrPropsFromSelection(s *goquery.Selection) ([]StrProp, error) {
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return nil, err
	}
	props := parsley.ScanBetweenSubStrs(htmlStr, "{{", "}}")
	out := make([]StrProp, 0)
	for _, prop := range props {
		sq := parsley.Squeeze(prop)
		val := strings.Replace(sq, "{{", "", 1)
		val = strings.Replace(val, "}}", "", 1)
		parts := strings.Split(val, ".")
		if len(parts) != 1 {
			continue
		}
		propIsForProp := false
		s.Find("for").Each(func(i int, forSel *goquery.Selection) {
			asAttr, _ := forSel.Attr("as")
			typeAttr, _ := forSel.Attr("type")
			if val == asAttr && typeAttr == "string" {
				propIsForProp = true
			}
		})
		if propIsForProp {
			continue
		}
		out = append(out, StrProp{
			Raw:     prop,
			Value:   val,
			AsParam: val + " string",
		})
	}
	return out, nil
}
