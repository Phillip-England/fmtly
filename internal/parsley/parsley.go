package parsley

import (
	"strings"
)

func Squeeze(str string) string {
	return strings.ReplaceAll(str, " ", "")
}

func MapChars(str string, fn func(i int, ch string) error) error {
	for index, char := range str {
		strCh := string(char)
		err := fn(index, strCh)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeLines(str string) []string {
	return strings.Split(str, "\n")
}

func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func CountStartingSpaces(str string) int {
	for i, ch := range str {
		char := string(ch)
		if char == " " {
			continue
		}
		return i
	}
	return 0
}

func SplitSpaces(str string) []string {
	return strings.Split(str, " ")
}

func CountLine(str string) int {
	lines := strings.Split(str, "\n")
	count := 0
	for i := 0; i < len(lines); i++ {
		count++
	}
	return count
}

func PrefiexLinesWith(str string, prefixWith string) string {
	var col []string
	for _, line := range MakeLines(str) {
		col = append(col, prefixWith+line)
	}
	return JoinLines(col)
}

func MapAny[T any](slice []T, mapFunc func(T) T) []T {
	mapped := make([]T, len(slice))
	for i, v := range slice {
		mapped[i] = mapFunc(v)
	}
	return mapped
}

func ForCollect[T any](slice []T, mapFunc func(T) string) string {
	var result string
	for _, v := range slice {
		result += mapFunc(v)
	}
	return result
}

func ReverseSlice[T any](s []T) []T {
	reversed := make([]T, len(s))
	for i, v := range s {
		reversed[len(s)-1-i] = v
	}
	return reversed
}

func InsertLinesAt(lines []string, insert []string, insertAt int) []string {
	var out []string
	for i, line := range lines {
		if i == insertAt {
			for _, line2 := range insert {
				out = append(out, line2)
			}
		}
		out = append(out, line)
	}
	return out
}

func MapLines(str string, fn func(i int, line string) error) error {
	for index, lineStr := range MakeLines(str) {
		err := fn(index, lineStr)
		if err != nil {
			return err
		}
	}
	return nil
}

func TrimStartingSpaces(s string) string {
	return strings.TrimLeft(s, " ")
}
