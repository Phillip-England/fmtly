package gtml

import (
	"fmt"
	"go/format"
	"slices"
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
	GetParams() []Param
	Print()
}

func NewFunc(elm Element, siblings []Element) (Func, error) {
	filtered := make([]Element, 0)
	for _, sibling := range siblings {
		if sibling.GetName() == elm.GetName() {
			continue
		}
		filtered = append(filtered, sibling)
	}
	siblings = filtered
	if elm.GetType() == KeyElementComponent || elm.GetType() == KeyElementPlaceholder {
		fn, err := NewGoComponentFunc(elm, siblings)
		if err != nil {
			return nil, err
		}
		return fn, nil
	}
	return nil, fmt.Errorf("provided element does not corrospond to a valid GoFunc: %s", elm.GetHtml())
}

// ##==================================================================
type GoComponentFunc struct {
	Element                 Element
	Vars                    []Var
	BuilderNames            []string
	Data                    string
	VarStr                  string
	Name                    string
	Params                  []Param
	ParamStr                string
	BuilderCalls            []string
	ReturnCalls             []string
	PlaceholderCalls        []Call
	OrderedPlaceholderCalls []string
}

func NewGoComponentFunc(elm Element, siblings []Element) (*GoComponentFunc, error) {
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
		func() error { return fn.initReturnCalls() },
		func() error { return fn.initPlaceholderCalls() },
		func() error { return fn.initOrderPlaceholderCalls(siblings) },
		func() error { return fn.initWriteCorrectPlaceholderCalls() },
		func() error { return fn.initFormatData() },
	)
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func (fn *GoComponentFunc) GetData() string    { return fn.Data }
func (fn *GoComponentFunc) SetData(str string) { fn.Data = str }
func (fn *GoComponentFunc) GetVars() []Var     { return fn.Vars }
func (fn *GoComponentFunc) GetParams() []Param { return fn.Params }
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
	if fn.Element.GetType() == KeyElementPlaceholder {
		goVar, err := NewVar(fn.Element)
		if err != nil {
			return err
		}
		fn.Vars = append(fn.Vars, goVar)
	}
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
	for i, v := range fn.Vars {
		// if a _component is also a _placeholder, we must only include the first var in its var string
		// not doing so will result in the outer func body becoming polluted with unneeded vars from the _placeholder element itself.
		// its vars will become present in the func body
		if fn.Element.GetType() == KeyElementPlaceholder {
			if i == 0 {
				data := v.GetData()
				str += data + "\n"
			}
		} else {
			data := v.GetData()
			str += data + "\n"
		}
	}
	// str = purse.PrefixLines(str, "\t")
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
		fn.Params = append(fn.Params, param)
		paramStrs = append(paramStrs, param.GetStr())
	}
	// paramStrs = purse.RemoveDuplicatesInSlice(paramStrs)
	filter := make([]string, 0)
	for _, str := range paramStrs {
		if slices.Contains(filter, str) {
			continue
		}
		filter = append(filter, str)
	}
	fn.ParamStr = strings.Join(filter, ", ")
	return nil
}

func (fn *GoComponentFunc) initData() error {
	series := ""
	if fn.Element.GetType() == KeyElementComponent {
		str, err := GetElementAsBuilderSeries(fn.Element, "builder")
		if err != nil {
			return err
		}
		series = str
	}
	if fn.Element.GetType() == KeyElementPlaceholder {
		goVar, err := NewVar(fn.Element)
		if err != nil {
			return err
		}
		series += "builder.WriteString(" + goVar.GetVarName() + "())"
	}
	// series = purse.PrefixLines(series, "\t")
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

func (fn *GoComponentFunc) initReturnCalls() error {
	lines := purse.MakeLines(fn.GetData())
	for _, line := range lines {
		line = purse.Flatten(line)
		if strings.Contains(line, "return ") {
			fn.ReturnCalls = append(fn.ReturnCalls, line)
		}
	}
	return nil
}

func (fn *GoComponentFunc) initPlaceholderCalls() error {
	for _, call := range fn.ReturnCalls {
		if strings.Contains(call, "ATTRID") {
			call = strings.Replace(call, "return ", "", 1)
			newCall, err := NewCall(call)
			if err != nil {
				return err
			}
			fn.PlaceholderCalls = append(fn.PlaceholderCalls, newCall)
		}
	}
	return nil
}

func (fn *GoComponentFunc) initOrderPlaceholderCalls(siblings []Element) error {
	ordered := make([]string, 0)
	for _, sib := range siblings {
		sibFunc, err := NewFunc(sib, siblings)
		if err != nil {
			return err
		}
		for _, sibParam := range sibFunc.GetParams() {
			for _, call := range fn.PlaceholderCalls {
				callParams := call.GetParams()
				for _, callParam := range callParams {
					clay := callParam
					clay = strings.Replace(clay, "ATTRID", "", 1)
					clay = strings.Replace(clay, "ATTRID", " ", 1)
					callParamParts := strings.Split(clay, " ")
					filtered := make([]string, 0)
					for _, part := range callParamParts {
						if len(part) == 0 {
							continue
						}
						filtered = append(filtered, part)
					}
					callParamParts = filtered
					// if len(callParamParts) != 2 {
					// 	return fmt.Errorf("somehow, we ended up with an attribute in our Component Func which is not wrapped in ATTRID: %s", call)
					// }
					callParamId := callParamParts[0]
					sibParamName := sibParam.GetName()
					if callParamId == sibParamName {
						writeAs := strings.Replace(callParam, "ATTRID", "", 1)
						i := strings.Index(writeAs, "ATTRID") + len("ATTRID")
						writeAs = writeAs[i:]
						ordered = append(ordered, writeAs)
					}
				}
			}
		}
		fn.OrderedPlaceholderCalls = ordered
	}
	return nil
}

func (fn *GoComponentFunc) initWriteCorrectPlaceholderCalls() error {
	for _, call := range fn.PlaceholderCalls {
		callStr := call.GetData()
		i := strings.Index(callStr, "(")
		callName := callStr[:i]
		paramStr := strings.Join(fn.OrderedPlaceholderCalls, ", ")
		fnCall := fmt.Sprintf(`%s(%s)`, callName, paramStr)
		fn.Data = strings.Replace(fn.Data, callStr, fnCall, 1)
	}
	return nil
}

func (fn *GoComponentFunc) initFormatData() error {
	newLines := make([]string, 0)
	lines := purse.MakeLines(fn.Data)
	indentCount := 0
	for _, line := range lines {
		if indentCount < 0 {
			break
		}
		tabs := strings.Repeat("\t", indentCount)
		newLines = append(newLines, tabs+line)
		if strings.HasSuffix(line, "{") {
			indentCount++
		}
		if strings.Contains(line, "return ") {
			indentCount--
		}

	}
	data := purse.JoinLines(newLines)
	code, err := format.Source([]byte(data))
	if err != nil {
		return err
	}
	data = string(code)
	fn.Data = data
	return nil
}
