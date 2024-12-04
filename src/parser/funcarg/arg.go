package funcarg

import (
	"github.com/phillip-england/purse"
)

type FuncArg interface {
	GetValue() string
	GetType() string
	Print()
}

func NewFuncArg(str string) (FuncArg, error) {
	str = purse.TrimLeadingSpaces(str)
	if len(str) == 0 {
		return nil, purse.Err(`
empty string provided to NewFuncArg, the string must contain a value`)
	}
	firstChar := string(str[0])
	lastChar := string(str[len(str)-1])
	if (firstChar == `"` && lastChar == `"`) || (firstChar == `'` && lastChar == `'`) {
		arg, err := NewArgStr(str)
		if err != nil {
			return nil, err
		}
		return arg, nil
	}
	arg, err := NewArgRaw(str)
	if err != nil {
		return nil, err
	}
	return arg, nil
}
