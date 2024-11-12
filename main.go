package main

import (
	"gotml/internal/gotml"
	"strings"
)

func main() {

	_, err := gotml.NewHtmlFileFromPath("./components/index.html")
	if err != nil {
		panic(err)
	}

}

func collect[T any](items []T, callback func(i int, item T) string) string {
	var builder strings.Builder
	for i, item := range items {
		builder.WriteString(callback(i, item))
	}
	return builder.String()
}
