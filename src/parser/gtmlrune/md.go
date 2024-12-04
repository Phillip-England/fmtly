package gtmlrune

import (
	"bytes"
	"fmt"
	"gtml/src/parser/funcarg"
	"html"
	"os"
	"strings"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

type Md struct {
	Data        string
	DecodedData string
	Value       string
	Type        string
	Location    string
	Args        []funcarg.FuncArg
}

func NewMd(data string) (*Md, error) {
	r := &Md{
		DecodedData: data,
		Data:        html.UnescapeString(data),
		Type:        KeyRuneMd,
	}
	err := fungi.Process(
		func() error { return r.initSetFuncArgs() },
		func() error { return r.initValue() },
	)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Md) Print()                     { fmt.Println(r.Data) }
func (r *Md) GetValue() string           { return r.Value }
func (r *Md) GetType() string            { return r.Type }
func (r *Md) GetDecodedData() string     { return r.DecodedData }
func (r *Md) GetLocation() string        { return r.Location }
func (r *Md) GetArgs() []funcarg.FuncArg { return r.Args }

func (r *Md) initSetFuncArgs() error {
	data := r.Data
	index := strings.Index(data, "(")
	argStr := data[index+1:]
	argStr = argStr[:len(argStr)-1]
	parts := strings.Split(argStr, ",")
	for _, part := range parts {
		arg, err := funcarg.NewFuncArg(part)
		if err != nil {
			return err
		}
		r.Args = append(r.Args, arg)
	}
	if len(r.Args) != 2 {
		return purse.Err(`
$md rune requires 2 args, a file path and a color theme: 
$md("/path/to/file.md", "dracula")`)
	}
	return nil
}

func (r *Md) initValue() error {
	mdPath := r.Args[0].GetValue()
	theme := r.Args[1].GetValue()
	mdFileContent, err := os.ReadFile(mdPath)
	if err != nil {
		return err
	}
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert([]byte(mdFileContent), &buf); err != nil {
		panic(err)
	}
	r.Value = buf.String()
	return nil
}
