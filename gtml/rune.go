package gtml

import (
	"fmt"
	"html"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyRuneProp = "$prop"
)

// ##==================================================================
func GetRuneNames() []string {
	return []string{KeyRuneProp, "$val", "$pipe"}
}

// ##==================================================================
type GtmlRune interface {
	Print()
	GetValue() string
	GetType() string
	GetDecodedData() string
}

func NewGtmlRune(runeStr string) (GtmlRune, error) {
	if strings.HasPrefix(runeStr, KeyRuneProp) {
		r, err := NewRuneProp(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, nil
}

// ##==================================================================
type RuneProp struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
}

func NewRuneProp(data string) (*RuneProp, error) {
	r := &RuneProp{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRuneProp,
	}
	err := fungi.Process(
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RuneProp) Print()                 { fmt.Println(r.Data) }
func (r *RuneProp) GetValue() string       { return r.Value }
func (r *RuneProp) GetType() string        { return r.Type }
func (r *RuneProp) GetDecodedData() string { return r.DecodedData }

func (r *RuneProp) initValue() error {
	index := strings.Index(r.Data, "(") + 1
	part := r.Data[index:]
	lastChar := string(part[len(part)-1])
	if lastChar != ")" {
		msg := purse.Fmt(`
invalid $prop rune found: %s`, r.Data)
		return fmt.Errorf(msg)
	}
	part = part[:len(part)-1]
	if strings.Count(part, "\"") != 2 {
		msg := purse.Fmt(`
invalid $prop rune found: %s
$prop must contain a single string wrapped in double quotes such as $prop("varName")
$prop may only contain characters; no symbols, numbers, or spaces
`, r.Data)
		return fmt.Errorf(msg)
	}
	whitelist := purse.GetAllLetters()
	whitelist = append(whitelist, "\"")
	if !purse.EnforeWhitelist(part, whitelist) {
		msg := purse.Fmt(`
invalid $prop rune found: %s
$prop must contain a single string wrapped in double quotes such as $prop("varName")
$prop may only contain characters; no symbols, numbers, or spaces
`, r.Data)
		return fmt.Errorf(msg)
	}
	val := part
	valFirstChar := string(val[0])
	valLastChar := string(val[len(val)-1])
	if valFirstChar != "\"" && valLastChar != "\"" {
		msg := purse.Fmt(`
invalid $prop rune found: %s
$prop must contain a single string wrapped in double quotes such as $prop("varName")
$prop may only contain characters; no symbols, numbers, or spaces
`, r.Data)
		return fmt.Errorf(msg)
	}
	val = strings.ReplaceAll(val, "\"", "")
	r.Value = purse.Squeeze(val)
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
