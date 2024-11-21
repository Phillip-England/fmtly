package gtml

import "fmt"

// ##==================================================================
type Attr interface {
	Print()
	GetKey() string
	GetValue() string
}

func NewAttr(key string, value string) (Attr, error) {
	attr, err := NewAttrPlaceholder(key, value)
	if err != nil {
		return nil, err
	}
	return attr, nil
}

// ##==================================================================
type AttrPlaceholder struct {
	Key   string
	Value string
}

func NewAttrPlaceholder(key string, value string) (*AttrPlaceholder, error) {
	attr := &AttrPlaceholder{
		Key:   key,
		Value: value,
	}
	return attr, nil
}

func (attr *AttrPlaceholder) Print() {
	fmt.Println("Key: " + attr.Key)
	fmt.Println("Value: " + attr.Value)
}
func (attr *AttrPlaceholder) GetKey() string   { return attr.Key }
func (attr *AttrPlaceholder) GetValue() string { return attr.Value }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
