package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyVarGoFor         = "VARGOFOR"
	KeyVarGoIf          = "VARGOIF"
	KeyVarGoElse        = "VARGOELSE"
	KeyVarGoPlaceholder = "VARGOPLACEHOLDER"
	KeyVarGoSlot        = "VARGOSLOT"
)

// ##==================================================================
type Var interface {
	GetData() string
	GetVarName() string
	GetType() string
	Print()
	GetBuilderName() string
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
	case KeyElementElse:
		v, err := NewVarGoElse(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case KeyElementPlaceholder:
		v, err := NewVarGoPlaceholder(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case KeyElementSlot:
		v, err := NewVarGoSlot(elm)
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

func (v *VarGoFor) GetData() string        { return v.Data }
func (v *VarGoFor) GetVarName() string     { return v.VarName }
func (v *VarGoFor) GetBuilderName() string { return v.BuilderName }
func (v *VarGoFor) GetType() string        { return v.Type }
func (v *VarGoFor) Print()                 { fmt.Print(v.Data) }

func (v *VarGoFor) initVarName() error {
	v.VarName = v.Element.GetAttrParts()[0]
	return nil
}

func (v *VarGoFor) initBasicInfo() error {
	attrParts := v.Element.GetAttrParts()
	v.VarName = attrParts[0] + "For" + v.Element.GetId()
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
		varsToWrite += inner.GetData()
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
%s := gtml.For(%s, func(i int, %s %s) string {
var %s strings.Builder
%s
%s
return %s.String()
})`+"\n", v.VarName, v.IterItems, v.IterItem, v.IterType, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
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
		varsToWrite += inner.GetData()
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
%s := gtml.If(%s, func() string {
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

// ##==================================================================
type VarGoElse struct {
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

func NewVarGoElse(elm Element) (*VarGoElse, error) {
	v := &VarGoElse{
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

func (v *VarGoElse) GetData() string        { return v.Data }
func (v *VarGoElse) GetVarName() string     { return v.VarName }
func (v *VarGoElse) GetBuilderName() string { return v.BuilderName }
func (v *VarGoElse) GetType() string        { return v.Type }
func (v *VarGoElse) Print()                 { fmt.Print(v.Data) }

func (v *VarGoElse) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "Else" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoElse
	return nil
}

func (v *VarGoElse) initVars() error {
	vars, err := GetElementVars(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoElse) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoElse) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *VarGoElse) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtml.Else(%s, func() string {
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

// ##==================================================================
type VarGoPlaceholder struct {
	Element       Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	Type          string
	ComponentName string
	Attrs         []Attr
	AttrsCalledAs []Attr
	ParamStr      string
	CallParams    []string
	CallParamStr  string
}

func NewVarGoPlaceholder(elm Element) (*VarGoPlaceholder, error) {
	v := &VarGoPlaceholder{
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

func (v *VarGoPlaceholder) GetData() string        { return v.Data }
func (v *VarGoPlaceholder) GetVarName() string     { return v.VarName }
func (v *VarGoPlaceholder) GetBuilderName() string { return v.BuilderName }
func (v *VarGoPlaceholder) GetType() string        { return v.Type }
func (v *VarGoPlaceholder) Print()                 { fmt.Print(v.Data) }

func (v *VarGoPlaceholder) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = strings.ToLower(attr) + "Placeholder" + v.Element.GetId()
	v.BuilderName = strings.ToLower(attr) + "Builder"
	v.ComponentName = attr
	v.Type = KeyVarGoPlaceholder
	return nil
}

func (v *VarGoPlaceholder) initAttrs() error {
	for _, a := range v.Element.GetSelection().Get(0).Attr {
		if strings.HasPrefix(a.Key, "_") {
			continue
		}
		attr, err := NewAttr(a.Key, a.Val)
		if err != nil {
			return err
		}
		v.Attrs = append(v.Attrs, attr)
	}
	for _, a := range v.Attrs {
		val := purse.Squeeze(a.GetValue())
		if val == "" {
			v.AttrsCalledAs = append(v.AttrsCalledAs, a)
			continue
		}
		if string(val[0]) == "@" {
			val = val[1:]
			attr, err := NewAttr(a.GetKey(), val)
			if err != nil {
				return err
			}
			v.AttrsCalledAs = append(v.AttrsCalledAs, attr)
			continue
		}
		if strings.HasPrefix(val, "{{") && strings.HasSuffix(val, "}}") {
			val = purse.RemoveAllSubStr(val, "{{", "}}")
			attr, err := NewAttr(a.GetKey(), val)
			if err != nil {
				return err
			}
			v.AttrsCalledAs = append(v.AttrsCalledAs, attr)
			continue
		}
		v.AttrsCalledAs = append(v.AttrsCalledAs, a)
	}
	return nil
}

func (v *VarGoPlaceholder) initVars() error {
	vars, err := GetElementVars(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoPlaceholder) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoPlaceholder) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *VarGoPlaceholder) initCallParams() error {
	for _, attr := range v.Attrs {
		if attr.GetType() == KeyAttrInitParam || attr.GetType() == KeyAttrAtParam {
			v.CallParams = append(v.CallParams, "ATTRID"+attr.GetKey()+"ATTRID"+attr.GetValue())
			continue
		}
		v.CallParams = append(v.CallParams, "ATTRID"+attr.GetKey()+"ATTRID\""+attr.GetValue()+"\"")
	}
	vars, err := GetElementVars(v.Element)
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

func (v *VarGoPlaceholder) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := func() string {
%s
return %s(%s)
}`+"\n", v.VarName, v.WriteVarsAs, v.ComponentName, v.CallParamStr))
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}

// ##==================================================================
type VarGoSlot struct {
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

func NewVarGoSlot(elm Element) (*VarGoSlot, error) {
	v := &VarGoSlot{
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

func (v *VarGoSlot) GetData() string        { return v.Data }
func (v *VarGoSlot) GetVarName() string     { return v.VarName }
func (v *VarGoSlot) GetBuilderName() string { return v.BuilderName }
func (v *VarGoSlot) GetType() string        { return v.Type }
func (v *VarGoSlot) Print()                 { fmt.Print(v.Data) }

func (v *VarGoSlot) initBasicInfo() error {
	attr := v.Element.GetAttr()
	v.VarName = attr + "Slot" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.BoolToCheck = attr
	v.Type = KeyVarGoSlot
	return nil
}

func (v *VarGoSlot) initVars() error {
	vars, err := GetElementVars(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *VarGoSlot) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *VarGoSlot) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *VarGoSlot) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtml.Slot(func() string {
var %s strings.Builder
%s
%s
return %s.String()
})`+"\n", v.VarName, v.BuilderName, v.WriteVarsAs, v.BuilderSeries, v.BuilderName))
	// v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}

// ##==================================================================
