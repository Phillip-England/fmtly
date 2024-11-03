package tag

import (
	"github.com/PuerkitoBio/goquery"
)

type ElseTag struct {
	Info *TagInfo
}

func NewElseTagFromSelection(s *goquery.Selection) (*ElseTag, error) {
	info, err := NewTagInfoFromSelection(s)
	if err != nil {
		return nil, err
	}
	t := &ElseTag{
		Info: info,
	}
	return t, nil
}
