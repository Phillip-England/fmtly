package gtmlfunc

import (
	"fmt"
	"go/format"
	"gtml/src/parser/call"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlvar"
	"gtml/src/parser/param"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

type GoComponentFunc struct {
	Element                 element.Element
	Vars                    []gtmlvar.Var
	BuilderNames            []string
	Data                    string
	VarStr                  string
	Name                    string
	Params                  []param.Param
	ParamStr                string
	BuilderCalls            []string
	ReturnCalls             []string
	PlaceholderCalls        []call.Call
	OrderedPlaceholderCalls []string
}

func NewGoComponentFunc(elm element.Element, siblings []element.Element) (*GoComponentFunc, error) {
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
func (fn *GoComponentFunc) GetData() string          { return fn.Data }
func (fn *GoComponentFunc) SetData(str string)       { fn.Data = str }
func (fn *GoComponentFunc) GetVars() []gtmlvar.Var   { return fn.Vars }
func (fn *GoComponentFunc) GetParams() []param.Param { return fn.Params }
func (fn *GoComponentFunc) Print()                   { fmt.Println(fn.GetData()) }

func (fn *GoComponentFunc) initName() error {
	compAttr, err := gqpp.ForceElementAttr(fn.Element.GetSelection(), element.KeyElementComponent)
	if err != nil {
		return err
	}
	fn.Name = compAttr
	return nil
}

func (fn *GoComponentFunc) initVars() error {
	if fn.Element.GetType() == element.KeyElementPlaceholder {
		goVar, err := gtmlvar.NewVar(fn.Element)
		if err != nil {
			return err
		}
		fn.Vars = append(fn.Vars, goVar)
	}
	err := element.WalkElementDirectChildren(fn.Element, func(child element.Element) error {
		goVar, err := gtmlvar.NewVar(child)
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
		if fn.Element.GetType() == element.KeyElementPlaceholder {
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
	params, err := param.NewParamsFromElement(fn.Element)
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
	goVar, err := gtmlvar.NewVar(fn.Element)
	if err != nil {
		return err
	}
	series := goVar.GetData()

	returnCall := ""
	if fn.Element.GetType() == element.KeyElementComponent {
		returnCall = fmt.Sprintf("%s()", strings.ToLower(fn.Name))
	}
	if fn.Element.GetType() == element.KeyElementPlaceholder {
		returnCall = fmt.Sprintf("%s()", goVar.GetVarName())
	}

	data := purse.RemoveFirstLine(fmt.Sprintf(`
func %s(%s) string {
%s
return gtmlEscape(%s)
}
`, fn.Name, fn.ParamStr, series, returnCall))
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
	for _, returnCall := range fn.ReturnCalls {
		if strings.Contains(returnCall, "ATTRID") {
			returnCall = strings.Replace(returnCall, "return ", "", 1)
			newCall, err := call.NewCall(returnCall)
			if err != nil {
				return err
			}
			fn.PlaceholderCalls = append(fn.PlaceholderCalls, newCall)
		}
	}
	return nil
}

func (fn *GoComponentFunc) initOrderPlaceholderCalls(siblings []element.Element) error {
	// Initialize an empty slice to hold the ordered placeholder call strings.
	ordered := make([]string, 0)

	// Iterate through all sibling elements.
	for _, sib := range siblings {
		// Skip processing if the sibling element has the same name as the current element.
		if fn.Element.GetName() == sib.GetName() {
			continue
		}

		// Retrieve parameters for the sibling element.
		params, err := param.NewParamsFromElement(sib)
		if err != nil {
			return err // Return any error encountered during parameter retrieval.
		}

		// Initialize slices for unique sibling parameters and already processed parameter names.
		sibParams := make([]param.Param, 0)
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
						if writeAs == "\"true\"" || writeAs == "\"false\"" {
							writeAs = strings.ReplaceAll(writeAs, "\"", "")
						}
						if writeAs == "\"\\\\true\"" || writeAs == "\"\\\\false\"" {
							writeAs = strings.ReplaceAll(writeAs, "\\", "")
						}
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
