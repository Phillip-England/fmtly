package comp

import "github.com/PuerkitoBio/goquery"

type Comp struct {
	Doc      *goquery.Document
	Html     string
	Name     string
	Props    []*Prop
	ForProps []*ForProp
	ForTags  []*ForTag
	IfTags   []*IfTag
	FmtTag   *FmtTag
	Output   string
}

type Prop struct {
	Raw   string
	Value string
}

type ForProp struct {
	Raw   string
	Value string
}

type ForTag struct {
	Selection *goquery.Selection
	Depth     int
	Html      string
	AttrStr   string
	AttrIn    string
	AttrTag   string
	AttrType  string
	AttrAs    string
}

type FmtTag struct {
	Selection *goquery.Selection
	Html      string
	AttrStr   string
	AttrName  string
	AttrTag   string
}

type IfTag struct {
	Selection     *goquery.Selection
	Depth         int
	Html          string
	AttrCondition string
	AttrTag       string
	ElseTag       *ElseTag
}

type ElseTag struct {
	Selection *goquery.Selection
	Html      string
}
