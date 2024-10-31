package transpile

import (
	"fmtly/internal/parsley"
	"strings"
)

type ForTag struct {
	Text             string
	InAttr           string
	TypeAttr         string
	TagAttr          string
	IsCloseTag       bool
	LineNumber       int
	Html             string
	Output           string
	EndTagLineNumber int
	IsInnerFor       bool
	IsRootFor        bool
}

func NewForTag(line *CompLine, tag *CompTag, comp *Comp) (*ForTag, error) {
	inAttr, _ := tag.GetAttr("in")
	typeAttr, _ := tag.GetAttr("type")
	tagAttr, _ := tag.GetAttr("tag")

	forTag := &ForTag{
		Text:       tag.Html,
		InAttr:     inAttr,
		TagAttr:    tagAttr,
		TypeAttr:   typeAttr,
		IsCloseTag: tag.IsCloseTag,
		LineNumber: line.Number,
	}

	err := forTag.setHtml(comp)
	if err != nil {
		return nil, err
	}

	err = forTag.setEndTagLineNumber()
	if err != nil {
		return nil, err
	}

	err = forTag.writeTagNames()
	if err != nil {
		return nil, err
	}

	err = forTag.setIsInnerFor()
	if err != nil {
		return nil, err
	}

	err = forTag.setCloseFors()
	if err != nil {
		return nil, err
	}

	err = forTag.setRootFors()
	if err != nil {
		return nil, err
	}

	return forTag, nil
}

func (tag *ForTag) setHtml(comp *Comp) error {
	if tag.IsCloseTag {
		return nil
	}

	isFirstLoop := true
	forCount := 0
	var col []string
	for i, ln := range comp.Lines {
		if i >= tag.LineNumber {
			if isFirstLoop {
				col = append(col, ln.Html)
				isFirstLoop = false
			} else {
				if forCount == -1 {
					break
				}
				sq := parsley.Squeeze(ln.Html)
				if strings.HasPrefix(sq, "<for") && !strings.HasPrefix(sq, "<form") {
					forCount++
				}
				if strings.HasPrefix(sq, "</for>") {
					forCount--
				}
				col = append(col, ln.Html)
			}
		}
	}
	str := strings.Join(col, "\n")
	tag.Html = str
	tag.Output = str
	return nil
}

func (tag *ForTag) setEndTagLineNumber() error {
	lineCount := parsley.CountLine(tag.Html)
	tag.EndTagLineNumber = lineCount + tag.LineNumber - 1
	return nil
}

func (tag *ForTag) writeTagNames() error {
	if tag.IsCloseTag {
		return nil
	}
	var out []string
	lines := parsley.MakeLines(tag.Output)
	for i, line := range lines {
		if i == len(lines)-1 || i == 0 {
			out = append(out, strings.Replace(line, "for", tag.TagAttr, 1))
			continue
		}
		out = append(out, line)
	}
	str := strings.Join(out, "\n")
	tag.Output = str
	return nil
}

func (tag *ForTag) setIsInnerFor() error {
	if strings.Contains(tag.InAttr, ".") {
		tag.IsInnerFor = true
	}
	return nil
}

func (tag *ForTag) setCloseFors() error {
	if strings.Contains(parsley.Squeeze(tag.Text), "</for>") {
		tag.IsCloseTag = true
	}
	return nil
}

func (tag *ForTag) setRootFors() error {
	if tag.IsCloseTag == false && tag.IsInnerFor == false {
		tag.IsRootFor = true
	}
	return nil
}
