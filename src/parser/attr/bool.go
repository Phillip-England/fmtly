package attr

import "fmt"

type Bool struct {
	Key   string
	Value string
	Type  string
}

func NewBool(key string, value string) (*Bool, error) {
	attr := &Bool{
		Key:   key,
		Value: value,
		Type:  KeyAttrBool,
	}
	return attr, nil
}

func (attr *Bool) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *Bool) GetKey() string                    { return attr.Key }
func (attr *Bool) GetValue() string                  { return attr.Value }
func (attr *Bool) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *Bool) GetType() string                   { return attr.Type }
