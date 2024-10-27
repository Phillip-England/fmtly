package fmtly

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//===============================
// COMPONENTS
//===============================

// an interface which all components share
type Component interface{}

// represents a fmtly component
type HTMLComponent struct {
	HTML     string
	Node     *goquery.Selection
	Name     string
	GoOutput string
	Tokens   []*Token
}

// creates a new Component
func NewHTMLComponent(compStr string, node *goquery.Selection) (*HTMLComponent, error) {
	comp := &HTMLComponent{
		HTML:     compStr,
		Node:     node,
		GoOutput: "func NAME(PROPS) (string) {\n\tBODY\n\tRETURN\n}",
	}
	err := comp.Init()
	if err != nil {
		return nil, err
	}
	return comp, nil
}

// init the components proper state
func (comp *HTMLComponent) Init() error {
	comp.RemoveEmptyLines()
	err := comp.SetName()
	if err != nil {
		return err
	}
	comp.SetDefineTagName()
	comp.ReadyToReturn()
	return nil
}

// sets the name of the comp using go query
func (comp *HTMLComponent) SetName() error {
	name, nameExists := comp.Node.Attr("name")
	if !nameExists {
		return fmt.Errorf("HTMLComponent does not contain name: %s", comp.HTML)
	}
	comp.GoOutput = strings.Replace(comp.GoOutput, "NAME", name, 1)
	return nil
}

// cleans comp.HTML of any empty lines
func (comp *HTMLComponent) RemoveEmptyLines() {
	var collect []string
	lines := strings.Split(comp.HTML, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		collect = append(collect, line)
	}
	comp.HTML = strings.Join(collect, "\n")
}

// replaces <define> using the tag="tagname" syntax
func (comp *HTMLComponent) SetDefineTagName() {
	tagName, tagNameExists := comp.Node.Attr("tag")
	if !tagNameExists {
		tagName = "div"
	}
	var collect []string
	lines := strings.Split(comp.HTML, "\n")
	for i, line := range lines {
		if i == 0 || i == len(lines)-1 {
			line = strings.Replace(line, "define", tagName, 1)
		}
		fmt.Println(line)
		collect = append(collect, line)
	}
	comp.HTML = strings.Join(collect, "\n")
}

// places comp.HTML into the return of the comp.GoOutput
func (comp *HTMLComponent) ReadyToReturn() {
	indentedCompStr := prefixByLine(comp.HTML, "\t\t")
	comp.GoOutput = replaceOne(comp.GoOutput, "RETURN", fmt.Sprintf("return `\n%s\n\t`", indentedCompStr), false)
}

//===============================
// TOKENS
//===============================

// an interface which all tokens share
type Token interface{}

//===============================
// HELPERS
//===============================

// gets all the file text out of a dir
func readDirToStr(root string) (string, error) {
	output := ""
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		output += string(f)
		return nil
	})
	if err != nil {
		return "", err
	}
	return output, nil
}

// get a go query doc ready
func getDoc(str string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// takes in a string, and adds a tab on each line
func prefixByLine(str string, prefix string) string {
	lines := strings.Split(str, "\n")
	var collect []string
	for _, line := range lines {
		line = prefix + line
		collect = append(collect, line)
	}
	return strings.Join(collect, "\n")
}

// replaces a single string with the option to leave behind an indicator for further replacing
func replaceOne(str string, old string, new string, indicate bool) string {
	if indicate {
		return strings.Replace(str, old, new+old, 1)
	}
	return strings.Replace(str, old, new, 1)
}
