package attr

import "fmt"

type Int struct {
	Key   string
	Value string
	Type  string
}

func NewInt(key string, value string) (*Int, error) {
	attr := &Int{
		Key:   key,
		Value: value,
		Type:  KeyAttrInt,
	}
	return attr, nil
}

func (attr *Int) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *Int) GetKey() string                    { return attr.Key }
func (attr *Int) GetValue() string                  { return attr.Value }
func (attr *Int) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *Int) GetType() string                   { return attr.Type }
