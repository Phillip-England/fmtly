package tagly

import (
	"fmt"
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"

	"github.com/PuerkitoBio/goquery"
)

type TagFmt struct {
	Name         string
	AttrName     string
	AttrTag      string
	Info         TagInfo
	Tags         []Tag
	Props        []TemplateProp
	StrProps     []TemplateProp
	ForStrProps  []TemplateProp
	ForTypeProps []TemplateProp
	ChildProp    TemplateProp
	TagFors      []TagFor
	TagIfs       []TagIf
}

func NewTagFmtsFromFilePath(path string) ([]TagFmt, error) {
	sel, err := gqpp.NewSelectionFromFilePath(path)
	if err != nil {
		return nil, err
	}
	out := make([]TagFmt, 0)
	var potErr error
	potErr = nil
	sel.Find("fmt").Each(func(i int, fmtSel *goquery.Selection) {
		fmtTag, err := NewTagFmtFromSelection(fmtSel)
		if err != nil {
			potErr = err
			return
		}
		out = append(out, fmtTag)
	})
	if potErr != nil {
		return nil, err
	}
	return out, nil
}

func NewTagFmtFromSelection(sel *goquery.Selection) (TagFmt, error) {
	tag := &TagFmt{}
	err := fungi.Process(
		func() error { return tag.setTagInfo(sel, sel, "name", "tag") },
		func() error { return tag.setAttrs(sel) },
		func() error { return tag.setName() },
		func() error { return tag.extractAllProps() },
		func() error { return tag.sortProps() },
		func() error { return tag.extractTagFors() },
		func() error { return tag.extractTagIfs() },
	)
	if err != nil {
		return *tag, err
	}

	return *tag, nil
}

func (tag *TagFmt) setTagInfo(root *goquery.Selection, selfSelection *goquery.Selection, attrsToExclude ...string) error {
	info, err := NewTagInfoFromSelection(root, selfSelection, attrsToExclude...)
	if err != nil {
		return err
	}
	tag.Info = info
	return nil
}

func (tag *TagFmt) setAttrs(s *goquery.Selection) error {
	nameAttr, exists := s.Attr("name")
	if !exists {
		return fmt.Errorf("<fmt> tag requires a 'name' attribute:\n\n%s", tag.Info.Html[0:50]+"...")
	}
	tagAttr, exists := s.Attr("tag")
	if !exists {
		tagAttr = "div"
	}
	tag.AttrTag = tagAttr
	tag.AttrName = nameAttr
	return nil
}

func (tag *TagFmt) setName() error {
	tag.Name = tag.AttrName
	return nil
}

func (tag *TagFmt) extractAllProps() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	props, err := NewTemplatePropsFromSelection(sel)
	if err != nil {
		return err
	}
	tag.Props = props
	return nil
}

func (tag *TagFmt) sortProps() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	strProps := make([]TemplateProp, 0)
	forTypeProps := make([]TemplateProp, 0)
	forStrProps := make([]TemplateProp, 0)
	childrenProps := make([]TemplateProp, 0)
	for _, prop := range tag.Props {
		shouldBreak := false
		sel.Find("for").Each(func(i int, forSel *goquery.Selection) {
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
		if prop.Value == "...children" {
			childrenProps = append(childrenProps, prop)
			continue
		}
		strProps = append(strProps, prop)
	}

	if len(childrenProps) > 1 {
		return fmt.Errorf("only one {{ children... }} prop allowed per <fmt> component:\n\n%s", tag.Info.Html)
	}
	tag.ForStrProps = forStrProps
	tag.ForTypeProps = forTypeProps
	tag.StrProps = strProps
	if len(childrenProps) == 1 {
		tag.ChildProp = childrenProps[0]
	}
	return nil
}

func (tag *TagFmt) extractTagFors() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	var potErr error
	potErr = nil
	sel.Find("for").Each(func(i int, forSel *goquery.Selection) {
		forTag, err := NewTagForFromSelection(sel, forSel)
		if err != nil {
			potErr = err
		}
		tag.TagFors = append(tag.TagFors, forTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (tag *TagFmt) extractTagIfs() error {
	sel, err := gqpp.NewSelectionFromHtmlStr(tag.Info.Html)
	if err != nil {
		return err
	}
	var potErr error
	potErr = nil
	sel.Find("if").Each(func(i int, forSel *goquery.Selection) {
		ifTag, err := NewTagIfFromSelection(sel, forSel)
		if err != nil {
			potErr = err
		}
		tag.TagIfs = append(tag.TagIfs, ifTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
