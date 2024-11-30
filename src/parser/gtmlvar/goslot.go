package gtmlvar

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlrune"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoSlot struct {
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

func NewGoSlot(elm element.Element) (*GoSlot, error) {
	v := &GoSlot{
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

func (v *GoSlot) GetData() string             { return v.Data }
func (v *GoSlot) GetVarName() string          { return v.VarName }
func (v *GoSlot) GetBuilderName() string      { return v.BuilderName }
func (v *GoSlot) GetType() string             { return v.Type }
func (v *GoSlot) GetElement() element.Element { return v.Element }
func (v *GoSlot) Print()                      { fmt.Print(v.Data) }

func (v *GoSlot) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "Slot" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoSlot
	return nil
}

func (v *GoSlot) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoSlot) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoSlot) initBuilderSeries() error {
	varCalls, err := GetVarsAsWriteStringCalls(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	runeCalls, err := gtmlrune.GetRunesAsWriteStringCalls(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	allCalls := make([]string, 0)
	allCalls = append(allCalls, varCalls...)
	allCalls = append(allCalls, runeCalls...)
	if len(allCalls) == 0 {
		singleCall := fmt.Sprintf("%s.WriteString(`%s`)", v.BuilderName, v.Element.GetHtml())
		allCalls = append(allCalls, singleCall)
	}
	// series, err := parser.GetElementAsBuilderSeries(v.Element, allCalls, v.BuilderName)
	// if err != nil {
	// 	return err
	// }
	// v.BuilderSeries = series
	return nil
}

func (v *GoSlot) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlSlot(func() string {
var %s strings.Builder
%s
%s
return %s.String()
})`+"\n", v.VarName, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
	// v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
