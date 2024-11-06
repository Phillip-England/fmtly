package tag

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FmtTag struct {
	TagInfo  *TagInfo
	GoOutput *GoOutput
}

func NewFmtTagFromSelection(selection *goquery.Selection) (*FmtTag, error) {
	info, err := NewTagInfoFromSelection(selection, []string{"name", "tag"})
	if err != nil {
		return nil, err
	}
	t := &FmtTag{
		TagInfo: info,
	}
	out, err := NewOutputFromTag(t)
	if err != nil {
		return nil, err
	}
	t.GoOutput = out
	return t, nil
}

func NewFmtTagsFromDir(dir string) ([]Tag, error) {
	var fmtTags []Tag
	err := filepath.Walk("./components", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fStr := string(f)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
		if err != nil {
			return err
		}
		var potErr error
		potErr = nil
		doc.Find("fmt").Each(func(i int, s *goquery.Selection) {
			ft, err := NewTagFromSelection(s)
			if err != nil {
				potErr = err
				return
			}
			innerFmt := s.Find("fmt")
			if innerFmt.Length() > 0 {
				return
			}
			fmtTags = append(fmtTags, ft)
		})
		if potErr != nil {
			return potErr
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fmtTags, nil
}

func (t *FmtTag) Info() *TagInfo { return t.TagInfo }
func (t *FmtTag) Out() string    { return t.GoOutput.Html }

func (t *FmtTag) MakeGoOutput() (string, error) {
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
	finalOut := fmt.Sprintf("func NAME(PARAMS) string {\n\treturn `%s`\n}", out)
	return finalOut, nil
}
