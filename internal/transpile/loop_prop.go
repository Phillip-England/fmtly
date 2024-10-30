package transpile

import "strings"

type LoopProp struct {
	Text  string
	Value string
}

func NewLoopProp(str string) (*LoopProp, error) {
	value := strings.ReplaceAll(str, "{", "")
	value = strings.ReplaceAll(value, "}", "")
	value = strings.ReplaceAll(value, " ", "")
	prop := &LoopProp{
		Text:  str,
		Value: value,
	}
	return prop, nil
}
