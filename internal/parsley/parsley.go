package parsley

import "strings"

func CountLeadingSpaces(s string) int {
	count := 0
	for _, char := range s {
		if char == ' ' {
			count++
		} else {
			break
		}
	}
	return count
}

func MakeLines(s string) []string {
	return strings.Split(s, "\n")
}

func JoinLines(slc []string) string {
	return strings.Join(slc, "\n")
}

func Squeeze(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func TrimLeadingSpaces(s string) string {
	leadingSpaces := CountLeadingSpaces(s)
	return s[leadingSpaces:]
}

func TrimLineNumber(lines []string, lineNumber int) []string {
	if lineNumber < 0 || lineNumber >= len(lines) {
		return lines
	}
	return append(lines[:lineNumber], lines[lineNumber+1:]...)
}

func TrimFirstLine(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	return lines[1:]
}

func TrimLastLine(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	return lines[:len(lines)-1]
}

func ReplaceFirstLine(input, newLine string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		lines[0] = newLine
	}
	return strings.Join(lines, "\n")
}

func GetFirstLine(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		return lines[0]
	}
	return ""
}

func GetLastLine(input string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		return lines[len(lines)-1]
	}
	return ""
}

func EqualsOneof(value string, options ...string) bool {
	for _, option := range options {
		if strings.EqualFold(value, option) { // Case-insensitive comparison
			return true
		}
	}
	return false
}
