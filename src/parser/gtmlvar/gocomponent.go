package gtmlvar

import (
	"fmt"
	"gtml/src/parser/element"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoComponent struct {
	Element       element.Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	Type          string
}

func NewGoComponent(elm element.Element) (*GoComponent, error) {
	v := &GoComponent{
		Element: elm,
	}
	err := fungi.Process(
		func() error { return v.initBasicInfo() },
		func() error { return v.initVars() },
		func() error { return v.initWriteVarsAs() },
		func() error { return v.initBuilderSeries() },
		func() error { return v.initData() },
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *GoComponent) GetData() string             { return v.Data }
func (v *GoComponent) GetVarName() string          { return v.VarName }
func (v *GoComponent) GetBuilderName() string      { return v.BuilderName }
func (v *GoComponent) GetType() string             { return v.Type }
func (v *GoComponent) GetElement() element.Element { return v.Element }
func (v *GoComponent) Print()                      { fmt.Print(v.Data) }

func (v *GoComponent) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = strings.ToLower(attr)
	v.BuilderName = v.VarName + "Builder"
	v.Type = KeyVarGoElse
	return nil
}

func (v *GoComponent) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoComponent) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoComponent) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *GoComponent) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := func() string {
var %s strings.Builder
%s
%s
return %s.String()
}`+"\n", v.VarName, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
	// v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
