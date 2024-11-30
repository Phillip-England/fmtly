package gtmlvar

import (
	"fmt"
	"gtml/src/parser/attr"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlrune"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoPlaceholder struct {
	Element       element.Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	Type          string
	ComponentName string
	Attrs         []attr.Attr
	ParamStr      string
	CallParams    []string
	CallParamStr  string
}

func NewGoPlaceholder(elm element.Element) (*GoPlaceholder, error) {
	v := &GoPlaceholder{
		Element: elm,
	}
	err := fungi.Process(
		func() error { return v.initBasicInfo() },
		func() error { return v.initAttrs() },
		func() error { return v.initVars() },
		func() error { return v.initWriteVarsAs() },
		func() error { return v.initBuilderSeries() },
		func() error { return v.initCallParams() },
		func() error { return v.initData() },
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *GoPlaceholder) GetData() string             { return v.Data }
func (v *GoPlaceholder) GetVarName() string          { return v.VarName }
func (v *GoPlaceholder) GetBuilderName() string      { return v.BuilderName }
func (v *GoPlaceholder) GetType() string             { return v.Type }
func (v *GoPlaceholder) GetElement() element.Element { return v.Element }
func (v *GoPlaceholder) Print()                      { fmt.Print(v.Data) }

func (v *GoPlaceholder) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = strings.ToLower(attr) + "Placeholder" + v.Element.GetId()
	v.BuilderName = strings.ToLower(attr) + "Builder"
	v.ComponentName = attr
	v.Type = KeyVarGoPlaceholder
	return nil
}

func (v *GoPlaceholder) initAttrs() error {
	for _, a := range v.Element.GetSelection().Get(0).Attr {
		if strings.HasPrefix(a.Key, "_") {
			continue
		}
		attr, err := attr.NewAttr(a.Key, a.Val)
		if err != nil {
			return err
		}
		v.Attrs = append(v.Attrs, attr)
	}
	return nil
}

func (v *GoPlaceholder) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoPlaceholder) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoPlaceholder) initBuilderSeries() error {
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

func (v *GoPlaceholder) initCallParams() error {
	for _, attr := range v.Attrs {
		v.CallParams = append(v.CallParams, "ATTRID"+attr.GetKey()+"ATTRID\""+attr.GetValue()+"\"")
	}
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	for _, inner := range vars {
		if inner.GetType() == KeyVarGoSlot {
			varName := inner.GetVarName()
			i := strings.Index(varName, "Slot")
			if i == -1 {
				return fmt.Errorf("_slot element found with a VarName which doesn't end in 'Slot' you need to check NewVarGoSlot")
			}
			firstPart := varName[:i]
			// secondPart := varName[i:]
			firstPart = "ATTRID" + firstPart + "ATTRID"
			v.CallParams = append(v.CallParams, firstPart+varName)
			continue
		}
		v.CallParams = append(v.CallParams, inner.GetVarName())
	}
	v.CallParamStr = strings.Join(v.CallParams, ", ")
	return nil
}

func (v *GoPlaceholder) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := func() string {
%s
return %s(%s)
}`+"\n", v.VarName, v.WriteVarsAs, v.ComponentName, v.CallParamStr))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
