package transpile

import (
	"fmt"
	"strings"
)

type CompTag struct {
	Html       string
	Parts      []string
	IsCloseTag bool
	Name       string
	StrAttrs   []string
	Attrs      []*CompAttr
}

func NewCompTag(str string) (*CompTag, error) {
	tag := &CompTag{
		Html: str,
	}
	tag.makeParts()
	tag.parseParts()
	err := tag.parseAttrs()
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (tag *CompTag) makeParts() {
	raw := strings.Replace(tag.Html, "<", "", 1)
	raw = strings.Replace(raw, ">", "", 1)
	parts := strings.Split(raw, " ")
	tag.Parts = parts
}

func (tag *CompTag) parseParts() {
	for i, part := range tag.Parts {
		if i == 0 {
			if strings.Contains(part, "/") {
				tag.IsCloseTag = true
				part = strings.Replace(part, "/", "", 1)
			}
			tag.Name = part
		}
		if i > 0 {
			if strings.Contains(part, "=\"") || strings.Contains(part, "'=") {
				tag.StrAttrs = append(tag.StrAttrs, part)
			}
		}
	}
}

func (tag *CompTag) parseAttrs() error {
	for _, strAttr := range tag.StrAttrs {
		compAttr, err := NewCompAttr(strAttr)
		if err != nil {
			return err
		}
		tag.Attrs = append(tag.Attrs, compAttr)
	}
	return nil
}

func (tag *CompTag) GetAttr(name string) (string, error) {
	outputValue := ""
	for _, attr := range tag.Attrs {
		if attr.Name == name {
			outputValue = attr.Value
		}
	}
	if outputValue == "" {
		return "", fmt.Errorf("attr: %s not found in tag: %s", name, tag.Html)
	}
	return outputValue, nil
}
