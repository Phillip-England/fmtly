package gotml

import (
	"gotml/internal/gqpp"
	"gotml/internal/purse"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Prop struct {
	Raw   string
	Value string
}

func NewPropsFromSelection(sel *goquery.Selection) ([]Prop, error) {
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	props := purse.ScanBetweenSubStrs(htmlStr, "{{", "}}")
	out := make([]Prop, 0)
	for _, prop := range props {
		sq := purse.Squeeze(prop)
		val := strings.Replace(sq, "{{", "", 1)
		val = strings.Replace(val, "}}", "", 1)
		out = append(out, Prop{
			Raw:   prop,
			Value: val,
		})
	}
	return out, nil
}
