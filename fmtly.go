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
type HtmlComponent struct {
	Html               string
	Node               *goquery.Selection
	Name               string
	Output             string
	BraceTokenDetails  []*BraceTokenDetails
	BracePropTokens    []*BracePropToken
	BraceForPropTokens []*BraceForPropToken
	ForTokens          []*ForToken
}

// creates a new Component
func NewHtmlComponent(compStr string, node *goquery.Selection) (*HtmlComponent, error) {
	comp := &HtmlComponent{
		Html:   compStr,
		Node:   node,
		Output: "func NAME(PROPS) (string) {\n\tBODY\n\tRETURN\n}",
	}
	err := comp.Init()
	if err != nil {
		return nil, err
	}
	return comp, nil
}

// init the components proper state
func (comp *HtmlComponent) Init() error {
	comp.RemoveEmptyLines()
	err := comp.SetName()
	if err != nil {
		return err
	}
	comp.SetDefineTagName()
	err = comp.CollectBasicTokenDetails()
	if err != nil {
		return err
	}
	err = comp.SortPropTokens()
	if err != nil {
		return err
	}
	err = comp.CollectForTokens()
	if err != nil {
		return err
	}
	comp.SetForLoops()
	// last step
	comp.ReadyToReturn()
	return nil
}

// sets the name of the comp using go query
func (comp *HtmlComponent) SetName() error {
	name, nameExists := comp.Node.Attr("name")
	if !nameExists {
		return fmt.Errorf("HtmlComponent does not contain name: %s", comp.Html)
	}
	comp.Output = strings.Replace(comp.Output, "NAME", name, 1)
	return nil
}

// cleans comp.Html of any empty lines
func (comp *HtmlComponent) RemoveEmptyLines() {
	var collect []string
	lines := strings.Split(comp.Html, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		collect = append(collect, line)
	}
	comp.Html = strings.Join(collect, "\n")
}

// replaces <define> using the tag="tagname" syntax
func (comp *HtmlComponent) SetDefineTagName() {
	tagName, tagNameExists := comp.Node.Attr("tag")
	if !tagNameExists {
		tagName = "div"
	}
	var collect []string
	lines := strings.Split(comp.Html, "\n")
	for i, line := range lines {
		if i == 0 || i == len(lines)-1 {
			line = strings.Replace(line, "define", tagName, 1)
		}
		collect = append(collect, line)
	}
	comp.Html = strings.Join(collect, "\n")
}

// collects just Prop and TypeProp tokens
func (comp *HtmlComponent) CollectBasicTokenDetails() error {
	col := collectUsingMarkers(comp.Html, "{{", "}}")
	for _, tok := range col {
		sq := squeezeStr(tok)
		det := &BraceTokenDetails{
			Text:     tok,
			Squeezed: sq,
			Value:    strings.Replace(strings.Replace(sq, "}}", "", 1), "{{", "", 1),
		}
		comp.BraceTokenDetails = append(comp.BraceTokenDetails, det)
	}
	return nil
}

// collects all the for tokens from comp.Html
func (comp *HtmlComponent) CollectForTokens() error {
	var potentialErr error
	filterLines(getLines(comp.Html), func(i int, line string) bool {
		sq := squeezeStr(line)
		if strContainsAll(sq, "<for", "in=", ">") && !strContainsAll(sq, "<form") {
			doc, err := getDoc(line)
			if err != nil {
				potentialErr = err
			}
			doc.Find("for").Each(func(i int, s *goquery.Selection) {
				propName, _ := s.Attr("in")
				goType, _ := s.Attr("type")
				tagName, _ := s.Attr("tag")
				tok, err := NewForToken(propName, goType, tagName, line)
				if err != nil {
					potentialErr = err
				}
				comp.ForTokens = append(comp.ForTokens, tok)
			})

		}
		return false
	})
	if potentialErr != nil {
		return potentialErr
	}
	return nil
}

// sorts the tokens into their specific type buckets
func (comp *HtmlComponent) SortPropTokens() error {
	for _, det := range comp.BraceTokenDetails {
		if strings.Count(det.Text, ".") == 1 && det.Squeezed[0:1] == "{{" {
			forPropToken, err := NewBraceForPropToken(det)
			if err != nil {
				return err
			}
			comp.BraceForPropTokens = append(comp.BraceForPropTokens, forPropToken)
			continue
		}
		propToken, err := NewPropToken(det)
		if err != nil {
			return err
		}
		comp.BracePropTokens = append(comp.BracePropTokens, propToken)
	}
	return nil
}

// sets the for loops in the GoOutput and overwrites it
func (comp *HtmlComponent) SetForLoops() {
	for _, tok := range comp.ForTokens {
		comp.Html = strings.Replace(comp.Html, tok.Text, tok.GoOutput, 1)
	}
	fmt.Println(comp.Html)
}

