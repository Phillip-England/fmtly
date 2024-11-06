package tag

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type IfTag struct {
	TagInfo  *TagInfo
	GoOutput *GoOutput
}

func NewIfTagFromSelection(s *goquery.Selection) (*IfTag, error) {
	info, err := NewTagInfoFromSelection(s, []string{"condition", "tag"})
	if err != nil {
		return nil, err
	}
	t := &IfTag{
		TagInfo: info,
	}
	out, err := NewOutputFromTag(t)
	if err != nil {
		return nil, err
	}
	t.GoOutput = out
	return t, nil
}

func (t *IfTag) Info() *TagInfo { return t.TagInfo }
func (t *IfTag) Out() string    { return t.GoOutput.Html }

func (t *IfTag) MakeGoOutput() (string, error) {

	attrTag, _ := t.Info().Selection.Attr("tag")
	attrCondition, _ := t.Info().Selection.Attr("condition")

	elseSelection := t.Info().Selection.Find("else")
	elseHtml, err := elseSelection.Html()
	if err != nil {
		return "", err
	}
	var elseOut string
	if t.Info().AttrStr == "" {
		elseOut = fmt.Sprintf("<%s>%s</%s>", attrTag, elseHtml, attrTag)
	} else {
		elseOut = fmt.Sprintf("<%s %s>%s</%s>", attrTag, t.Info().AttrStr, elseHtml, attrTag)

	}

	selectionCopy := t.Info().Selection.Clone()
	selectionCopy.Find("else").Remove()
	htmlStr, err := selectionCopy.Html()
	if err != nil {
		return "", err
	}
	var ifOut string
	if t.Info().AttrStr == "" {
		ifOut = fmt.Sprintf("<%s>%s</%s>", attrTag, htmlStr, attrTag)
	} else {
		ifOut = fmt.Sprintf("<%s %s>%s</%s>", attrTag, t.Info().AttrStr, htmlStr, attrTag)

	}

	finalOut := fmt.Sprintf("` + If(%s, `%s`, `%s`) + `", attrCondition, ifOut, elseOut)

	return finalOut, nil
}
