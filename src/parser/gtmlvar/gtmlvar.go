package gtmlvar

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlrune"
	"strings"

	"github.com/phillip-england/purse"
)

type Var interface {
	GetData() string
	GetVarName() string
	GetType() string
	Print()
	GetBuilderName() string
	GetElement() element.Element
}

func NewVar(elm element.Element) (Var, error) {
	switch elm.GetType() {
	case element.KeyElementComponent:
		v, err := NewGoComponent(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementFor:
		v, err := NewGoFor(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementIf:
		v, err := NewGoIf(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementElse:
		v, err := NewGoElse(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementPlaceholder:
		v, err := NewGoPlaceholder(elm)
		if err != nil {
			return nil, err
		}
		return v, nil
	case element.KeyElementSlot:
		v, err := NewGoSlot(elm)
		if err != nil {
			return nil, err
		}
		return v, nil

	}

	return nil, fmt.Errorf("element does not corrospond to a valid Var: %s", elm.GetHtml())
}

func NewVarsFromElement(elm element.Element) ([]Var, error) {
	vars := make([]Var, 0)
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		innerVar, err := NewVar(child)
		if err != nil {
			return nil
		}
		vars = append(vars, innerVar)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func GetVarsAsWriteStringCalls(elm element.Element, builderName string) ([]string, error) {
	calls := make([]string, 0)
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		newVar, err := NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, GetFullVarList()...) {
			if varType == KeyVarGoPlaceholder {
				call := fmt.Sprintf("%s.WriteString(%s())", builderName, newVar.GetVarName())
				calls = append(calls, call)
				return nil
			}
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, newVar.GetVarName())
			calls = append(calls, call)
		}
		return nil
	})
	if err != nil {
		return calls, err
	}
	return calls, nil
}

func GetElementAsBuilderSeries(elm element.Element, builderName string) (string, error) {
	clay := elm.GetHtml()
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		childHtml := child.GetHtml()
		newVar, err := NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, KeyVarGoElse, KeyVarGoFor, KeyVarGoIf, KeyVarGoPlaceholder, KeyVarGoSlot) {
			if varType == KeyVarGoPlaceholder {
				call := fmt.Sprintf("%s.WriteString(%s())", builderName, newVar.GetVarName())
				clay = strings.Replace(clay, childHtml, call, 1)
			}
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, newVar.GetVarName())
			clay = strings.Replace(clay, childHtml, call, 1)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	runes, err := gtmlrune.NewRunesFromElement(elm)
	if err != nil {
		return "", err
	}
	for _, rn := range runes {
		if rn.GetType() == gtmlrune.KeyRuneProp {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			clay = strings.Replace(clay, rn.GetDecodedData(), call, 1)
		}
		if rn.GetType() == gtmlrune.KeyRuneVal {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			clay = strings.Replace(clay, rn.GetDecodedData(), call, 1)
		}
		if rn.GetType() == gtmlrune.KeyRunePipe {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			clay = strings.Replace(clay, rn.GetDecodedData(), call, 1)
		}
		if rn.GetType() == gtmlrune.KeyRuneSlot {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			clay = strings.Replace(clay, rn.GetDecodedData(), call, 1)
		}
	}

	if strings.Index(clay, builderName) == -1 {
		singleCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		return singleCall, nil
	}
	series := ""
	for {
		builderIndex := strings.Index(clay, builderName)
		if builderIndex == -1 {
			break
		}
		htmlPart := clay[:builderIndex]
		if htmlPart != "" {
			htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, htmlPart)
			series += htmlCall + "\n"
			clay = strings.Replace(clay, htmlPart, "", 1)
		}
		endBuilderIndex := strings.Index(clay, ")")
		loopCount := 0
		for {
			loopCount++
			nextChar := string(clay[endBuilderIndex+loopCount])
			if nextChar == ")" {
				endBuilderIndex = endBuilderIndex + loopCount
				continue
			}
			break
		}
		builderPart := clay[:endBuilderIndex+1]
		series += builderPart + "\n"
		clay = strings.Replace(clay, builderPart, "", 1)
	}
	if len(clay) > 0 {
		htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		series += htmlCall + "\n"
	}
	return series, nil
}
