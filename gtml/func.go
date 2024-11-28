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
	if elm.GetType() == KeyElementComponent || elm.GetType() == KeyElementPlaceholder {
		fn, err := NewGoComponentFunc(elm, filtered)
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
		func() error { return fn.initParams() },
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
	fn.VarStr = str
	return nil
}

func (fn *GoComponentFunc) initParams() error {
	params, err := GetElementParams(fn.Element)
	if err != nil {
		return err
	}
	fn.Params = params
	strs := make([]string, 0)
	for _, param := range params {
		strs = append(strs, param.GetStr())
	}
	fn.ParamStr = strings.Join(strs, ", ")
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
	// Initialize an empty slice to hold the ordered placeholder call strings.
	ordered := make([]string, 0)

	// Iterate through all sibling elements.
	for _, sib := range siblings {
		// Skip processing if the sibling element has the same name as the current element.
		if fn.Element.GetName() == sib.GetName() {
			continue
		}

		// Retrieve parameters for the sibling element.
		params, err := GetElementParams(sib)
		if err != nil {
			return err // Return any error encountered during parameter retrieval.
		}

		// Initialize slices for unique sibling parameters and already processed parameter names.
		sibParams := make([]Param, 0)
		found := make([]string, 0)

		// Filter out duplicate parameters from the sibling element.
		for _, param := range params {
			if purse.SliceContains(found, param.GetStr()) {
				continue // Skip if the parameter has already been processed.
			}
			sibParams = append(sibParams, param) // Add unique parameters.
			found = append(found, param.GetStr())
		}

		// Iterate through each unique sibling parameter.
		for _, sibParam := range sibParams {
			// Loop through all placeholder calls in the current function.
			for _, call := range fn.PlaceholderCalls {
				// Retrieve parameters for the current placeholder call.
				callParams := call.GetParams()

				// Process each parameter in the call.
				for _, callParam := range callParams {
					clay := callParam

					// Clean up the parameter string by removing "ATTRID" substrings.
					clay = strings.Replace(clay, "ATTRID", "", 1)
					clay = strings.Replace(clay, "ATTRID", " ", 1)

					// Split the cleaned parameter string into parts.
					callParamParts := strings.Split(clay, " ")

					// Filter out empty parts from the split results.
					filtered := make([]string, 0)
					for _, part := range callParamParts {
						if len(part) == 0 {
							continue
						}
						filtered = append(filtered, part)
					}
					callParamParts = filtered

					// Extract the identifier (first part) of the cleaned parameter string.
					callParamId := callParamParts[0]
					sibParamName := sibParam.GetName()

					// Check if the call parameter ID matches the sibling parameter name.
					if callParamId == sibParamName {
						// Extract the part of the call parameter string after "ATTRID" and append it to the ordered list.
						writeAs := strings.Replace(callParam, "ATTRID", "", 1)
						i := strings.Index(writeAs, "ATTRID") + len("ATTRID")
						writeAs = writeAs[i:]
						ordered = append(ordered, writeAs)
					}
				}
			}
		}
	}

	// Set the ordered placeholder calls for the function.
	fn.OrderedPlaceholderCalls = ordered

	// Return nil to indicate successful execution.
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
