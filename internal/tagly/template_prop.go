package tagly

import (
	"strings"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TemplateProp struct {
	Raw   string
	Value string
}

func NewTemplatePropsFromSelection(sel *goquery.Selection) ([]TemplateProp, error) {
	htmlStr, err := gqpp.GetHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	props := parsley.ScanBetweenSubStrs(htmlStr, "{{", "}}")
	out := make([]TemplateProp, 0)
	for _, prop := range props {
		sq := parsley.Squeeze(prop)
		val := strings.Replace(sq, "{{", "", 1)
		val = strings.Replace(val, "}}", "", 1)
		out = append(out, TemplateProp{
			Raw:   prop,
			Value: val,
		})
	}
	return out, nil
}
