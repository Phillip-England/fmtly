package gtml

import (
	"fmt"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"

	"github.com/phillip-england/purse"
)

// ##==================================================================
type Func interface {
	GetData() string
	SetData(str string)
	GetVars() []Var
}

func NewFunc(elm Element) (Func, error) {
	if elm.GetType() == "COMPONENT" {
		fn, err := NewGoComponentFunc(elm)
		if err != nil {
			return nil, err
		}
		return fn, nil
	}
	return nil, fmt.Errorf("provided element does not corrospond to a valid GoFunc: %s", elm.GetHtml())
}

func PrintGoFunc(fn Func) {
	fmt.Println(fn.GetData())
}

// ##==================================================================
type GoComponentFunc struct {
	Element  Element
	Vars     []Var
	Data     string
	VarStr   string
	Name     string
	ParamStr string
}

func NewGoComponentFunc(elm Element) (*GoComponentFunc, error) {
	fn := &GoComponentFunc{
		Element: elm,
	}
	err := fungi.Process(
		func() error { return fn.initName() },
		func() error { return fn.initVars() },
		func() error { return fn.initVarStr() },
		func() error { return fn.initParamStr() },
		func() error { return fn.initData() },
	)
	if err != nil {
		return nil, err
	}
	_, err = GetElementParams(fn.Element)
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func (fn *GoComponentFunc) GetData() string    { return fn.Data }
func (fn *GoComponentFunc) SetData(str string) { fn.Data = str }
func (fn *GoComponentFunc) GetVars() []Var     { return fn.Vars }

func (fn *GoComponentFunc) initName() error {
	compAttr, err := gqpp.ForceElementAttr(fn.Element.GetSelection(), "_component")
	if err != nil {
		return err
	}
	fn.Name = compAttr
	return nil
}

func (fn *GoComponentFunc) initVars() error {
	err := WalkElementDirectChildren(fn.Element, func(child Element) error {
		goVar, err := NewGoVar(child)
		if err != nil {
			return err
		}
		fn.Vars = append(fn.Vars, goVar)
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (fn *GoComponentFunc) initVarStr() error {
	str := ""
	for _, v := range fn.Vars {
		data := v.GetData()
		str += data + "\n"
	}
	str = purse.PrefixLines(str, "\t")
	fn.VarStr = str
	return nil
}

func (fn *GoComponentFunc) initParamStr() error {
	params, err := GetElementParams(fn.Element)
	if err != nil {
		return err
	}
	fn.ParamStr = params
	return nil
}

func (fn *GoComponentFunc) initData() error {
	series, err := GetElementAsBuilderSeries(fn.Element, "builder")
	if err != nil {
		return err
	}
	series = purse.PrefixLines(series, "\t")
	data := purse.RemoveFirstLine(fmt.Sprintf(`
func %s(%s) string {
	var builder strings.Builder
%s
%s
	return builder.String()
}
	`, fn.Name, fn.ParamStr, fn.VarStr, series))
	data = purse.RemoveEmptyLines(data)
	fn.Data = data
	return nil
}
