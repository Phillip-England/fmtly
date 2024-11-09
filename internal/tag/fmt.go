package tag

import (
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type FmtTag struct {
	Info     TagInfo
	StrProps []StrProp
	ForTags  []ForTag
	IfTags   []IfTag
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
		func() error { return t.extractStrProps(s) },
		func() error { return t.extractForTags() },
		func() error { return t.extractIfTags() },
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

// func (t *FmtTag) setParamStr() error {
// 	paramStr := ""
// 	for _, prop := range t.StrProps {
// 		newParam := prop.AsParam + ", "
// 		if strings.Contains(paramStr, newParam) {
// 			continue
// 		}
// 		paramStr += newParam
// 	}
// 	for _, tag := range t.ForTags {
// 		if tag.AsParam == "" || len(tag.AsParam) == 0 {
// 			continue
// 		}
// 		paramStr += tag.AsParam + ", "
// 	}
// 	for _, tag := range t.IfTags {
// 		paramStr += tag.AsParam + ", "
// 	}
// 	paramStr = paramStr[:len(paramStr)-2]
// 	t.ParamStr = paramStr
// 	return nil
// }
