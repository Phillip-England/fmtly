package gtml

import "github.com/PuerkitoBio/goquery"

type ComponentElement struct {
	Value    *goquery.Selection
	Children []Element
}
