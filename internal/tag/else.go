package tag

import (
	"github.com/PuerkitoBio/goquery"
)

type ElseTag struct {
	TagInfo  *TagInfo
	GoOutput *GoOutput
}

func NewElseTagFromSelection(s *goquery.Selection) (*ElseTag, error) {
	info, err := NewTagInfoFromSelection(s, []string{})
	if err != nil {
		return nil, err
	}
	t := &ElseTag{
		TagInfo: info,
	}
	out, err := NewOutputFromTag(t)
	if err != nil {
		return nil, err
	}
	t.GoOutput = out
	return t, nil
}

func (t *ElseTag) Info() *TagInfo                { return t.TagInfo }
func (t *ElseTag) MakeGoOutput() (string, error) { return "", nil }
func (t *ElseTag) Out() string                   { return t.GoOutput.Html }
