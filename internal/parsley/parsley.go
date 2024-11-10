package parsley

import (
	"strings"
)

func MakeLines(s string) []string {
	return strings.Split(s, "\n")
}

func JoinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

func ReplaceLastSubStr(s, old, new string) string {
	pos := strings.LastIndex(s, old)
	if pos == -1 {
		return s
	}
	return s[:pos] + new + s[pos+len(old):]
}

func GetFirstLine(s string) string {
	lines := MakeLines(s)
	if len(lines) == 0 {
		return s
	}
	return lines[0]
}

func GetLastLine(s string) string {
	lines := MakeLines(s)
	if len(lines) == 0 {
		return s
	}
	return lines[len(lines)-1]
}

func RemoveAllSubStr(s string, subs ...string) string {
	for _, sub := range subs {
		s = strings.ReplaceAll(s, sub, "")
	}
	return s
}

func CountLeadingSpaces(line string) int {
	count := 0
	for _, char := range line {
		if char != ' ' {
			break
		}
		count++
	}
	return count
}

func PrefixLines(str, prefix string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return strings.Join(lines, "\n")
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

func TrimLeadingSpaces(str string) string {
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimLeft(line, " ")
	}
	return strings.Join(lines, "\n")
}

func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func BackTick() string {
	return "`"
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

func Squeeze(s string) string {
	return strings.ReplaceAll(s, " ", "")
}

func ScanBetweenSubStrs(s, start, end string) []string {
	var out []string
	inSearch := false
	searchStr := ""

	i := 0
	for i < len(s) {
		// Check for the start delimiter
		if !inSearch && i+len(start) <= len(s) && s[i:i+len(start)] == start {
			inSearch = true
			searchStr = start // Start capturing, including the start delimiter
			i += len(start)
			continue
		}

		// If we're in search mode, start capturing until we find the end delimiter
		if inSearch {
			// Check for the end delimiter
			if i+len(end) <= len(s) && s[i:i+len(end)] == end {
				searchStr += end // Include the end delimiter
				out = append(out, searchStr)
				searchStr = ""
				inSearch = false
				i += len(end)
				continue
			}
			// Append current character to searchStr if still inside the delimiters
			searchStr += string(s[i])
		}

		i++
	}

	return out
}

func RemoveFirstLine(input string) string {
	index := strings.Index(input, "\n")
	if index == -1 {
		return ""
	}
	return input[index+1:]
}

func RemoveTrailingEmptyLines(input string) string {
	// Split the input string into lines
	lines := strings.Split(input, "\n")

	// Remove empty lines from the end
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}

	// Join the cleaned lines back into a single string
	return strings.Join(lines, "\n")
}

func RemoveEmptyLines(input string) string {
	// Split the input string into lines
	lines := strings.Split(input, "\n")

	// Create a slice to hold the non-empty lines
	var cleanedLines []string

	// Loop through the lines and add non-empty lines to cleanedLines
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			cleanedLines = append(cleanedLines, line)
		}
	}

	// Join the non-empty lines back into a single string
	return strings.Join(cleanedLines, "\n")
}
