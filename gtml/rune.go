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

const (
	KeyLocationAttribute = "KEYLOCATIONATTRIBUTE"
	KeyLocationElsewhere = "KEYLOCATIONELSEWHERE"
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
	GetLocation() string
}

func NewGtmlRune(runeStr string, location string) (GtmlRune, error) {
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
	Location    string
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
func (r *RuneProp) GetLocation() string    { return r.Location }

func (r *RuneProp) initValue() error {
	index := strings.Index(r.Data, "(") + 1
	part := r.Data[index:]
	lastChar := string(part[len(part)-1])
	if lastChar != ")" {
		msg := purse.Fmt(`
invalid $prop rune found: %s`, r.Data)
		return fmt.Errorf(msg)
	}
	val := part[:len(part)-1]

	valFirstChar := string(val[0])
	valLastChar := string(val[len(val)-1])
	valIsSingleQuotes := false
	valIsDoubleQuotes := false
	if valFirstChar != "\"" && valLastChar != "\"" {
		valIsDoubleQuotes = true
	}
	if valFirstChar != "'" && valLastChar != "'" {
		valIsSingleQuotes = true
	}
	if !valIsDoubleQuotes && !valIsSingleQuotes {
		msg := purse.Fmt(`
	invalid $prop rune found: %s
	$prop must contain a single string wrapped in quotes such as $prop("varName")
	$prop may only contain characters; no symbols, numbers, or spaces
	`, r.Data)
		return fmt.Errorf(msg)
	}

	if valIsDoubleQuotes {
		if strings.Count(val, "\"") > 2 {
			msg := purse.Fmt(`
			invalid $prop rune found: %s
			$prop must contain a single string wrapped in quotes such as $prop("varName")
			$prop may only contain characters; no symbols, numbers, or spaces
			`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	if valIsSingleQuotes {
		if strings.Count(val, "'") > 2 {
			msg := purse.Fmt(`
			invalid $prop rune found: %s
			$prop must contain a single string wrapped in quotes such as $prop("varName")
			$prop may only contain characters; no symbols, numbers, or spaces
			`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	whitelist := purse.GetAllLetters()
	if valIsSingleQuotes {
		whitelist = append(whitelist, "\"")
	}
	if valIsDoubleQuotes {
		whitelist = append(whitelist, "'")
	}
	if !purse.EnforeWhitelist(val, whitelist) {
		msg := purse.Fmt(`
invalid $prop rune found: %s
$prop must contain a single string wrapped in quotes such as $prop("varName")
$prop may only contain characters; no symbols, numbers, or spaces
`, r.Data)
		return fmt.Errorf(msg)
	}

	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	r.Value = purse.Squeeze(val)
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
