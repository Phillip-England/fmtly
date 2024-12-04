package gtmlrune

import (
	"fmt"
	"gtml/src/parser/funcarg"
	"html"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type Val struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
	Args        []funcarg.FuncArg
}

func NewVal(data string) (*Val, error) {
	r := &Val{
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

func (r *Val) Print()                     { fmt.Println(r.Data) }
func (r *Val) GetValue() string           { return r.Value }
func (r *Val) GetType() string            { return r.Type }
func (r *Val) GetDecodedData() string     { return r.DecodedData }
func (r *Val) GetLocation() string        { return r.Location }
func (r *Val) GetArgs() []funcarg.FuncArg { return r.Args }

func (r *Val) initValue() error {
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
