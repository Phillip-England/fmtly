package param

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlrune"
	"strings"

	"github.com/phillip-england/purse"
)

type Param interface {
	GetStr() string
	GetName() string
	GetType() string
	Print()
}

func NewParam(name string, typeof string) (Param, error) {
	param := NewParamGoFunc(name, typeof)
	return param, nil
}

func NewParamsFromElement(elm element.Element) ([]Param, error) {
	params := make([]Param, 0)
	// pulling params from runes from the root and its elements
	err := element.WalkElementChildrenIncludingRoot(elm, func(child element.Element) error {
		runes, err := gtmlrune.NewRunesFromElement(child)
		if err != nil {
			return err
		}
		for _, rn := range runes {
			if rn.GetType() == gtmlrune.KeyRuneProp || rn.GetType() == gtmlrune.KeyRuneSlot {
				param, err := NewParam(rn.GetValue(), "string")
				if err != nil {
					return err
				}
				params = append(params, param)
			}
		}
		return nil
	})
	if err != nil {
		return params, err
	}
	// pulling element specific params
	elementSpecificParams := make([]Param, 0)
	err = element.WalkElementChildrenIncludingRoot(elm, func(child element.Element) error {
		params := make([]Param, 0)
		elmType := child.GetType()
		if elmType == element.KeyElementComponent {
			return nil
		}

		if elmType == element.KeyElementPlaceholder {
			return nil
		}

		if elmType == element.KeyElementSlot {
			return nil
		}
		if elmType == element.KeyElementElse {
			param, err := NewParam(child.GetAttr(), "bool")
			if err != nil {
				return err
			}
			params = append(params, param)
		}
		if elmType == element.KeyElementFor {
			parts := child.GetAttrParts()
			iterItems := parts[2]
			if strings.Contains(iterItems, ".") {
				return nil
			}
			iterType := parts[3]
			param, err := NewParam(iterItems, iterType)
			if err != nil {
				return err
			}
			params = append(params, param)
		}
		if elmType == element.KeyElementIf {
			param, err := NewParam(child.GetAttr(), "bool")
			if err != nil {
				return err
			}
			params = append(params, param)
		}

		elementSpecificParams = append(elementSpecificParams, params...)
		return nil
	})
	if err != nil {
		return params, err
	}
	// merging the params
	params = append(params, elementSpecificParams...)
	filtered := make([]Param, 0)
	found := make([]string, 0)
	for _, p := range params {
		if !purse.SliceContains(found, p.GetStr()) {
			found = append(found, p.GetStr())
			filtered = append(filtered, p)
			continue
		}
	}
	return filtered, nil
}

type ParamGoFunc struct {
	Name string
	Type string
}

func NewParamGoFunc(name string, typeof string) *ParamGoFunc {
	return &ParamGoFunc{
		Name: name,
		Type: typeof,
	}
}

func (param *ParamGoFunc) GetStr() string  { return param.Name + " " + param.Type }
func (param *ParamGoFunc) GetName() string { return param.Name }
func (param *ParamGoFunc) GetType() string { return param.Type }
func (param *ParamGoFunc) Print()          { fmt.Println(param.GetStr()) }
