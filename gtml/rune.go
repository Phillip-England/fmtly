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
	KeyRuneSlot = "$slot"
	KeyRuneVal  = "$val"
	KeyRunePipe = "$pipe"
)

const (
	KeyLocationAttribute = "KEYLOCATIONATTRIBUTE"
	KeyLocationElsewhere = "KEYLOCATIONELSEWHERE"
)

// ##==================================================================
func GetRuneNames() []string {
	return []string{KeyRuneProp, KeyRuneSlot, KeyRuneVal, KeyRunePipe}
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
	if strings.HasPrefix(runeStr, KeyRuneSlot) {
		r, err := NewRuneSlot(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(runeStr, KeyRuneVal) {
		r, err := NewRuneVal(runeStr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if strings.HasPrefix(runeStr, KeyRunePipe) {
		r, err := NewRunePipe(runeStr)
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
type RuneSlot struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
}

func NewRuneSlot(data string) (*RuneSlot, error) {
	r := &RuneSlot{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRuneSlot,
	}
	err := fungi.Process(
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RuneSlot) Print()                 { fmt.Println(r.Data) }
func (r *RuneSlot) GetValue() string       { return r.Value }
func (r *RuneSlot) GetType() string        { return r.Type }
func (r *RuneSlot) GetDecodedData() string { return r.DecodedData }
func (r *RuneSlot) GetLocation() string    { return r.Location }

func (r *RuneSlot) initValue() error {
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
	if valFirstChar == "\"" && valLastChar == "\"" {
		valIsDoubleQuotes = true
	}
	if valFirstChar == "'" && valLastChar == "'" {
		valIsSingleQuotes = true
	}
	if !valIsDoubleQuotes && !valIsSingleQuotes {
		msg := purse.Fmt(`
invalid $slot rune found: %s
$slot must contain a single string wrapped in quotes such as $slot("varName")
$slot may only contain characters; no symbols, numbers, or spaces
	`, r.Data)
		return fmt.Errorf(msg)
	}

	if valIsDoubleQuotes {
		if strings.Count(val, "\"") > 2 {
			msg := purse.Fmt(`
invalid $slot rune found: %s
$slot must contain a single string wrapped in quotes such as $slot("varName")
$slot may only contain characters; no symbols, numbers, or spaces
			`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	if valIsSingleQuotes {
		if strings.Count(val, "'") > 2 {
			msg := purse.Fmt(`
invalid $slot rune found: %s
$slot must contain a single string wrapped in quotes such as $slot("varName")
$slot may only contain characters; no symbols, numbers, or spaces
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
invalid $slot rune found: %s
$slot must contain a single string wrapped in quotes such as $slot("varName")
$slot may only contain characters; no symbols, numbers, or spaces
`, r.Data)
		return fmt.Errorf(msg)
	}

	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	r.Value = purse.Squeeze(val)
	return nil
}

// ##==================================================================
type RuneVal struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
}

func NewRuneVal(data string) (*RuneVal, error) {
	r := &RuneVal{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRuneVal,
	}
	err := fungi.Process(
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RuneVal) Print()                 { fmt.Println(r.Data) }
func (r *RuneVal) GetValue() string       { return r.Value }
func (r *RuneVal) GetType() string        { return r.Type }
func (r *RuneVal) GetDecodedData() string { return r.DecodedData }
func (r *RuneVal) GetLocation() string    { return r.Location }

func (r *RuneVal) initValue() error {
	index := strings.Index(r.Data, "(") + 1
	part := r.Data[index:]
	lastChar := string(part[len(part)-1])
	if lastChar != ")" {
		msg := purse.Fmt(`
			invalid $val rune found: %s`, r.Data)
		return fmt.Errorf(msg)
	}
	val := part[:len(part)-1]

	valFirstChar := string(val[0])
	valLastChar := string(val[len(val)-1])
	valIsSingleQuotes := false
	valIsDoubleQuotes := false
	if valFirstChar == "\"" && valLastChar == "\"" {
		valIsDoubleQuotes = true
	}
	if valFirstChar == "'" && valLastChar == "'" {
		valIsSingleQuotes = true
	}

	if valIsDoubleQuotes || valIsSingleQuotes {
		msg := purse.Fmt(`
invalid $val rune found: %s
$val must contain a single value (not a string) such as $val(someValue)
$val may only contain characters; no symbols, numbers, or spaces`, r.Data)
		return fmt.Errorf(msg)
	}

	if valIsDoubleQuotes {
		if strings.Count(val, "\"") > 2 {
			msg := purse.Fmt(`
invalid $val rune found: %s
$val must contain a single value (not a string) such as $val(someValue)
$val may only contain characters; no symbols, numbers, or spaces`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	if valIsSingleQuotes {
		if strings.Count(val, "'") > 2 {
			msg := purse.Fmt(`
invalid $val rune found: %s
$val must contain a single value (not a string) such as $val(someValue)
$val may only contain characters; no symbols, numbers, or spaces`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	whitelist := purse.GetAllLetters()
	whitelist = append(whitelist, ".")
	if !purse.EnforeWhitelist(val, whitelist) {
		msg := purse.Fmt(`
invalid $val rune found: %s
$val must contain a single value (not a string) such as $val(someValue)
$val may only contain characters; no symbols, numbers, or spaces`, r.Data)
		return fmt.Errorf(msg)
	}

	r.Value = purse.Squeeze(val)
	return nil
}

// ##==================================================================
type RunePipe struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
}

func NewRunePipe(data string) (*RunePipe, error) {
	r := &RunePipe{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRunePipe,
	}
	err := fungi.Process(
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RunePipe) Print()                 { fmt.Println(r.Data) }
func (r *RunePipe) GetValue() string       { return r.Value }
func (r *RunePipe) GetType() string        { return r.Type }
func (r *RunePipe) GetDecodedData() string { return r.DecodedData }
func (r *RunePipe) GetLocation() string    { return r.Location }

func (r *RunePipe) initValue() error {
	index := strings.Index(r.Data, "(") + 1
	part := r.Data[index:]
	lastChar := string(part[len(part)-1])
	if lastChar != ")" {
		msg := purse.Fmt(`
invalid $pipe rune found: %s`, r.Data)
		return fmt.Errorf(msg)
	}
	val := part[:len(part)-1]

	valFirstChar := string(val[0])
	valLastChar := string(val[len(val)-1])
	valIsSingleQuotes := false
	valIsDoubleQuotes := false
	if valFirstChar == "\"" && valLastChar == "\"" {
		valIsDoubleQuotes = true
	}
	if valFirstChar == "'" && valLastChar == "'" {
		valIsSingleQuotes = true
	}

	if valIsDoubleQuotes || valIsSingleQuotes {
		msg := purse.Fmt(`
invalid $pipe rune found: %s
$pipe must contain a single value (not a string) such as $pipe(someValue)
$pipe may only contain characters; no symbols, numbers, or spaces`, r.Data)
		return fmt.Errorf(msg)
	}

	if valIsDoubleQuotes {
		if strings.Count(val, "\"") > 2 {
			msg := purse.Fmt(`
invalid $pipe rune found: %s
$pipe must contain a single value (not a string) such as $pipe(someValue)
$pipe may only contain characters; no symbols, numbers, or spaces`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	if valIsSingleQuotes {
		if strings.Count(val, "'") > 2 {
			msg := purse.Fmt(`
invalid $pipe rune found: %s
$pipe must contain a single value (not a string) such as $pipe(someValue)
$pipe may only contain characters; no symbols, numbers, or spaces`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	whitelist := purse.GetAllLetters()
	if !purse.EnforeWhitelist(val, whitelist) {
		msg := purse.Fmt(`
invalid $pipe rune found: %s
$pipe must contain a single value (not a string) such as $pipe(someValue)
$pipe may only contain characters; no symbols, numbers, or spaces`, r.Data)
		return fmt.Errorf(msg)
	}

	r.Value = purse.Squeeze(val)
	return nil
}

// ##==================================================================

// ##==================================================================
