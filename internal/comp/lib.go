package comp

import (
	"fmtly/internal/parsley"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

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

func GetHtmlFromSelection(s *goquery.Selection) (string, error) {
	out := ""
	htmlStr, err := goquery.OuterHtml(s)
	if err != nil {
		return out, err
	}
	lines := parsley.MakeLines(htmlStr)
	lastLine := lines[len(lines)-1]
	spaces := parsley.CountLeadingSpaces(lastLine)
	for i, line := range lines {
		if i == 0 {
			out += strings.Repeat(" ", spaces) + line + "\n"
			continue
		}
		out += line + "\n"
	}
	return out, nil
}
