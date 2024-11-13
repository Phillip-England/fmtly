package gtml

import (
	"fmt"
	"strings"
)

type GoForLoop struct {
	VarName          string
	BuilderName      string
	IterItem         string
	IterType         string
	Props            []Prop
	HtmlInput        string
	WriteStringCalls []string
}

func NewGoForLoop(varName string, iterItem string, iterType string, props []Prop, htmlInput string) (GoForLoop, error) {
	loop := &GoForLoop{
		VarName:     varName + "Loop",
		BuilderName: varName + "Builder",
		IterItem:    iterItem,
		IterType:    iterType,
		Props:       props,
		HtmlInput:   htmlInput,
	}
	htmlStr := loop.HtmlInput
	for _, prop := range props {
		writeStringCall := fmt.Sprintf("%s.WriteSting(%s)", loop.BuilderName, prop.Value)
		loop.WriteStringCalls = append(loop.WriteStringCalls, writeStringCall)
		htmlStr = strings.Replace(htmlStr, prop.Raw, writeStringCall, 1)
	}
	fmt.Println(htmlStr)
	out := fmt.Sprintf("%s := collect(%s, func(i int, %s %s) {\n\tvar %s strings.Builder\n\treturn `RETURN`\n})", loop.VarName, loop.IterItem, loop.VarName, loop.IterType, loop.BuilderName)
	fmt.Println(out)
	return *loop, nil
}
