package transpile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// transpiles all the components in a dir to a specific .go output file
func CompToGo(rootDir string, outputFile string) ([]*FmtComp, error) {
	comps, err := getComps(rootDir)
	if err != nil {
		return nil, err
	}
	for _, comp := range comps {
		setLines(comp)
		filterEmptyLines(comp)
		setStartingSpaces(comp)
		squeezeLines(comp)
		squeezeComp(comp)
		setLineNumbers(comp)
		setLineCount(comp)
		setNoIndents(comp)
		setTags(comp)
		stripTags(comp)
		setTagParts(comp)
		setTagNames(comp)
		markClosingTags(comp)
		setTagStrAttrs(comp)
		setTagAttrs(comp)
		setFmtTag(comp)
		setFmtTagAttrs(comp)
		setForTags(comp)
		setForTagAttrs(comp)
		prepareCompOutput(comp)
	}
	return nil, nil
}

type FmtComp struct {
	Lines     []*HtmlLine
	Tags      []*HtmlTag
	Html      string
	Squeezed  string
	LineCount int
	FmtTag    *FmtTag
	ForTags   []*ForTag
	Output    string
}

type HtmlLine struct {
	Text          string
	StartingSpace int
	Squeezed      string
	Number        int
	NoIndent      string
}

type HtmlTag struct {
	Text       string
	Name       string
	Stripped   string
	Parts      []string
	IsCloseTag bool
	StrAttrs   []string
	Attrs      []*HtmlAttr
	LineNumber int
}

type HtmlAttr struct {
	Name  string
	Value string
}

type FmtTag struct {
	Tag      *HtmlTag
	AttrName string
	AttrTag  string
}

type ForTag struct {
	Tag      *HtmlTag
	AttrIn   string
	AttrType string
	AttrTag  string
}

func getComps(dir string) ([]*FmtComp, error) {
	var comps []*FmtComp
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str := string(f)
		comp := &FmtComp{
			Html: str,
		}
		comps = append(comps, comp)
		return nil
	})
	if err != nil {
		return comps, err
	}
	return comps, nil
}

func setLines(comp *FmtComp) {
	lines := strings.Split(comp.Html, "\n")
	for _, line := range lines {
		htmlLine := &HtmlLine{
			Text: line,
		}
		comp.Lines = append(comp.Lines, htmlLine)
	}
}

func filterEmptyLines(comp *FmtComp) {
	var coll []*HtmlLine
	var text []string
	for _, line := range comp.Lines {
		if len(line.Text) == 0 {
			continue
		}
		coll = append(coll, line)
		text = append(text, line.Text)
	}
	comp.Lines = coll
	comp.Html = strings.Join(text, "\n")
}

func setStartingSpaces(comp *FmtComp) {
	for _, line := range comp.Lines {
		index := strings.Index(line.Text, "<")
		line.StartingSpace = index
	}
}

func squeezeLines(comp *FmtComp) {
	for _, line := range comp.Lines {
		sq := strings.ReplaceAll(line.Text, " ", "")
		line.Squeezed = sq
	}
}

func squeezeComp(comp *FmtComp) {
	comp.Squeezed = strings.ReplaceAll(comp.Html, " ", "")
}

func setLineNumbers(comp *FmtComp) {
	for i, line := range comp.Lines {
		line.Number = i
	}
}

func setLineCount(comp *FmtComp) {
	comp.LineCount = len(comp.Lines)
}

func setNoIndents(comp *FmtComp) {
	for _, line := range comp.Lines {
		startI := strings.Index(line.Text, "<")
		line.NoIndent = line.Text[startI:]
	}
}

func setTags(comp *FmtComp) {
	for _, line := range comp.Lines {
		inTag := false
		var tags []string
		tag := ""
		for _, ch := range line.NoIndent {
			char := string(ch)
			if char == "<" {
				inTag = true
			}
			if inTag {
				tag += char
				if char == ">" {
					tags = append(tags, tag)
					inTag = false
					tag = ""
				}
			}
		}
		for _, tg := range tags {
			comp.Tags = append(comp.Tags, &HtmlTag{
				Text:       tg,
				LineNumber: line.Number,
			})
		}
	}
}

