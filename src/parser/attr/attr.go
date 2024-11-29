package attr

import (
	"strconv"
	"strings"

	"github.com/phillip-england/purse"
)

type Attr interface {
	Print()
	GetKey() string
	GetValue() string
	GetKeyValuePair() (string, string)
	GetType() string
}

func NewAttr(key string, value string) (Attr, error) {
	if strings.Contains(key, "-") {
		key = purse.KebabToCamelCase(key)
	}
	if value == "" {
		attr, err := NewEmpty(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	sqValue := purse.Squeeze(value)
	if sqValue == "true" || sqValue == "false" {
		attr, err := NewBool(key, value)
		if err != nil {
			return nil, err
		}
		return attr, nil
	}
	_, err := strconv.Atoi(sqValue)
	if err == nil {
		attr, err := NewInt(key, value)
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
