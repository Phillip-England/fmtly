package main

import (
	"fmt"
	"fmtly/internal/tag"
	"strings"
)

func main() {

	fmtTags, err := tag.NewFmtTagsFromDir("./components.go")
	if err != nil {
		panic(err)
	}

	for _, t := range fmtTags {
		for _, sc := range t.Scopes() {
			fmt.Println(sc.Name())
		}
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
