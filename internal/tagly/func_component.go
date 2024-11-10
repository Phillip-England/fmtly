package tagly

import (
	"tagly/internal/fungi"
)

type FuncComponent struct {
	Name     string
	Code     string
	ParamStr string
}

func NewFuncComponentFromTagFmt(fmtTag TagFmt) (FuncComponent, error) {
	comp := &FuncComponent{}
	err := fungi.Process(
	// func() error { return comp.copyFmtTag(fmtTag) },
	// func() error { return comp.setName() },
	// func() error { return comp.setParamStr() },
	// func() error { return comp.transpileTagsToGo() },
	// func() error { return comp.captureReturnHtml() },
	)
	if err != nil {
		return *comp, err
	}
	return *comp, nil
}
