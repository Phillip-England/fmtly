package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyPropForType = "FORTYPE"
	KeyPropForStr  = "FORSTR"
	KeyPropStr     = "STR"
)

// ##==================================================================
type Prop interface {
	GetRaw() string
	GetValue() string
	GetType() string
	Print()
}

func NewProps(elm Element) ([]Prop, error) {
	strProps := purse.ScanBetweenSubStrs(elm.GetHtml(), "{{", "}}")
	props := make([]Prop, 0)
	for _, prop := range strProps {
		val := purse.Squeeze(prop)
		val = purse.RemoveAllSubStr(val, "{{", "}}")
		if val == "" {
			return nil, fmt.Errorf("empty prop tag provided: %s", elm.GetHtml())
		}
		if strings.Count(val, ".") == 1 && len(strings.Split(val, ".")) == 2 {
			propForType, err := NewPropForType(prop, val)
			if err != nil {
				return nil, err
			}
			props = append(props, propForType)
			continue
		}
		propStr, err := NewPropStr(prop, val)
		if err != nil {
			return nil, err
		}
		props = append(props, propStr)
	}
	morphed := make([]Prop, 0)
	for _, prop := range props {
		if prop.GetType() == KeyPropStr {
			// check if an Element's children are _for
			err := WalkElementChildren(elm, func(child Element) error {
				if child.GetType() == KeyElementFor {
					attrParts := child.GetAttrParts()
					iterItem := attrParts[0]
					// if a for element of _for="item of items []string" exists, and a StrProp has the value of "item", then it is a ForStrProp
					if prop.GetValue() == iterItem {
						forStrProp, err := NewPropForStr(prop.GetRaw(), prop.GetValue())
						if err != nil {
							return err
						}
						morphed = append(morphed, forStrProp)
						return nil
					}
					morphed = append(morphed, prop)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			continue
		}
		morphed = append(morphed, prop)
	}
	return morphed, nil
}

func PropAsWriteString(prop Prop, builderName string) string {
	return fmt.Sprintf("%s.WriteString(%s)", builderName, prop.GetValue())
}

// ##==================================================================
type PropForType struct {
	Raw   string
	Value string
	Type  string
}

func NewPropForType(raw string, val string) (*PropForType, error) {
	prop := &PropForType{
		Raw:   raw,
		Value: val,
		Type:  KeyPropForType,
	}
	return prop, nil
}

func (prop *PropForType) GetRaw() string   { return prop.Raw }
func (prop *PropForType) GetValue() string { return prop.Value }
func (prop *PropForType) GetType() string  { return prop.Type }
func (prop *PropForType) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s", prop.Raw, prop.Value))
}

// ##==================================================================
type PropForStr struct {
	Raw   string
	Value string
	Type  string
}

func NewPropForStr(raw string, val string) (*PropForStr, error) {
	prop := &PropForStr{
		Raw:   raw,
		Value: val,
		Type:  KeyPropForStr,
	}
	return prop, nil
}

func (prop *PropForStr) GetRaw() string   { return prop.Raw }
func (prop *PropForStr) GetValue() string { return prop.Value }
func (prop *PropForStr) GetType() string  { return prop.Type }
func (prop *PropForStr) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s", prop.Raw, prop.Value))
}

// ##==================================================================
type PropStr struct {
	Raw   string
	Value string
	Type  string
}

func NewPropStr(raw string, val string) (*PropStr, error) {
	prop := &PropStr{
		Raw:   raw,
		Value: val,
		Type:  KeyPropStr,
	}
	return prop, nil
}

func (prop *PropStr) GetRaw() string   { return prop.Raw }
func (prop *PropStr) GetValue() string { return prop.Value }
func (prop *PropStr) GetType() string  { return prop.Type }
func (prop *PropStr) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s", prop.Raw, prop.Value))
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
