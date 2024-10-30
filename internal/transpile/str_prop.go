package transpile

import "strings"

type StrProp struct {
	Text  string
	Value string
}

func NewStrProp(str string) (*StrProp, error) {
	value := strings.ReplaceAll(str, "{", "")
	value = strings.ReplaceAll(value, "}", "")
	value = strings.ReplaceAll(value, " ", "")
	prop := &StrProp{
		Text:  str,
		Value: value,
	}
	return prop, nil
}
