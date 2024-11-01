package comp

import "github.com/PuerkitoBio/goquery"

type Comp struct {
	Doc     *goquery.Document
	Html    string
	Name    string
	ForTags []*ForTag
}

type ForTag struct {
	Selection *goquery.Selection
	Depth     int
	Html      string
}
