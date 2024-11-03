package main

import (
	"fmtly/internal/tag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	var fmtTags []*tag.FmtTag
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
			ft, err := tag.NewFmtTagFromSelection(s)
			if err != nil {
				potErr = err
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
		panic(err)
	}

}

func collectStr[T any](slice []T, mapper func(i int, t T) string) string {
	var builder strings.Builder
	for i, t := range slice {
		builder.WriteString(mapper(i, t))
	}
	return builder.String()
}

func ifElse(cond bool, ifTrue string, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// names := []string{"Alice", "Bob", "Charlie"}
// result = parsley.CollectStr(names, func(i int, name string) string {
// 	return fmt.Sprintf("<li>%s</li>", name)
// })
