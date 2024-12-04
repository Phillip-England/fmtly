package funcarg

import (
	"fmt"
	"strings"
)

type ArgStr struct {
	Value string
	Type  string
}

func (arg *ArgStr) GetValue() string { return arg.Value }
func (arg *ArgStr) GetType() string  { return arg.Type }
func (arg *ArgStr) Print()           { fmt.Println(arg.Value) }

func NewArgStr(str string) (*ArgStr, error) {
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.ReplaceAll(str, "'", "")
	arg := &ArgStr{
		Value: str,
		Type:  KeyFuncArgRaw,
	}
	return arg, nil
}
