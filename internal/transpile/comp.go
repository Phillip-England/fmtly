package transpile

import (
	"fmtly/internal/parsley"
	"strings"
)

type Comp struct {
	Html      string
	Output    string
	Lines     []*CompLine
	LoopProps []*LoopProp
	StrProps  []*StrProp
	FmtTag    *FmtTag
	ForTags   []*ForTag
	IfTags    []*IfTag
	ElseTags  []*ElseTag
}

func NewFmtComp(compStr string) (*Comp, error) {
	comp := &Comp{
		Html:   compStr,
		Output: compStr,
	}
	err := comp.makeLines()
	if err != nil {
		return nil, err
	}
	err = comp.makeProps()
	if err != nil {
		return nil, err
	}
	err = comp.setFmtTag()
	if err != nil {
		return nil, err
	}
	err = comp.setForTags()
	if err != nil {
		return nil, err
	}
	err = comp.setIfTags()
	if err != nil {
		return nil, err
	}
	err = comp.setElseTags()
	if err != nil {
		return nil, err
	}
	return comp, nil
}

func (comp *Comp) makeLines() error {
	for i, line := range parsley.MakeLines(comp.Html) {
		compLine, err := NewCompLine(line, i)
		if err != nil {
			return err
		}
		comp.Lines = append(comp.Lines, compLine)
	}
	return nil
}

func (comp *Comp) makeProps() error {
	inProp := false
	propStr := ""
	var slice []string
	err := parsley.MapChars(comp.Html, func(i int, ch string) error {
		if i+2 > len(comp.Html)-1 {
			return nil
		}
		ck := comp.Html[i : i+2]
		if ck == "{{" {
			inProp = true
		}
		if inProp {
			propStr += ch
		}
		if ck == "}}" {
			inProp = false
			propStr += string(comp.Html[i+1])
			slice = append(slice, propStr)
			propStr = ""
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, str := range slice {
		// loop props will always contain a "."
		if strings.Count(str, ".") == 1 {
			loopProp, err := NewLoopProp(str)
			if err != nil {
				return err
			}
			comp.LoopProps = append(comp.LoopProps, loopProp)
			continue
		}
		// str props will not contains a "."
		strProp, err := NewStrProp(str)
		if err != nil {
			return err
		}
		comp.StrProps = append(comp.StrProps, strProp)
	}
	return nil
}

func (comp *Comp) setFmtTag() error {
	for _, line := range comp.Lines {
		for _, tag := range line.Tags {
			if tag.Name == "fmt" {
				fmtTag, err := NewFmtTag(line, tag, comp)
				if err != nil {
					return err
				}
				comp.FmtTag = fmtTag
			}
		}
	}
	return nil
}

func (comp *Comp) setForTags() error {
	for _, line := range comp.Lines {
		for _, tag := range line.Tags {
			if tag.Name == "for" {
				forTag, err := NewForTag(line, tag, comp)
				if err != nil {
					return err
				}
				comp.ForTags = append(comp.ForTags, forTag)
			}
		}
	}
	return nil
}

func (comp *Comp) setIfTags() error {
	for _, line := range comp.Lines {
		for _, tag := range line.Tags {
			if tag.Name == "if" {
				ifTag, err := NewIfTag(line, tag, comp)
				if err != nil {
					return err
				}
				comp.IfTags = append(comp.IfTags, ifTag)
			}
		}
	}
	return nil
}

func (comp *Comp) setElseTags() error {
	for _, line := range comp.Lines {
		for _, tag := range line.Tags {
			if tag.Name == "else" {
				elseTag, err := NewElseTag(line, tag, comp)
				if err != nil {
					return err
				}
				comp.ElseTags = append(comp.ElseTags, elseTag)
			}
		}
	}
	return nil
}
