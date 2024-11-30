package gtmlrune

import (
	"fmt"
	"gtml/src/parser/element"
	"strings"

	"github.com/phillip-england/purse"
)

type GtmlRune interface {
	Print()
	GetValue() string
	GetType() string
	GetDecodedData() string
	GetLocation() string
}

func NewGtmlRune(runeStr string, location string) (GtmlRune, error) {
	if strings.HasPrefix(runeStr, KeyRuneProp) {
		r, err := NewProp(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(runeStr, KeyRuneSlot) {
		r, err := NewSlot(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(runeStr, KeyRuneVal) {
		r, err := NewVal(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(runeStr, KeyRunePipe) {
		r, err := NewPipe(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, nil
}

func NewRunesFromStr(s string) ([]GtmlRune, error) {
	runes := make([]GtmlRune, 0)
	parts := purse.ScanBetweenSubStrs(s, "$", ")")
	clay := s
	for _, part := range parts {
		index := strings.Index(part, "(")
		if index == -1 {
			continue
		}
		name := part[:index]
		if !purse.SliceContains(GetRuneNames(), name) {
			continue
		}
		index = strings.Index(clay, part)
		if index == -1 {
			continue // Skip if the part is not found in `clay`
		}
		potentialEqualSignIndex := index - 2
		if potentialEqualSignIndex < 0 || potentialEqualSignIndex+2 > len(clay) {
			continue // Skip if accessing `clay[potentialEqualSignIndex : potentialEqualSignIndex+2]` would go out of bounds
		}
		potentialAttrStart := clay[potentialEqualSignIndex : potentialEqualSignIndex+2]
		attrLocation := KeyLocationElsewhere
		if potentialAttrStart == "=\"" || potentialAttrStart == "='" {
			attrLocation = KeyLocationAttribute
		}
		r, err := NewGtmlRune(part, attrLocation)
		if err != nil {
			return runes, err
		}
		runes = append(runes, r)
	}
	return runes, nil
}

func NewRunesFromElement(elm element.Element) ([]GtmlRune, error) {
	elmHtml, err := element.GetElementHtmlWithoutChildren(elm)
	if err != nil {
		return make([]GtmlRune, 0), err
	}
	rns, err := NewRunesFromStr(elmHtml)
	if err != nil {
		return rns, err
	}
	return rns, nil
}

func GetRunesAsWriteStringCalls(elm element.Element, builderName string) ([]string, error) {
	calls := make([]string, 0)
	runes, err := NewRunesFromElement(elm)
	if err != nil {
		return calls, err
	}
	for _, rn := range runes {
		if rn.GetType() == KeyRuneProp {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			calls = append(calls, call)
		}
		if rn.GetType() == KeyRuneVal {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			calls = append(calls, call)
		}
		if rn.GetType() == KeyRunePipe {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			calls = append(calls, call)
		}
		if rn.GetType() == KeyRuneSlot {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, rn.GetValue())
			calls = append(calls, call)
		}
	}
	return calls, nil
}
