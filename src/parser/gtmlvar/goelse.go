package gtmlvar

import (
	"fmt"
	"gtml/src/parser"
	"gtml/src/parser/element"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoElse struct {
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

func NewGoElse(elm element.Element) (*GoElse, error) {
	v := &GoElse{
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

func (v *GoElse) GetData() string        { return v.Data }
func (v *GoElse) GetVarName() string     { return v.VarName }
func (v *GoElse) GetBuilderName() string { return v.BuilderName }
func (v *GoElse) GetType() string        { return v.Type }
func (v *GoElse) Print()                 { fmt.Print(v.Data) }

func (v *GoElse) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "Else" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoElse
	return nil
}

func (v *GoElse) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoElse) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoElse) initBuilderSeries() error {
	series, err := parser.GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *GoElse) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlElse(%s, func() string {
var %s strings.Builder
%s
%s
if !%s {
	return %s.String()
}
return ""
})`+"\n", v.VarName, v.BoolToCheck, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BoolToCheck, v.BuilderName))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
