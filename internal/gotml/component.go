package gotml

import (
	"fmt"
	"gotml/internal/fungi"
	"gotml/internal/gqpp"
	"gotml/internal/purse"

	"github.com/PuerkitoBio/goquery"
)

type Component struct {
	Name      string
	Selection *goquery.Selection
	Html      string
	ForTags   []ForTag
}

func NewComponentFromSelection(sel *goquery.Selection) (Component, error) {
	comp := &Component{
		Selection: sel,
	}
	err := fungi.Process(
		func() error { return comp.setHtml() },
		func() error { return comp.setName() },
		func() error { return comp.setForTags() },
	)
	if err != nil {
		return *comp, err
	}
	return *comp, nil
}

func (comp *Component) setHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(comp.Selection)
	if err != nil {
		return err
	}
	comp.Html = purse.FlattenStr(htmlStr)
	return nil
}

func (comp *Component) setName() error {
	nameAttr, exists := comp.Selection.Attr("_name")
	if !exists {
		return fmt.Errorf("no _name='' found on: %s...", purse.SnipStrAtIndex(comp.Html, 30))
	}
	comp.Name = nameAttr
	return nil
}

func (comp *Component) setForTags() error {
	var potErr error
	comp.Selection.Find("*").Each(func(i int, sel *goquery.Selection) {
		_, exists := sel.Attr("_for")
		if exists {
			forTag, err := NewForTagFromSelection(comp.Selection, sel)
			if err != nil {
				potErr = err
				return
			}
			comp.ForTags = append(comp.ForTags, forTag)
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
