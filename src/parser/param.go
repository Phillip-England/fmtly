package parser

import "fmt"

// ##==================================================================
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

// ##==================================================================
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

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
