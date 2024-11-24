package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
type Call interface {
	GetData() string
	Print()
	GetParams() []string
}

func NewCall(str string) (Call, error) {
	if strings.Contains(str, "ATTRID") {
		call, err := NewCallPlaceholder(str)
		if err != nil {
			return nil, err
		}
		return call, nil
	}
	return nil, fmt.Errorf("func call as string cannot be converted into type 'Call': %s", str)
}

// ##==================================================================
type CallPlaceholder struct {
	Data   string
	Params []string
}

func NewCallPlaceholder(str string) (*CallPlaceholder, error) {
	call := &CallPlaceholder{
		Data: str,
	}
	err := fungi.Process(
		func() error { return call.initParams() },
	)
	if err != nil {
		return nil, err
	}
	return call, nil
}

func (call *CallPlaceholder) GetData() string     { return call.Data }
func (call *CallPlaceholder) Print()              { fmt.Println(call.Data) }
func (call *CallPlaceholder) GetParams() []string { return call.Params }

func (call *CallPlaceholder) initParams() error {
	data := call.Data
	i := strings.Index(data, "(") + 1
	data = data[i:]
	data = purse.ReplaceLastInstanceOf(data, ")", "")
	if strings.HasSuffix(data, "\"") {
		i := strings.Index(data, "\"")
		firstHalf := data[:i]
		firstHalf = purse.Squeeze(firstHalf)
		secondHalf := data[i:]
		data = firstHalf + secondHalf
	} else {
		data = purse.Squeeze(data)
	}
	parts := strings.Split(data, ",")
	call.Params = parts
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
