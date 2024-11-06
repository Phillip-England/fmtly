package tag

import (
	"fmtly/internal/fungi"
	"strings"
)

type GoOutput struct {
	Tag  Tag
	Html string
}

func NewOutputFromTag(tag Tag) (*GoOutput, error) {
	out := &GoOutput{
		Tag:  tag,
		Html: tag.Info().Html,
	}
	err := fungi.Process(
		func() error { return out.setOutput() },
	)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (out *GoOutput) setOutput() error {
	tags, err := NewTagsFromHtmlStr(out.Html)
	if err != nil {
		return err
	}
	if len(tags) == 0 {
		return nil
	}
	for _, tag := range tags {
		output, err := tag.MakeGoOutput()
		if err != nil {
			return err
		}
		out.Html = strings.Replace(out.Html, tag.Info().Html, output, 1)
	}
	return out.setOutput()
}
