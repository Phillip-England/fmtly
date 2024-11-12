package tagly

import (
	"fmt"
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagFor struct {
	Info      TagInfo
	AttrTag   string
	AttrType  string
	AttrAs    string
	AttrIn    string
	IsRootFor bool
}

func NewTagForFromSelection(root *goquery.Selection, ogSel *goquery.Selection) (TagFor, error) {
	t := &TagFor{}
	err := fungi.Process(
		func() error { return t.setTagInfo(root, ogSel, "in", "as", "tag", "type") },
		func() error { return t.setAttrs(ogSel) },
		func() error { return t.setIsRootFor(root, ogSel) },
	)
	if err != nil {
		return *t, err
	}
	return *t, nil
}

func (tag TagFor) AsStr() (string, error) {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return "", err
	}
	newTag, err := gqpp.GetHtmlFromSelectionWithNewTag(sel, tag.AttrTag, tag.Info.AttrStr)
	if err != nil {
		return "", err
	}
	for _, prop := range tag.Info.Props {
		newTag = strings.Replace(newTag, prop.Raw, "`+"+prop.Value+"+`", 1)
	}
	return newTag, nil
}

func (tag TagFor) GetInfo() TagInfo { return tag.Info }

func (tag TagFor) TranspileToGo() (string, error) {
	htmlStr, err := tag.AsStr()
	if err != nil {
		return "", err
	}
	htmlStr = parsley.WrapStr(htmlStr, "`", "`")
	parts := strings.Split(htmlStr, "+")
	buildParts := make([]string, 0)
	for _, part := range parts {
		buildParts = append(buildParts, fmt.Sprintf(`%sBuilder.WriteString(%s)`, tag.AttrAs, part))
	}
	htmlStr = strings.Join(buildParts, " ")
	goCode := fmt.Sprintf(`builder.WriteString(%scollectStr(%s, func(i int, %s %s) string { var %sBuilder strings.Builder %s return %sBuilder.String() })%s)`, parsley.BackTick(), tag.AttrIn, tag.AttrAs, tag.AttrType, tag.AttrAs, htmlStr, tag.AttrAs, parsley.BackTick())
	return parsley.FlattenStr(goCode), nil
}

func (tag *TagFor) setTagInfo(root *goquery.Selection, ogSel *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, ogSel, attrsToExclude...)
	if err != nil {
		return err
	}
	tag.Info = info
	return nil
}

func (tag *TagFor) setAttrs(ogSel *goquery.Selection) error {
	inAttr, exists := ogSel.Attr("in")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'in' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	typeAttr, exists := ogSel.Attr("type")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'type' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	asAttr, exists := ogSel.Attr("as")
	if !exists {
		return fmt.Errorf("<for> tag requires a 'type' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	tagAttr, exists := ogSel.Attr("tag")
	if !exists {
		tagAttr = "div"
	}
	tag.AttrIn = inAttr
	tag.AttrType = typeAttr
	tag.AttrAs = asAttr
	tag.AttrTag = tagAttr
	return nil
}

func (tag *TagFor) setIsRootFor(root *goquery.Selection, ogSel *goquery.Selection) error {
	rootHtml, err := gqpp.GetHtmlFromSelection(root)
	if err != nil {
		return err
	}
	rootHtml = parsley.FlattenStr(rootHtml)
	ogSelHtml, err := gqpp.GetHtmlFromSelection(ogSel)
	if err != nil {
		return err
	}
	ogSelHtml = parsley.FlattenStr(ogSelHtml)
	var potErr error
	potErr = nil
	root.Find("for").Each(func(i int, forToCheck *goquery.Selection) {
		forToCheckHtml, err := gqpp.GetHtmlFromSelection(forToCheck)
		if err != nil {
			potErr = err
			return
		}
		forToCheckHtml = parsley.FlattenStr(forToCheckHtml)
		if ogSelHtml == forToCheckHtml {
			forCount := 0
			err := gqpp.ClimbTreeUntil(forToCheck, func(parent *goquery.Selection) bool {
				parentHtml, err := gqpp.GetHtmlFromSelection(parent)
				if err != nil {
					potErr = err
					return true
				}
				parentHtml = parsley.FlattenStr(parentHtml)
				if parentHtml == rootHtml {
					return true
				}
				parentNodeName := goquery.NodeName(parent)
				if parentNodeName == "for" {
					forCount++
				}
				return false
			})
			if forCount == 0 {
				tag.IsRootFor = true
			}
			if err != nil {
				potErr = err
				return
			}
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
