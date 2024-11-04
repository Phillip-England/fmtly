package tag

import (
	"github.com/PuerkitoBio/goquery"
)

type ElseTag struct {
	Info *TagInfo
}

func NewElseTagFromSelection(s *goquery.Selection) (*ElseTag, error) {
	info, err := NewTagInfoFromSelection(s, "else", []string{})
	if err != nil {
		return nil, err
	}
	t := &ElseTag{
		Info: info,
	}
	return t, nil
}

func (t *ElseTag) Html() string  { return t.Info.Html }
func (t *ElseTag) Name() string  { return t.Info.Name }
func (t *ElseTag) Scopes() []Tag { return t.Info.Scopes }
