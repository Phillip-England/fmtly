package funcarg

import "fmt"

type ArgRaw struct {
	Value string
	Type  string
}

func (arg *ArgRaw) GetValue() string { return arg.Value }
func (arg *ArgRaw) GetType() string  { return arg.Type }
func (arg *ArgRaw) Print()           { fmt.Println(arg.Value) }

func NewArgRaw(str string) (*ArgRaw, error) {
	arg := &ArgRaw{
		Value: str,
		Type:  KeyFuncArgRaw,
	}
	return arg, nil
}
