package gtmlrune

import (
	"fmt"
	"html"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type Pipe struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
}

func NewPipe(data string) (*Pipe, error) {
	r := &Pipe{
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

func (r *Pipe) Print()                 { fmt.Println(r.Data) }
func (r *Pipe) GetValue() string       { return r.Value }
func (r *Pipe) GetType() string        { return r.Type }
func (r *Pipe) GetDecodedData() string { return r.DecodedData }
func (r *Pipe) GetLocation() string    { return r.Location }

func (r *Pipe) initValue() error {
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
