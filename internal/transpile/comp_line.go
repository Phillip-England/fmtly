package transpile

import (
	"fmtly/internal/parsley"
)

type CompLine struct {
	Html           string
	Number         int
	StartingSpaces int
	Tags           []*CompTag
}

func NewCompLine(str string, index int) (*CompLine, error) {
	line := &CompLine{
		Html:           str,
		Number:         index,
		StartingSpaces: parsley.CountStartingSpaces(str),
	}
	err := line.makeCompTags()
	if err != nil {
		return nil, err
	}
	return line, nil
}

func (line *CompLine) makeCompTags() error {
	inHtml := false
	currentTag := ""
	var tags []string
	err := parsley.MapChars(line.Html, func(i int, ch string) error {
		if ch == "<" {
			inHtml = true
		}
		if inHtml {
			currentTag += ch
		}
		if ch == ">" {
			inHtml = false
			tags = append(tags, currentTag)
			currentTag = ""

		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, tag := range tags {
		compTag, err := NewCompTag(tag)
		if err != nil {
			return err
		}
		line.Tags = append(line.Tags, compTag)
	}
	return nil
}
