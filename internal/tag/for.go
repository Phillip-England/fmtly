package tag

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	TagInfo  *TagInfo
	GoOutput *GoOutput
	ElseTag  *ElseTag
}

func NewForTagFromSelection(selection *goquery.Selection) (*ForTag, error) {
	info, err := NewTagInfoFromSelection(selection, []string{"in", "as", "tag", "type"})
	if err != nil {
		return nil, err
	}
	elseSelection := info.Selection.Find("else")
	elseTag, err := NewElseTagFromSelection(elseSelection)
	if err != nil {
		return nil, err
	}
	t := &ForTag{
		TagInfo: info,
		ElseTag: elseTag,
	}
	out, err := NewOutputFromTag(t)
	if err != nil {
		return nil, err
	}
	t.GoOutput = out
	return t, nil
}

func (t *ForTag) Info() *TagInfo { return t.TagInfo }
func (t *ForTag) Out() string    { return t.GoOutput.Html }

func (t *ForTag) MakeGoOutput() (string, error) {
	attrTag, _ := t.Info().Selection.Attr("tag")
	htmlStr, err := t.Info().Selection.Html()
	if err != nil {
		return "", err
	}
	var out string
	if t.Info().AttrStr == "" {
		out = fmt.Sprintf("<%s>%s</%s>", attrTag, htmlStr, attrTag)
	} else {
		out = fmt.Sprintf("<%s %s>%s</%s>", attrTag, t.Info().AttrStr, htmlStr, attrTag)

	}
	return out, nil
}