func stripTags(comp *FmtComp) {
	for _, tag := range comp.Tags {
		copy := tag.Text
		copy = strings.Replace(copy, "<", "", 1)
		copy = strings.Replace(copy, ">", "", 1)
		tag.Stripped = copy
	}
}

func setTagParts(comp *FmtComp) {
	for _, tag := range comp.Tags {
		tag.Parts = strings.Split(tag.Stripped, " ")
	}
}

func setTagNames(comp *FmtComp) {
	for _, tag := range comp.Tags {
		if len(tag.Parts) != 0 {
			tag.Name = strings.Replace(tag.Parts[0], "/", "", 1)
		}
	}
}

func markClosingTags(comp *FmtComp) {
	for _, tag := range comp.Tags {
		if len(tag.Parts) != 0 {
			first := tag.Parts[0]
			if strings.Contains(first, "/") {
				tag.IsCloseTag = true
			}
		}
	}
}

func setTagStrAttrs(comp *FmtComp) {
	for _, tag := range comp.Tags {
		if tag.IsCloseTag {
			continue
		}
		if len(tag.Parts) >= 1 {
			parts := tag.Parts[1:]
			if len(parts) == 0 {
				continue
			}
			var attrs []string
			for _, attr := range parts {
				lastChar := string(attr[len(attr)-1])
				if strings.Contains(attr, "=\"") && lastChar == "\"" || strings.Contains(attr, "='") && lastChar == "'" {
					attrs = append(attrs, attr)
				}
			}
			tag.StrAttrs = attrs
		}
	}
}

func setTagAttrs(comp *FmtComp) {
	for _, tag := range comp.Tags {
		for _, attr := range tag.StrAttrs {
			parts := strings.Split(attr, "=")
			if len(parts) != 2 {
				continue
			}
			name := parts[0]
			value := parts[1]
			value = value[1 : len(value)-1]
			tag.Attrs = append(tag.Attrs, &HtmlAttr{
				Name:  name,
				Value: value,
			})
		}
	}
}

func setFmtTag(comp *FmtComp) {
	for _, tag := range comp.Tags {
		if tag.IsCloseTag {
			continue
		}
		if comp.FmtTag != nil {
			break
		}
		if tag.Name == "fmt" {
			comp.FmtTag = &FmtTag{
				Tag: tag,
			}
		}
	}
}

func setFmtTagAttrs(comp *FmtComp) {
	for _, attr := range comp.FmtTag.Tag.Attrs {
		if attr.Name == "name" {
			comp.FmtTag.AttrName = attr.Value
		}
		if attr.Name == "tag" {
			comp.FmtTag.AttrTag = attr.Value
		}
	}
}

func setForTags(comp *FmtComp) {
	for _, tag := range comp.Tags {
		if tag.IsCloseTag {
			continue
		}
		if tag.Name == "for" {
			comp.ForTags = append(comp.ForTags, &ForTag{
				Tag: tag,
			})
		}
	}
}

func setForTagAttrs(comp *FmtComp) {
	for _, tag := range comp.ForTags {
		for _, attr := range tag.Tag.Attrs {
			if attr.Name == "in" {
				tag.AttrIn = attr.Value
			}
			if attr.Name == "type" {
				tag.AttrType = attr.Value
			}
			if attr.Name == "tag" {
				tag.AttrTag = attr.Value
			}
		}
	}
}

func prepareCompOutput(comp *FmtComp) {
	comp.Output = fmt.Sprintf("func NAME(PARAM) string {\nreturn `\n%s\n`\n}", comp.Html)
	lines := strings.Split(comp.Output, "\n")
	var col []string
	for i, line := range lines {
		if i == 0 {
			col = append(col, line)
			continue
		}
		if i == 1 {
			col = append(col, "\t"+line)
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "`" {
			col = append(col, "\t"+line)
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "}" {
			col = append(col, line)
			continue
		}
		col = append(col, "\t\t"+line)
	}
	comp.Output = strings.Join(col, "\n")
	fmt.Println(comp.Output)
}
