package parsley

import "strings"

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
