package tagly

import (
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagInfo struct {
	Html         string
	Depth        int
	AttrStr      string
	Props        []TemplateProp
	StrProps     []TemplateProp
	ForStrProps  []TemplateProp
	ForTypeProps []TemplateProp
}

func NewTagInfoFromSelection(root *goquery.Selection, ogSel *goquery.Selection, attrsToExclude ...string) (TagInfo, error) {
	tagInfo := &TagInfo{}
	err := fungi.Process(
		func() error { return tagInfo.setHtml(ogSel) },
		func() error { return tagInfo.setDepth(root) },
		func() error { return tagInfo.setAttrStr(attrsToExclude...) },
		func() error { return tagInfo.extractAllProps(ogSel) },
		func() error { return tagInfo.sortProps(ogSel) },
	)
	if err != nil {
		return *tagInfo, err
	}
	return *tagInfo, nil
}

func (info *TagInfo) setHtml(ogSel *goquery.Selection) error {
	htmlStr, err := gqpp.GetHtmlFromSelection(ogSel)
	if err != nil {
		return err
	}
	info.Html = parsley.FlattenStr(htmlStr)
	return nil
}

func (info *TagInfo) setDepth(root *goquery.Selection) error {
	sel, err := gqpp.NewSelectionFromHtmlStr(info.Html)
	if err != nil {
		return err
	}
	depth, err := gqpp.CalculateNodeDepth(root, sel)
	if err != nil {
		return err
	}
	info.Depth = depth
	return nil
}

func (info *TagInfo) setAttrStr(attrsToExclude ...string) error {
	s, err := gqpp.NewSelectionFromHtmlStr(info.Html)
	if err != nil {
		return err
	}
	attrStr := gqpp.GetAttrStr(s, attrsToExclude...)
	info.AttrStr = attrStr
	return nil
}

func (info *TagInfo) extractAllProps(ogSel *goquery.Selection) error {
	props, err := NewTemplatePropsFromSelection(ogSel)
	if err != nil {
		return err
	}
	info.Props = props
	return nil
}

func (info *TagInfo) sortProps(ogSel *goquery.Selection) error {
	strProps := make([]TemplateProp, 0)
	forTypeProps := make([]TemplateProp, 0)
	forStrProps := make([]TemplateProp, 0)
	for _, prop := range info.Props {
		shouldBreak := false
		ogSel.Find("for").Each(func(i int, forSel *goquery.Selection) {
			asAttr, _ := forSel.Attr("as")
			if prop.Value == asAttr {
				forStrProps = append(forStrProps, prop)
				shouldBreak = true
			}
		})
		if shouldBreak {
			break
		}
		if strings.Count(prop.Raw, ".") == 1 {
			forTypeProps = append(forTypeProps, prop)
			continue
		}
		strProps = append(strProps, prop)
	}
	info.ForStrProps = forStrProps
	info.ForTypeProps = forTypeProps
	info.StrProps = strProps
	return nil
}
