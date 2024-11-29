package attr

import "fmt"

type Str struct {
	Key   string
	Value string
	Type  string
}

func NewAttrStr(key string, value string) (*Str, error) {
	attr := &Str{
		Key:   key,
		Value: value,
		Type:  KeyAttrStr,
	}
	return attr, nil
}

func (attr *Str) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *Str) GetKey() string                    { return attr.Key }
func (attr *Str) GetValue() string                  { return attr.Value }
func (attr *Str) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *Str) GetType() string                   { return attr.Type }
