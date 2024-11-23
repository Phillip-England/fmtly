package gtml

import (
	"fmt"
	"go/format"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"

	"github.com/phillip-england/purse"
)

// ##==================================================================
type Func interface {
	GetData() string
	SetData(str string)
	GetVars() []Var
	Print()
}

func NewFunc(elm Element) (Func, error) {
	if elm.GetType() == KeyElementComponent {
		fn, err := NewGoComponentFunc(elm)
		if err != nil {
			return nil, err
		}
		return fn, nil
	}
	return nil, fmt.Errorf("provided element does not corrospond to a valid GoFunc: %s", elm.GetHtml())
}

// ##==================================================================
type GoComponentFunc struct {
	Element      Element
	Vars         []Var
	BuilderNames []string
	Data         string
	VarStr       string
	Name         string
	ParamStr     string
	BuilderCalls []string
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
		func() error { return fn.initBuilderNames() },
		func() error { return fn.initBuilderCalls() },
	)
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func (fn *GoComponentFunc) GetData() string    { return fn.Data }
func (fn *GoComponentFunc) SetData(str string) { fn.Data = str }
func (fn *GoComponentFunc) GetVars() []Var     { return fn.Vars }
func (fn *GoComponentFunc) Print()             { fmt.Println(fn.GetData()) }
func (fn *GoComponentFunc) initName() error {
	compAttr, err := gqpp.ForceElementAttr(fn.Element.GetSelection(), KeyElementComponent)
	if err != nil {
		return err
	}
	fn.Name = compAttr
	return nil
}

func (fn *GoComponentFunc) initVars() error {
	err := WalkElementDirectChildren(fn.Element, func(child Element) error {
		goVar, err := NewVar(child)
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
	paramStrs := make([]string, 0)
	for _, param := range params {
		paramStrs = append(paramStrs, param.GetStr())
	}
	paramStrs = purse.RemoveDuplicatesInSlice(paramStrs)
	fn.ParamStr = strings.Join(paramStrs, ",")
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
	code, err := format.Source([]byte(data))
	if err != nil {
		return err
	}
	data = string(code)
	data = purse.RemoveEmptyLines(data)
	fn.Data = data
	return nil
}

func (fn *GoComponentFunc) initBuilderNames() error {
	fn.BuilderNames = append(fn.BuilderNames, "builder")
	for _, v := range fn.Vars {
		fn.BuilderNames = append(fn.BuilderNames, v.GetBuilderName())
	}
	return nil
}

func (fn *GoComponentFunc) initBuilderCalls() error {
	lines := purse.MakeLines(fn.GetData())
	for _, line := range lines {
		line = purse.Flatten(line)
		for _, name := range fn.BuilderNames {
			if strings.HasPrefix(line, name+".") {
				fn.BuilderCalls = append(fn.BuilderCalls, line)
			}
		}
	}
	return nil
}
