package gtml

import (
	"fmt"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyVarGoFor = "VARGOFOR"
	KeyVarGoIf  = "VARGOIF"
)

// ##==================================================================
type Var interface {
	GetData() string
	GetVarName() string
	GetType() string
}

func NewVar(elm Element) (Var, error) {
	switch elm.GetType() {
	case KeyElementFor:
		v, err := NewVarGoFor(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case KeyElementIf:
		v, err := NewVarGoIf(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, fmt.Errorf("element does not corrospond to a valid Var: %s", elm.GetHtml())
}

// ##==================================================================
type VarGoFor struct {
	Element       Element
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

func NewVarGoFor(elm Element) (*VarGoFor, error) {
	v := &VarGoFor{
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

func (v *VarGoFor) GetData() string    { return v.Data }
func (v *VarGoFor) GetVarName() string { return v.VarName }
func (v *VarGoFor) GetType() string    { return v.Type }

func (v *VarGoFor) initVarName() error {
	v.VarName = v.Element.GetAttrParts()[0]
	return nil
}

func (v *VarGoFor) initBasicInfo() error {
	attrParts := v.Element.GetAttrParts()
	v.VarName = attrParts[0] + "For"
	v.BuilderName = attrParts[0] + "Builder"
	v.IterItems = attrParts[2]
	v.IterItem = attrParts[0]
	v.IterType = purse.RemoveAllSubStr(attrParts[3], "[]")
	v.Type = KeyVarGoFor
	return nil
}

func (v *VarGoFor) initVars() error {
	vars, err := GetElementVars(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoFor) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += "\t" + inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoFor) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *VarGoFor) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlFor(%s, func(i int, %s %s) string {
	var %s strings.Builder
%s
%s
	return %s.String()
})`, v.VarName, v.IterItems, v.IterItem, v.IterType, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}

// ##==================================================================
type VarGoIf struct {
	Element       Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	BoolToCheck   string
	Type          string
}

func NewVarGoIf(elm Element) (*VarGoIf, error) {
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

func (v *VarGoIf) GetData() string    { return v.Data }
func (v *VarGoIf) GetVarName() string { return v.VarName }
func (v *VarGoIf) GetType() string    { return v.Type }

func (v *VarGoIf) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "If"
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoIf
	return nil
}

func (v *VarGoIf) initVars() error {
	vars, err := GetElementVars(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoIf) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += "\t" + inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoIf) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
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
})`, v.VarName, v.BoolToCheck, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BoolToCheck, v.BuilderName))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
