package gtmlvar

import (
	"fmt"
	"gtml/src/parser"
	"gtml/src/parser/element"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoFor struct {
	Element       element.Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	IterItems     string
	IterItem      string
	IterType      string
	BuilderSeries string
	Type          string
}

func NewGoFor(elm element.Element) (*GoFor, error) {
	v := &GoFor{
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

func (v *GoFor) GetData() string        { return v.Data }
func (v *GoFor) GetVarName() string     { return v.VarName }
func (v *GoFor) GetBuilderName() string { return v.BuilderName }
func (v *GoFor) GetType() string        { return v.Type }
func (v *GoFor) Print()                 { fmt.Print(v.Data) }

func (v *GoFor) initVarName() error {
	v.VarName = v.Element.GetAttrParts()[0]
	return nil
}

func (v *GoFor) initBasicInfo() error {
	attrParts := v.Element.GetAttrParts()
	v.VarName = attrParts[0] + "For" + v.Element.GetId()
	v.BuilderName = attrParts[0] + "Builder"
	v.IterItems = attrParts[2]
	v.IterItem = attrParts[0]
	v.IterType = purse.RemoveAllSubStr(attrParts[3], "[]")
	v.Type = KeyVarGoFor
	return nil
}

func (v *GoFor) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoFor) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoFor) initBuilderSeries() error {
	series, err := parser.GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *GoFor) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlFor(%s, func(i int, %s %s) string {
var %s strings.Builder
%s
%s
return %s.String()
})`+"\n", v.VarName, v.IterItems, v.IterItem, v.IterType, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
