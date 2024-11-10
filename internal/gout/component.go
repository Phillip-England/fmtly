package gout

import (
	"fmt"
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"
	"tagly/internal/tag"
)

type ComponentFunc struct {
	Name     string
	Shell    string
	FmtTag   *tag.FmtTag
	ParamStr string
}

func NewComponentFuncFromFmtTag(fmtTag tag.FmtTag) (ComponentFunc, error) {
	comp := &ComponentFunc{}
	err := fungi.Process(
		func() error { return comp.makeShell() },
		func() error { return comp.copyFmtTag(fmtTag) },
		func() error { return comp.setName() },
		func() error { return comp.setParamStr() },
		func() error { return comp.transpileTagsToGo() },
	)
	if err != nil {
		return *comp, err
	}
	return *comp, nil
}

func (comp *ComponentFunc) makeShell() error {
	comp.Shell = parsley.RemoveFirstLine(fmt.Sprintf(`
func NAME(PARAMS) string {
	return %s
		RETURN
	%s
}`, parsley.BackTick(), parsley.BackTick()))
	return nil
}

func (comp *ComponentFunc) copyFmtTag(fmtTag tag.FmtTag) error {
	copy := fmtTag
	comp.FmtTag = &copy
	return nil
}

func (comp *ComponentFunc) setName() error {
	comp.Name = comp.FmtTag.Name
	return nil
}

func (comp *ComponentFunc) setParamStr() error {
	paramStr := ""
	for _, prop := range comp.FmtTag.StrProps {
		newParam := prop.AsParam + ", "
		if strings.Contains(paramStr, newParam) {
			continue
		}
		paramStr += newParam
	}
	for _, tag := range comp.FmtTag.ForTags {
		if tag.AsParam == "" || len(tag.AsParam) == 0 {
			continue
		}
		paramStr += tag.AsParam + ", "
	}
	for _, tag := range comp.FmtTag.IfTags {
		paramStr += tag.AsParam + ", "
	}
	paramStr = paramStr[:len(paramStr)-2]
	comp.ParamStr = paramStr
	return nil
}

func (comp *ComponentFunc) transpileTagsToGo() error {
	clay := comp.FmtTag.Info.Html
	for {
		foundTag := false
		for _, t := range comp.FmtTag.Tags {
			goCode, err := t.TranspileToGo()
			if err != nil {
				return err
			}
			searchHtml := t.GetInfo().Html
			if strings.Contains(clay, searchHtml) {
				foundTag = true
			}
			clay = strings.Replace(clay, searchHtml, goCode, 1)
			newSel, err := gqpp.NewSelectionFromHtmlStr(clay)
			if err != nil {
				return err
			}
			tag.NewFmtTagFromSelection(newSel)
		}
		if foundTag == false {
			break
		}
	}
	fmt.Println(clay)
	return nil
}
