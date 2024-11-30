package call

import (
	"fmt"
	"strings"
)

type Call interface {
	GetData() string
	Print()
	GetParams() []string
}

func NewCall(str string) (Call, error) {
	if strings.Contains(str, "ATTRID") {
		call, err := NewPlaceholder(str)
		if err != nil {
			return nil, err
		}
		return call, nil
	}
	return nil, fmt.Errorf("func call as string cannot be converted into type 'Call': %s", str)
}
