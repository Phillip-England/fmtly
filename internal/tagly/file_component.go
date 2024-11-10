package tagly

import (
	"tagly/internal/fungi"
)

type FileComponent struct {
	SrcPath        string
	FuncComponents []FuncComponent
}

func NewFileComponentFromFileTagly(fTagly FileTagly) (FileComponent, error) {
	tempFile := &FileComponent{
		SrcPath: fTagly.Path,
	}
	err := fungi.Process(
		func() error { return tempFile.generateComponentFuncs(fTagly) },
	)
	if err != nil {
		return *tempFile, err
	}

	return *tempFile, nil
}

func (fComp *FileComponent) generateComponentFuncs(fTagly FileTagly) error {
	for _, fmtTag := range fTagly.TagFmts {
		compFunc, err := NewFuncComponentFromTagFmt(fmtTag)
		if err != nil {
			return err
		}
		fComp.FuncComponents = append(fComp.FuncComponents, compFunc)
	}
	return nil
}
