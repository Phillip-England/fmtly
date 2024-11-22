package gtml

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyAttrPlaceholder = "ATTRPLACEHOLDER"
	KeyAttrEmpty       = "ATTREMPTY"
	KeyAttrInt         = "ATTRINT"
	KeyAttrBool        = "ATTRBOOL"
	KeyAttrStr         = "ATTRSTR"
	KeyAttrAtParam     = "ATTRATPARAM"
)

// ##==================================================================
type Attr interface {
	Print()
	GetKey() string
	GetValue() string
	GetKeyValuePair() (string, string)
	GetType() string
}

func NewAttr(key string, value string) (Attr, error) {
	if value == "" {
		attr, err := NewAttrEmpty(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	sqValue := purse.Squeeze(value)
	firstChar := string(sqValue[0])
	if firstChar == "@" {
		attr, err := NewAttrAtParam(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	if strings.HasPrefix(sqValue, "{{") && strings.HasSuffix(sqValue, "}}") {
		attr, err := NewAttrPlaceholder(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	if sqValue == "true" || sqValue == "false" {
		attr, err := NewAttrBool(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	_, err := strconv.Atoi(sqValue)
	// if the attr can be converted to an int
	if err == nil {
		attr, err := NewAttrInt(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	attr, err := NewAttrStr(key, value)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

// ##==================================================================
type AttrPlaceholder struct {
	Key   string
	Value string
	Type  string
}

func NewAttrPlaceholder(key string, value string) (*AttrPlaceholder, error) {
	attr := &AttrPlaceholder{
		Key:   key,
		Value: value,
		Type:  KeyAttrPlaceholder,
	}
	return attr, nil
}

func (attr *AttrPlaceholder) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrPlaceholder) GetKey() string                    { return attr.Key }
func (attr *AttrPlaceholder) GetValue() string                  { return attr.Value }
func (attr *AttrPlaceholder) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrPlaceholder) GetType() string                   { return attr.Type }

// ##==================================================================
type AttrEmpty struct {
	Key   string
	Value string
	Type  string
}

func NewAttrEmpty(key string, value string) (*AttrEmpty, error) {
	attr := &AttrEmpty{
		Key:   key,
		Value: value,
		Type:  KeyAttrEmpty,
	}
	return attr, nil
}

func (attr *AttrEmpty) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrEmpty) GetKey() string                    { return attr.Key }
func (attr *AttrEmpty) GetValue() string                  { return attr.Value }
func (attr *AttrEmpty) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrEmpty) GetType() string                   { return attr.Type }

// ##==================================================================
type AttrAtParam struct {
	Key   string
	Value string
	Type  string
}

func NewAttrAtParam(key string, value string) (*AttrEmpty, error) {
	attr := &AttrEmpty{
		Key:   key,
		Value: value,
		Type:  KeyAttrAtParam,
	}
	return attr, nil
}

func (attr *AttrAtParam) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrAtParam) GetKey() string                    { return attr.Key }
func (attr *AttrAtParam) GetValue() string                  { return attr.Value }
func (attr *AttrAtParam) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrAtParam) GetType() string                   { return attr.Type }

// ##==================================================================
type AttrStr struct {
	Key   string
	Value string
	Type  string
}

func NewAttrStr(key string, value string) (*AttrEmpty, error) {
	attr := &AttrEmpty{
		Key:   key,
		Value: value,
		Type:  KeyAttrStr,
	}
	return attr, nil
}

func (attr *AttrStr) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrStr) GetKey() string                    { return attr.Key }
func (attr *AttrStr) GetValue() string                  { return attr.Value }
func (attr *AttrStr) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrStr) GetType() string                   { return attr.Type }

// ##==================================================================
type AttrBool struct {
	Key   string
	Value string
	Type  string
}

func NewAttrBool(key string, value string) (*AttrEmpty, error) {
	attr := &AttrEmpty{
		Key:   key,
		Value: value,
		Type:  KeyAttrBool,
	}
	return attr, nil
}

func (attr *AttrBool) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrBool) GetKey() string                    { return attr.Key }
func (attr *AttrBool) GetValue() string                  { return attr.Value }
func (attr *AttrBool) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrBool) GetType() string                   { return attr.Type }

// ##==================================================================
type AttrInt struct {
	Key   string
	Value string
	Type  string
}

func NewAttrInt(key string, value string) (*AttrEmpty, error) {
	attr := &AttrEmpty{
		Key:   key,
		Value: value,
		Type:  KeyAttrInt,
	}
	return attr, nil
}

func (attr *AttrInt) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *AttrInt) GetKey() string                    { return attr.Key }
func (attr *AttrInt) GetValue() string                  { return attr.Value }
func (attr *AttrInt) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *AttrInt) GetType() string                   { return attr.Type }
