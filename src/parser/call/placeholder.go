package call

import (
	"fmt"
	"gtml/src/parser/gtmlrune"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type Placeholder struct {
	Data   string
	Params []string
}

func NewPlaceholder(str string) (*Placeholder, error) {
	call := &Placeholder{
		Data: str,
	}
	err := fungi.Process(
		func() error { return call.initParams() },
		func() error { return call.initRunes() },
	)
	if err != nil {
		return nil, err
	}
	return call, nil
}

func (call *Placeholder) GetData() string     { return call.Data }
func (call *Placeholder) Print()              { fmt.Println(call.Data) }
func (call *Placeholder) GetParams() []string { return call.Params }

func (call *Placeholder) initParams() error {
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

func (call *Placeholder) initRunes() error {
	for i, param := range call.Params {
		rns, err := gtmlrune.NewRunesFromStr(param)
		if err != nil {
			return err
		}
		for _, rn := range rns {
			if rn.GetType() == gtmlrune.KeyRuneProp || rn.GetType() == gtmlrune.KeyRunePipe {
				runeVal := rn.GetValue()
				param = strings.Replace(param, rn.GetDecodedData(), runeVal, 1)
				param = purse.RemoveAllSubStr(param, "'", "\"")
				call.Params[i] = param
			} else {
				purse.Fmt(`
_placeholder component found with the following rune in its attributes: %s
only $prop and $pipe runes are usable within _placeholder component attributes`, rn.GetDecodedData())
			}
		}
	}
	return nil
}
