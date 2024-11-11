package tagly

import (
	"fmt"
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"

	"github.com/PuerkitoBio/goquery"
)

type TagFmt struct {
	Name     string
	AttrName string
	AttrTag  string
	Info     TagInfo
	Tags     []Tag
	TagFors  []TagFor
	TagIfs   []TagIf
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
		func() error { return tag.extractTagFors() },
		func() error { return tag.extractTagIfs() },
		func() error { return tag.combineTags() },
		func() error { return tag.sortTagsByDepth() },
	)
	if err != nil {
		return *tag, err
	}

	return *tag, nil
}

func (tag TagFmt) AsStr() (string, error) {
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

func (tag TagFmt) GetInfo() TagInfo { return tag.Info }

func (tag TagFmt) TranspileToGo() (string, error) {
	clay := tag.CloneAsPointer()
	out := parsley.FlattenStr(tag.Info.Html)
	targetLength := len(clay.Tags)
	modifiedTags := make([]Tag, 0)
	for {
		if len(modifiedTags) == targetLength {
			break
		}
		for _, innerTag := range clay.Tags {
			innerHtml := parsley.FlattenStr(innerTag.GetInfo().Html)
			goCode, err := innerTag.TranspileToGo()
			if err != nil {
				return "", err
			}
			out = strings.Replace(out, innerHtml, goCode, 1)
			modifiedTags = append(modifiedTags, innerTag)
			newSel, err := gqpp.NewSelectionFromHtmlStr(out)
			if err != nil {
				return "", err
			}
			newFmt, err := NewTagFmtFromSelection(newSel)
			if err != nil {
				return "", err
			}
			clay = newFmt.CloneAsPointer()
		}
	}
	fmt.Println(out)
	return "", nil
}

func (tag *TagFmt) GetGoFuncParamStr() (string, error) {
	paramStr := ""
	for _, strProp := range tag.Info.StrProps {
		paramStr += strProp.FmtAsGoParam("string")
	}
	slice := strings.Split(paramStr, " ")
	filtered := parsley.RemoveDuplicatesInSlice(slice)
	paramStr = strings.Join(filtered, " ")
	for _, tagFor := range tag.TagFors {
		if !tagFor.IsRootFor {
			continue
		}
		paramStr += fmt.Sprintf("%s %s, ", tagFor.AttrAs, tagFor.AttrType)
	}
	for _, tagIf := range tag.TagIfs {
		paramStr += fmt.Sprintf("%s bool, ", tagIf.AttrCondition)
	}
	return paramStr, nil
}

func (tag TagFmt) CloneAsPointer() *TagFmt {
	tagCopy := tag
	tagPointer := &tagCopy
	return tagPointer
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

func (tag *TagFmt) combineTags() error {
	tags := make([]Tag, 0)
	for _, tagFor := range tag.TagFors {
		tags = append(tags, tagFor)
	}
	for _, tagIf := range tag.TagFors {
		tags = append(tags, tagIf)
	}
	return nil
}

func (tag *TagFmt) sortTagsByDepth() error {
	sorted := make([]Tag, 0)
	highestDepth := 0
	targetLength := len(tag.Tags)
	for _, innerTag := range tag.Tags {
		depth := innerTag.GetInfo().Depth
		if depth > highestDepth {
			highestDepth = depth
		}
	}
	for {
		for _, innerTag := range tag.Tags {
			depth := innerTag.GetInfo().Depth
			if depth == highestDepth {
				sorted = append(sorted, innerTag)
			}
		}
		highestDepth--
		if len(sorted) == targetLength {
			break
		}
	}
	return nil
}
