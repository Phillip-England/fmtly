package gtml

import (
	"fmt"
	"gtml/internal/fungi"
	"gtml/internal/gqpp"
	"gtml/internal/purse"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	Selection    *goquery.Selection
	Html         string
	ForAttr      string
	ForAttrParts []string
	ItemName     string
	SliceName    string
	SliceType    string
	ItemType     string
	Depth        int
	Props        []Prop
	GoForLoop    GoForLoop
}

func NewForTagFromSelection(componentSel, forSel *goquery.Selection) (ForTag, error) {
	tag := &ForTag{
		Selection: forSel,
	}
	err := fungi.Process(
		func() error { return tag.setHtml() },
		func() error { return tag.setAttrDetails() },
		func() error { return tag.setDepth(componentSel) },
		func() error { return tag.setProps() },
		func() error { return tag.setGoForLoop() },
	)
	if err != nil {
		return *tag, err
	}
	return *tag, nil
}

func (tag *ForTag) setHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(tag.Selection)
	if err != nil {
		return err
	}
	tag.Html = purse.FlattenStr(htmlStr)
	return nil
}

func (tag *ForTag) setAttrDetails() error {
	forAttr, exists := tag.Selection.Attr("_for")
	if !exists {
		return fmt.Errorf("no _for='' found on: %s", purse.SnipStrAtIndex(tag.Html, 50))
	}
	parts := strings.Split(forAttr, " ")
	if len(parts) != 4 {
		return fmt.Errorf("_for='ITEM of SLICE SLICETYPE' format not followed: %s...", purse.SnipStrAtIndex(tag.Html, 50))
	}
	tag.ForAttrParts = parts
	tag.ForAttr = forAttr
	tag.ItemName = parts[0]
	tag.SliceName = parts[2]
	tag.SliceType = parts[3]
	tag.ItemType = purse.RemoveAllSubStr(parts[3], "[", "]")
	return nil
}

func (tag *ForTag) setDepth(componentSel *goquery.Selection) error {
	depth, err := gqpp.CalculateNodeDepth(componentSel, tag.Selection)
	if err != nil {
		return err
	}
	tag.Depth = depth
	return nil
}

func (tag *ForTag) setProps() error {
	props, err := NewPropsFromSelection(tag.Selection)
	if err != nil {
		return err
	}
	for _, prop := range props {
		if prop.Value == tag.ItemName {
			tag.Props = append(tag.Props, prop)
		}
		parts := strings.Split(prop.Value, ".")
		if strings.Count(prop.Value, ".") == 1 && len(parts) == 2 {
			firstPart := parts[0]
			if firstPart == tag.ItemName {
				tag.Props = append(tag.Props, prop)
			}
		}
	}
	return nil
}

func (tag *ForTag) setGoForLoop() error {
	loop, err := NewGoForLoop(tag.ItemName, tag.SliceName, tag.ItemType, tag.Props, tag.Html)
	if err != nil {
		return err
	}
	tag.GoForLoop = loop
	return nil
}
