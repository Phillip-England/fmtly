package gtmlvar

import (
	"fmt"
	"gtml/src/parser"
	"gtml/src/parser/element"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type VarGoIf struct {
	Element       element.Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	BoolToCheck   string
	Type          string
}

func NewGoIf(elm element.Element) (*VarGoIf, error) {
	v := &VarGoIf{
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

func (v *VarGoIf) GetData() string        { return v.Data }
func (v *VarGoIf) GetVarName() string     { return v.VarName }
func (v *VarGoIf) GetBuilderName() string { return v.BuilderName }
func (v *VarGoIf) GetType() string        { return v.Type }
func (v *VarGoIf) Print()                 { fmt.Print(v.Data) }

func (v *VarGoIf) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "If" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoIf
	return nil
}

func (v *VarGoIf) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoIf) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoIf) initBuilderSeries() error {
	series, err := parser.GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *VarGoIf) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlIf(%s, func() string {
var %s strings.Builder
%s
%s
if %s {
return %s.String()
}
return ""
})`+"\n", v.VarName, v.BoolToCheck, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BoolToCheck, v.BuilderName))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
