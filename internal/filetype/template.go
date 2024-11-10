package filetype

import (
	"tagly/internal/fungi"
	"tagly/internal/gout"
)

type TemplateFile struct {
	SrcPath        string
	ComponentFuncs []gout.ComponentFunc
}

func NewTemplateFileFromTaglyFile(tagFile TaglyFile) (TemplateFile, error) {
	tempFile := &TemplateFile{
		SrcPath: tagFile.Path,
	}
	err := fungi.Process(
		func() error { return tempFile.generateComponentFuncs(tagFile) },
	)
	if err != nil {
		return *tempFile, err
	}
	return *tempFile, nil
}

func (tf *TemplateFile) generateComponentFuncs(tagFile TaglyFile) error {
	for _, fmtTag := range tagFile.FmtTags {
		compFunc, err := gout.NewComponentFuncFromFmtTag(fmtTag)
		if err != nil {
			return err
		}
		tf.ComponentFuncs = append(tf.ComponentFuncs, compFunc)
	}
	return nil
}
