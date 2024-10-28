package transpile

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// transpiles all the components in a dir to a specific .go output file
func CompToGo(rootDir string, outputFile string) ([]*Comp, error) {

	comps, err := getComps(rootDir)
	if err != nil {
		return nil, err
	}

	for _, comp := range comps {

		setCompLines(comp)
		filterEmptyLines(comp)
		setStartingSpaces(comp)
		squeezeCompLines(comp)

		for _, line := range comp.Lines {
			fmt.Println(line.Squeezed)
		}
	}

	return nil, nil
}

// a component to be transpiled
type Comp struct {
	Lines    []*CompLine
	Html     string
	Squeezed string
}

// a line in a component
type CompLine struct {
	Text          string
	StartingSpace int
	Squeezed      string
}

// gets all the components from a dir
func getComps(dir string) ([]*Comp, error) {
	var comps []*Comp
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		str := string(f)
		comp := &Comp{
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

// collects all the CompLines from a Comp
func setCompLines(comp *Comp) {
	lines := strings.Split(comp.Html, "\n")
	for _, line := range lines {
		compLine := &CompLine{
			Text: line,
		}
		comp.Lines = append(comp.Lines, compLine)
	}
}

// remove all the empty lines from a Comp
func filterEmptyLines(comp *Comp) {
	var coll []*CompLine
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

// set StartingSpaces in our CompLines
func setStartingSpaces(comp *Comp) {
	for _, line := range comp.Lines {
		index := strings.Index(line.Text, "<")
		line.StartingSpace = index
	}
}

// sets the "Squeezed" value in all CompLine
func squeezeCompLines(comp *Comp) {
	for _, line := range comp.Lines {
		sq := strings.ReplaceAll(line.Text, " ", "")
		line.Squeezed = sq
	}
}

// sets the "Squeezed" value in all Comp
func squeezeComp(comp *Comp) {
	comp.Squeezed = strings.ReplaceAll(comp.Html, " ", "")
}
