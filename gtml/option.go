package gtml

import (
	"fmt"

	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyOptionWatch = "--watch"
)

// ##==================================================================
func getOptionList() []string {
	return []string{KeyOptionWatch}
}

// ##==================================================================
type Option interface {
	Print()
	GetType() string
}

func NewOption(arg string) (Option, error) {
	match := purse.FindMatchInStrSlice(getOptionList(), arg)
	if match == "" {
		return nil, fmt.Errorf("invalid option selected: %s\nRun 'gtml help' for usage.", arg)
	}
	switch match {
	case KeyOptionWatch:
		opt, err := NewOptionWatch()
		if err != nil {
			return nil, err
		}
		return opt, err
	}
	return nil, fmt.Errorf("invalid option selected: %s\nRun 'gtml help' for usage.", arg)
}

// ##==================================================================
type OptionWatch struct {
	Type string
}

func NewOptionWatch() (*OptionWatch, error) {
	opt := &OptionWatch{
		Type: KeyOptionWatch,
	}
	return opt, nil
}

func (opt *OptionWatch) GetType() string { return opt.Type }
func (opt *OptionWatch) Print()          { fmt.Println(opt.Type) }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
