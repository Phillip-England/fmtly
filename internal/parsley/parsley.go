package parsley

import (
	"fmt"
	"os"
	"strings"
)

func GetTick() string {
	return "`"
}

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

func TrimLineNumber(text string, lineNumber int) string {
	lines := strings.Split(text, "\n")
	if lineNumber < 0 || lineNumber >= len(lines) {
		return text // Return the original text if the line number is out of range
	}
	lines = append(lines[:lineNumber], lines[lineNumber+1:]...)
	return strings.Join(lines, "\n")
}

func TrimFirstLine(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return "" // Return empty string if there's only one or no lines
	}
	return strings.Join(lines[1:], "\n")
}

func TrimLastLine(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) <= 1 {
		return "" // Return empty string if there's only one or no lines
	}
	return strings.Join(lines[:len(lines)-1], "\n")
}

func ReplaceFirstLine(input, newLine string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		lines[0] = newLine
	}
	return strings.Join(lines, "\n")
}

func ReplaceLastLine(input, newLine string) string {
	lines := strings.Split(input, "\n")
	if len(lines) > 0 {
		lines[len(lines)-1] = newLine
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

func LastIndexOf(s string, ch rune) int {
	for i := len(s) - 1; i >= 0; i-- {
		if rune(s[i]) == ch {
			return i
		}
	}
	return -1
}

func TabLinesBy(text string, numTabs int) string {
	lines := strings.Split(text, "\n")
	tabPrefix := strings.Repeat("\t", numTabs)

	for i, line := range lines {
		lines[i] = tabPrefix + line
	}

	return strings.Join(lines, "\n")
}

func FilterChars(input string, filters ...string) string {
	var result strings.Builder

	for _, char := range input {
		charStr := string(char)
		for _, filter := range filters {
			if charStr == filter {
				result.WriteString(charStr)
				break
			}
		}
	}

	return result.String()
}

func TrimOuterEmptyLines(text string) string {
	lines := MakeLines(text)

	// Find the index of the first non-empty line
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}

	// Find the index of the last non-empty line
	end := len(lines) - 1
	for end >= 0 && strings.TrimSpace(lines[end]) == "" {
		end--
	}

	// If all lines are empty, return an empty string
	if start > end {
		return ""
	}

	// Return the lines between the first and last non-empty line
	return JoinLines(lines[start : end+1])
}

func RemoveEmptyLines(input string) string {
	var result []string
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func FlattenLines(lines []string) []string {
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " \t")
	}
	return lines
}

func FlattenStr(str string) string {
	lines := MakeLines(str)
	flat := FlattenLines(lines)
	return strings.Join(flat, "")
}

func Log(message string, path string) error {
	// Open the file in write-only mode, creating/truncating it each time
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Write the message to the file
	_, err = file.WriteString(message + "\n")
	if err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

func CountLines(s string) int {
	if s == "" {
		return 0
	}
	return len(strings.Split(s, "\n"))
}
