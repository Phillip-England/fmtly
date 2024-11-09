package tag

import (
	"strings"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

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
