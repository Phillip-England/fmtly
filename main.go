package main

import (
	"fmtly/internal/tags"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	var fmtTags []*tags.Fmt
	filepath.Walk("./components", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		ft, err := tags.NewFmtFromStr(string(f))
		if err != nil {
			panic(err)
		}
		fmtTags = append(fmtTags, ft)
		return nil
	})

	CollectStr(customers, func(i int, customer []*Customer) string {
		return `
			<li in="customers" as="customer" type="[]*Customer" tag="li"><p>{{ customer.Name }}</p>` +
			CollectStr(customer.Friends, func(i int, friend []*Friend) string {
				return `
						<div in="customer.Friends" tag="div" as="friend" type="[]*Friend"><p>{{ friend.Name }}</p><p>{{ friend.Age }}</p></div>`
			}) +
			`</li>`
	})

}

func CollectStr[T any](slice []T, mapper func(i int, t T) string) string {
	var builder strings.Builder
	for i, t := range slice {
		builder.WriteString(mapper(i, t))
	}
	return builder.String()
}

// names := []string{"Alice", "Bob", "Charlie"}
// result = parsley.CollectStr(names, func(i int, name string) string {
// 	return fmt.Sprintf("<li>%s</li>", name)
// })
