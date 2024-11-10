package gout

import (
	"strings"
	"tagly/internal/fungi"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"
	"tagly/internal/tag"
)

type ComponentFunc struct {
	Name       string
	Code       string
	FmtTag     *tag.FmtTag
	ParamStr   string
	ReturnHtml string
}

func NewComponentFuncFromFmtTag(fmtTag tag.FmtTag) (ComponentFunc, error) {
	comp := &ComponentFunc{}
	err := fungi.Process(
		func() error { return comp.copyFmtTag(fmtTag) },
		func() error { return comp.setName() },
		func() error { return comp.setParamStr() },
		func() error { return comp.transpileTagsToGo() },
		func() error { return comp.captureReturnHtml() },
	)
	if err != nil {
		return *comp, err
	}
	return *comp, nil
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
	clay := parsley.FlattenStr(comp.FmtTag.Info.Html)
	for {
		newSel, err := gqpp.NewSelectionFromHtmlStr(clay)
		if err != nil {
			return err
		}
		newFmt, err := tag.NewFmtTagFromSelection(newSel)
		if err != nil {
			return err
		}
		if len(newFmt.Tags) == 0 {
			break
		}
		for _, t := range newFmt.Tags {
			goCode, err := t.TranspileToGo()
			if err != nil {
				return err
			}
			goCode = parsley.FlattenStr(goCode)
			tagHtml := parsley.FlattenStr(t.GetInfo().Html)
			if strings.Contains(clay, tagHtml) {
				clay = strings.Replace(clay, tagHtml, "`+"+goCode+"+`", 1)
			}
		}
	}
	finalSel, err := gqpp.NewSelectionFromHtmlStr(clay)
	if err != nil {
		return nil
	}
	finalFmt, err := tag.NewFmtTagFromSelection(finalSel)
	if err != nil {
		return err
	}
	finalCompFunc, err := finalFmt.TranspileToGo()
	if err != nil {
		return err
	}
	finalCompFunc = strings.Replace(finalCompFunc, "PARAMS", comp.ParamStr, 1)
	comp.Code = finalCompFunc
	return nil
}

func (comp *ComponentFunc) captureReturnHtml() error {
	indexOfReturn := strings.Index(comp.Code, "return")
	sliced := comp.Code[indexOfReturn:len(comp.Code)]
	indexOfFirstBackTick := strings.Index(sliced, "`")
	indexOfFinalBackTick := strings.LastIndex(sliced, "`")
	sliced = sliced[indexOfFirstBackTick+1 : indexOfFinalBackTick]
	sel, err := gqpp.NewSelectionFromHtmlStr(sliced)
	if err != nil {
		return err
	}
	htmlStr, err := gqpp.GetHtmlFromSelection(sel)
	if err != nil {
		return err
	}
	comp.ReturnHtml = htmlStr
	return nil

}
