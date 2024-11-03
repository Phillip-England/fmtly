package main

import (
	"fmtly/internal/tag"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var fmtTags []*tag.Fmt
	filepath.Walk("./components", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		ft, err := tag.NewFmtTagFromStr(string(f))
		if err != nil {
			panic(err)
		}
		fmtTags = append(fmtTags, ft)
		return nil
	})

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
