package attr

import "fmt"

type Empty struct {
	Key   string
	Value string
	Type  string
}

func NewEmpty(key string, value string) (*Empty, error) {
	attr := &Empty{
		Key:   key,
		Value: value,
		Type:  KeyAttrEmpty,
	}
	return attr, nil
}

func (attr *Empty) Print()                            { fmt.Println(attr.Key + ":" + attr.Value) }
func (attr *Empty) GetKey() string                    { return attr.Key }
func (attr *Empty) GetValue() string                  { return attr.Value }
func (attr *Empty) GetKeyValuePair() (string, string) { return attr.Key, attr.Value }
func (attr *Empty) GetType() string                   { return attr.Type }
