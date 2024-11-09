package tag

import (
	"strings"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

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