// places comp.Html into the return of the comp.GoOutput
func (comp *HtmlComponent) ReadyToReturn() {
	indentedCompStr := prefixByLine(comp.Html, "\t\t")
	comp.Output = replaceOne(comp.Output, "RETURN", fmt.Sprintf("return `\n%s\n\t`", indentedCompStr), false)
}

//===============================
// TOKENS
//===============================

// an interface which all tokens share
type BraceToken interface{}

// details which concern every token
type BraceTokenDetails struct {
	Text     string
	Squeezed string
	Value    string
}

//===============================
// FORPROP
//===============================

// a token to represent data pulled from a Go Type
type BraceForPropToken struct {
	Details *BraceTokenDetails
}

// creates a new PropTypeToken
func NewBraceForPropToken(details *BraceTokenDetails) (*BraceForPropToken, error) {
	tok := &BraceForPropToken{
		Details: details,
	}
	return tok, nil
}

//===============================
// PROP
//===============================

// a token to represent a simple func prop
type BracePropToken struct {
	Details *BraceTokenDetails
}

// creates a new PropTypeToken
func NewPropToken(details *BraceTokenDetails) (*BracePropToken, error) {
	tok := &BracePropToken{
		Details: details,
	}
	return tok, nil
}

//===============================
// FOR
//===============================

type ForToken struct {
	PropName   string
	GoType     string
	TagName    string
	Text       string
	GoOutput   string
	RenamedTag string
}

func NewForToken(propName string, goType string, tagName string, text string) (*ForToken, error) {
	tok := &ForToken{
		PropName: propName,
		GoType:   goType,
		TagName:  tagName,
		Text:     text,
	}
	err := tok.Init()
	if err != nil {
		return nil, err
	}
	return tok, nil
}

func (tok *ForToken) Init() error {
	tok.SetRenamedTag()
	tok.SetGoOutput()
	return nil
}

func (tok *ForToken) SetRenamedTag() {
	output := tok.Text
	output = strings.Replace(output, "for", tok.TagName, 1)
	tok.RenamedTag = output
}

func (tok *ForToken) SetGoOutput() {
	tabs := ""
	for i := 0; i < strings.Count(tok.Text, "\t"); i++ {
		tabs += "\t"
	}
	output := fmt.Sprintf("%sMapAny(MakeMap(%s,`\n%s%s", tabs, tok.PropName, tabs, tok.RenamedTag)
	tok.GoOutput = output
}

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

// goes through a string and collects all instances of substrings between provided markers
func collectUsingMarkers(str string, start string, end string) []string {
	var col []string
	for i, c := range str {
		if i+len(start) > len(str)-1 {
			break
		}
		ch := string(c)
		startSearch := ch + string(str[i+len(start)-1])
		if start == startSearch {
			colStr := ""
			for i2, c2 := range str {
				if i2 >= i {
					if i2+len(end) > len(str)-1 {
						colStr = ""
						break
					}
					ch2 := string(c2)
					endSearch := ch2 + string(str[i2+len(end)-1])
					colStr += ch2
					if endSearch == end {
						colStr += string(str[i2+1])
						col = append(col, colStr)
						colStr = ""
						break
					}
				}
			}
		}

	}
	return col
}

// removes all the empty spaces from a string
func squeezeStr(str string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, " ", ""), "\n", ""), "\t", "")
}

// if any line in a string is blank space, this func will remove them
func removeEmptyLines(str string) string {
	var col []string
	for _, line := range strings.Split(str, "\n") {
		if len(line) != 0 {
			col = append(col, line)
		}
	}
	return strings.Join(col, "\n")
}

// determines if a string contains all of the provided chars
func strContainsAll(str string, subStrings ...string) bool {
	for _, ss := range subStrings {
		if !strings.Contains(str, ss) {
			return false
		}
	}
	return true
}

// gets you a slice of lines from a str
func getLines(str string) []string {
	return strings.Split(str, "\n")
}

// takes a slice of strings and combines them with new lines
func joinLines(lines []string) string {
	return strings.Join(lines, "\n")
}

// iterate through the lines of a string, returning them based on a condition
func filterLines(lines []string, fn func(i int, line string) bool) []string {
	var col []string
	for i, line := range lines {
		shouldCollect := fn(i, line)
		if shouldCollect {
			col = append(col, line)
		}
	}
	return col
}

// iterate through the characters in a string, returning them based off of a condition
func filterChars(str string, fn func(i int, char string) bool) string {
	output := ""
	for i, ch := range str {
		char := string(ch)
		shouldCollect := fn(i, char)
		if shouldCollect {
			output += char
		}
	}
	return output
}

// given an index and a string, will return a chunck string[i+int], if able, "" if not
func getLoopChunck(i int, size int, str string) string {
	if i+size > len(str)-1 {
		return ""
	}
	return str[i : i+size]
}

// CountLeadingTabs counts the number of tabs at the beginning of a given line.
func countLeadingTabs(line string) int {
	count := 0
	for _, char := range line {
		if char == '\t' {
			count++
		} else {
			break
		}
	}
	return count
}
