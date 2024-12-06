package gtmlvar

import (
	"fmt"
	"gtml/src/parser/element"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type GoMd struct {
	Element       element.Element
	VarName       string
	BuilderName   string
	Vars          []Var
	WriteVarsAs   string
	Data          string
	BuilderSeries string
	Type          string
	MdFilePath    string
	MdTheme       string
}

func NewGoMd(elm element.Element) (*GoMd, error) {
	v := &GoMd{
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

func (v *GoMd) GetData() string             { return v.Data }
func (v *GoMd) GetVarName() string          { return v.VarName }
func (v *GoMd) GetBuilderName() string      { return v.BuilderName }
func (v *GoMd) GetType() string             { return v.Type }
func (v *GoMd) GetElement() element.Element { return v.Element }
func (v *GoMd) Print()                      { fmt.Print(v.Data) }

func (v *GoMd) initVarName() error {
	v.VarName = v.Element.GetAttrParts()[0]
	return nil
}

func (v *GoMd) initBasicInfo() error {
	attr := v.Element.GetAttr()
	attr = strings.ReplaceAll(attr, "/", "")
	attr = strings.ReplaceAll(attr, ".", "")
	v.VarName = attr + "Md" + v.Element.GetId()
	v.BuilderName = attr + "Builder"
	v.MdFilePath = v.Element.GetAttr()
	sel := v.Element.GetSelection()
	theme, exists := sel.Attr("_md-theme")
	if !exists {
		theme = "dracula"
	}
	v.MdTheme = theme
	v.Type = KeyVarGoMd
	return nil
}

func (v *GoMd) initVars() error {
	vars, err := NewVarsFromElement(v.Element)
	if err != nil {
		return err
	}
	v.Vars = vars
	return nil
}

func (v *GoMd) initWriteVarsAs() error {
	varsToWrite := ""
	for _, inner := range v.Vars {
		varsToWrite += inner.GetData()
	}
	v.WriteVarsAs = varsToWrite
	return nil
}

func (v *GoMd) initBuilderSeries() error {
	series, err := GetElementAsBuilderSeries(v.Element, v.BuilderName)
	if err != nil {
		return err
	}
	v.BuilderSeries = series
	return nil
}

func (v *GoMd) initData() error {
	v.Data = purse.RemoveFirstLine(fmt.Sprintf(`
%s := gtmlMd("%s", "%s")`+"\n", v.VarName, v.MdFilePath, v.MdTheme))
	// v.Data = purse.RemoveEmptyLines(v.Data)
	return nil
}
