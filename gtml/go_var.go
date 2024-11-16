package gtml

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
type GoVar interface {
	GetData() string
}

func NewGoVar(elm Element) (GoVar, error) {
	match := GetElementType(elm)
	switch match {
	case "for":
		v, err := NewGoLoopVar(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	htmlStr, err := GetElementHtml(elm)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("element does not corrospond to a valid GoToken: %s", htmlStr)
}

func PrintGoVar(v GoVar) {
	fmt.Println(v.GetData())
}

func GetGoVarName(v GoVar) (string, error) {
	parts := strings.Split(v.GetData(), ":=")
	if len(parts) == 0 {
		return "", fmt.Errorf("GoVar does not contain a := symbol, so we cannot parse the name: %s", v.GetData())
	}
	firstPart := purse.Squeeze(parts[0])
	return firstPart, nil
}

// ##==================================================================
type GoLoopVar struct {
	Element     Element
	VarName     string
	BuilderName string
	Vars        []GoVar
	WriteVarsAs string
	Data        string
	IterItems   string
	IterItem    string
	IterType    string
}

func NewGoLoopVar(elm Element) (*GoLoopVar, error) {
	v := &GoLoopVar{
		Element: elm,
	}
	err := fungi.Process(
		func() error { return v.initBasicInfo() },
		func() error { return v.initVars() },
		func() error { return v.initWriteVarsAs() },
		func() error { return v.initData() },
	)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (v *GoLoopVar) GetData() string { return v.Data }

func (v *GoLoopVar) initBasicInfo() error {
	attrParts, err := ForceElementAttrParts(v.Element, "_for", 4)
	if err != nil {
		return err
	}
	v.VarName = attrParts[0] + "Loop"
	v.BuilderName = attrParts[0] + "Builder"
	v.IterItems = attrParts[2]
	v.IterItem = attrParts[0]
	v.IterType = purse.RemoveAllSubStr(attrParts[3], "[]")
	return nil
}

func (v *GoLoopVar) initVars() error {
	err := WalkElementDirectChildren(v.Element, func(child Element) error {
		innerVar, err := NewGoVar(child)
		if err != nil {
			return nil
		}
		v.Vars = append(v.Vars, innerVar)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (v *GoLoopVar) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += "\t" + inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoLoopVar) initData() error {
	htmlStr, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := collect(%s, func(i int, %s %s) string {
	var %s strings.Builder
%s
%s
	return %s.String()
})`, v.VarName, v.IterItems, v.IterItem, v.IterType, v.BuilderName, v.WriteVarsAs, htmlStr, v.BuilderName))
	code, err := format.Source([]byte(v.Data))
	if err != nil {
		return err
	}
	v.Data = string(code)
	v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
