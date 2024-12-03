package gtmlrune

import (
	"fmt"
	"html"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

type Md struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
}

func NewMd(data string) (*Md, error) {
	r := &Md{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRuneMd,
	}
	err := fungi.Process(
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Md) Print()                 { fmt.Println(r.Data) }
func (r *Md) GetValue() string       { return r.Value }
func (r *Md) GetType() string        { return r.Type }
func (r *Md) GetDecodedData() string { return r.DecodedData }
func (r *Md) GetLocation() string    { return r.Location }

func (r *Md) initValue() error {
	index := strings.Index(r.Data, "(") + 1
	part := r.Data[index:]
	lastChar := string(part[len(part)-1])
	if lastChar != ")" {
		msg := purse.Fmt(`
invalid $md rune found: %s`, r.Data)
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
invalid $md rune found: %s
$md must containing a single string pointing to a .md file $md("./some/file.md")
	`, r.Data)
		return fmt.Errorf(msg)
	}

	if valIsDoubleQuotes {
		if strings.Count(val, "\"") > 2 {
			msg := purse.Fmt(`
invalid $md rune found: %s
$md must containing a single string pointing to a .md file $md("./some/file.md")
			`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	if valIsSingleQuotes {
		if strings.Count(val, "'") > 2 {
			msg := purse.Fmt(`
invalid $md rune found: %s
$md must containing a single string pointing to a .md file $md("./some/file.md")
			`, r.Data)
			return fmt.Errorf(msg)
		}
	}

	// whitelist := purse.GetAllLetters()
	// if valIsSingleQuotes {
	// 	whitelist = append(whitelist, "\"")
	// }
	// if valIsDoubleQuotes {
	// 	whitelist = append(whitelist, "'")
	// }
	// 	if !purse.EnforeWhitelist(val, whitelist) {
	// 		msg := purse.Fmt(`
	// invalid $md rune found: %s
	// $md must contain a single string wrapped in quotes such as $md("# Markdown Content")
	// $md may only contain characters; no symbols, numbers, or spaces
	// `, r.Data)
	// 		return fmt.Errorf(msg)
	// 	}

	val = strings.ReplaceAll(val, "\"", "")
	val = strings.ReplaceAll(val, "'", "")
	// here we have the md content, just convert it and write it into the HTML :)
	r.Value = val
	return nil
}
