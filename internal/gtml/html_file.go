package gtml

import (
	"gtml/internal/fungi"
	"gtml/internal/gqpp"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type HtmlFile struct {
	Path       string
	Text       string
	Selection  *goquery.Selection
	Components []Component
}

func NewHtmlFileFromPath(path string) (HtmlFile, error) {
	file := &HtmlFile{
		Path: path,
	}
	err := fungi.Process(
		func() error { return file.setText() },
		func() error { return file.setSelection() },
		func() error { return file.setComponents() },
	)
	if err != nil {
		return *file, err
	}
	return *file, nil
}

func (file *HtmlFile) setText() error {
	fBytes, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}
	fStr := string(fBytes)
	file.Text = fStr
	return nil
}

func (file *HtmlFile) setSelection() error {
	sel, err := gqpp.NewSelectionFromFilePath(file.Path)
	if err != nil {
		return err
	}
	file.Selection = sel
	return nil
}

func (file *HtmlFile) setComponents() error {
	comps := make([]Component, 0)
	var potErr error
	file.Selection.Find("*").Each(func(i int, sel *goquery.Selection) {
		_, exists := sel.Attr("_name")
		if exists {
			comp, err := NewComponentFromSelection(sel)
			if err != nil {
				potErr = err
				return
			}
			comps = append(comps, comp)
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}
