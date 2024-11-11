package tagly

import (
	"tagly/internal/fungi"
)

type FuncComponent struct {
	Name        string
	Code        string
	ParamStr    string
	OriginalFmt TagFmt
}

func NewFuncComponentFromTagFmt(ogFmt TagFmt) (FuncComponent, error) {
	comp := &FuncComponent{}
	err := fungi.Process(
		func() error { return comp.setOrignalFmt(ogFmt) },
		func() error { return comp.setName() },
		func() error { return comp.setParamStr() },
		func() error { return comp.generateGoFunc() },
	// func() error { return comp.captureReturnHtml() },
	)
	if err != nil {
		return *comp, err
	}
	return *comp, nil
}

func (comp *FuncComponent) setOrignalFmt(ogFmt TagFmt) error {
	comp.OriginalFmt = ogFmt
	return nil
}

func (comp *FuncComponent) setName() error {
	comp.Name = comp.OriginalFmt.Name
	return nil
}

func (comp *FuncComponent) setParamStr() error {
	paramStr, err := comp.OriginalFmt.GetGoFuncParamStr()
	if err != nil {
		return err
	}
	comp.ParamStr = paramStr
	return nil
}

func (comp *FuncComponent) generateGoFunc() error {
	_, err := comp.OriginalFmt.TranspileToGo()
	if err != nil {
		return err
	}
	return nil
}
