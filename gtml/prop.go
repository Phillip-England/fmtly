package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
type Prop interface {
	GetRaw() string
	GetValue() string
	GetType() string
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
	filtered := make([]Prop, 0)
	for _, prop := range props {
		if prop.GetType() == "STR" {
			err := WalkElementChildren(elm, func(child Element) error {
				if elm.GetType() == "FOR" {
					attrParts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), "_for", 4)
					if err != nil {
						return err
					}
					iterItem := attrParts[0]
					// here, we filter out any PropStr which is found to be in a _for="str of strs []string"
					if prop.GetValue() == iterItem {
						return nil
					}
					filtered = append(filtered, prop)
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			continue
		}
		filtered = append(filtered, prop)
	}
	return filtered, nil
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
		Type:  "FORTYPE",
	}
	return prop, nil
}

func (prop *PropForType) GetRaw() string   { return prop.Raw }
func (prop *PropForType) GetValue() string { return prop.Value }
func (prop *PropForType) GetType() string  { return prop.Type }

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
		Type:  "FORSTR",
	}
	return prop, nil
}

func (prop *PropForStr) GetRaw() string   { return prop.Raw }
func (prop *PropForStr) GetValue() string { return prop.Value }
func (prop *PropForStr) GetType() string  { return prop.Type }

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
		Type:  "STR",
	}
	return prop, nil
}

func (prop *PropStr) GetRaw() string   { return prop.Raw }
func (prop *PropStr) GetValue() string { return prop.Value }
func (prop *PropStr) GetType() string  { return prop.Type }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
