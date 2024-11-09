package tag

import (
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

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

type ForProp struct {
	Raw   string
	Value string
}

func NewForPropsFromSelection(s *goquery.Selection) ([]ForProp, error) {
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return nil, err
	}
	typeAttr, _, err := gqpp.AttrFromStr(htmlStr, "type")
	if err != nil {
		return nil, err
	}
	if strings.Contains(typeAttr, "*") {
		typeAttr = strings.Replace(typeAttr, "*", "", 1)
	}
	asAttr, _, err := gqpp.AttrFromStr(htmlStr, "as")
	if err != nil {
		return nil, err
	}
	props := parsley.ScanBetweenSubStrs(htmlStr, "{{", "}}")
	out := make([]ForProp, 0)
	for _, prop := range props {
		sq := parsley.Squeeze(prop)
		val := strings.Replace(sq, "{{", "", 1)
		val = strings.Replace(val, "}}", "", 1)
		parts := strings.Split(val, ".")
		if len(parts) > 2 {
			continue
		}
		if val == asAttr && typeAttr == "string" {
			out = append(out, ForProp{
				Raw:   prop,
				Value: val,
			})
			continue
		}
		if len(parts) != 2 {
			continue
		}
		if parts[0] == asAttr {
			out = append(out, ForProp{
				Raw:   prop,
				Value: val,
			})
		}
	}
	return out, nil
}
