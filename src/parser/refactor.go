package parser

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlrune"
	"gtml/src/parser/gtmlvar"
	"strings"

	"github.com/phillip-england/purse"
)

func GetElementAsBuilderSeries(elm element.Element, builderName string) (string, error) {
	clay := elm.GetHtml()
	err := element.WalkElementDirectChildren(elm, func(child element.Element) error {
		childHtml := child.GetHtml()
		newVar, err := gtmlvar.NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, gtmlvar.KeyVarGoElse, gtmlvar.KeyVarGoFor, gtmlvar.KeyVarGoIf, gtmlvar.KeyVarGoPlaceholder, gtmlvar.KeyVarGoSlot) {
			if varType == gtmlvar.KeyVarGoPlaceholder {
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
	err = element.WalkElementRunes(elm, func(rn gtmlrune.GtmlRune) error {
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
		return nil
	})
	if err != nil {
		return "", err
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
