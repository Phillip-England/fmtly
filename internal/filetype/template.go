package filetype

import (
	"fmt"
	"tagly/internal/gout"
)

type TemplateFile struct{}

func NewTemplateFileFromTaglyFile(f TaglyFile) (TemplateFile, error) {
	tf := &TemplateFile{}
	for _, fmtTag := range f.FmtTags {
		compFunc, err := gout.NewComponentFuncFromFmtTag(fmtTag)
		if err != nil {
			return *tf, err
		}
		fmt.Println(compFunc)
	}
	return *tf, nil
}
