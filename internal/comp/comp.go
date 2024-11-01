package comp

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReadDir(targetDir string) error {
	str, err := readDirToStr(targetDir)
	if err != nil {
		return err
	}
	strComps := filterStrComps(str)
	docs, err := makeQueryDocs(strComps)
	if err != nil {
		return err
	}
	comps, err := makeRawComps(docs)
	if err != nil {
		return err
	}
	for _, comp := range comps {
		if err := setCompHtml(comp); err != nil {
			return err
		}
		if err := setCompName(comp); err != nil {
			return err
		}
		if err := setRawForTags(comp); err != nil {
			return err
		}
		for _, tag := range comp.ForTags {
			if err := setForTagHtml(tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func readDirToStr(targetDir string) (string, error) {
	var ss []string
	err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		fStr := string(f)
		ss = append(ss, fStr)
		return nil
	})
	if err != nil {
		return "", err
	}
	return strings.Join(ss, "\n"), nil
}

func filterStrComps(str string) []string {
	var slc []string
	inComp := false
	lines := strings.Split(str, "\n")
	var comp []string
	for _, line := range lines {
		sq := strings.ReplaceAll(line, " ", "")
		if strings.Contains(sq, "<fmt") {
			inComp = true
		}
		if inComp {
			comp = append(comp, line)
		}
		if strings.Contains(sq, "</fmt") {
			slc = append(slc, strings.Join(comp, "\n"))
			comp = make([]string, 0)
			inComp = false
		}
	}
	return slc
}

func makeQueryDocs(comps []string) ([]*goquery.Document, error) {
	var docs []*goquery.Document
	for _, comp := range comps {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(comp))
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func makeRawComps(docs []*goquery.Document) ([]*Comp, error) {
	var comps []*Comp
	for _, doc := range docs {
		comps = append(comps, &Comp{
			Doc: doc,
		})
	}
	return comps, nil
}

func setCompHtml(comp *Comp) error {
	str, err := comp.Doc.Find("body").Html()
	if err != nil {
		return err
	}
	comp.Html = str
	return nil
}

func setCompName(comp *Comp) error {
	tag := comp.Doc.Find("fmt")
	name, exists := tag.Attr("name")
	if !exists {
		return fmt.Errorf("<fmt> tag does not contain a name:\n\n%s", comp.Html)
	}
	comp.Name = name
	return nil
}

func setRawForTags(comp *Comp) error {
	comp.Doc.Find("for").Each(func(i int, s *goquery.Selection) {
		parentCount := 0
		s.Parents().Each(func(i2 int, s2 *goquery.Selection) {
			tagName := goquery.NodeName(s2)
			if tagName == "for" {
				parentCount++
			}
		})
		comp.ForTags = append(comp.ForTags, &ForTag{
			Selection: s,
			Depth:     parentCount,
		})
	})
	return nil
}

func setForTagHtml(tag *ForTag) error {
	htmlStr, err := goquery.OuterHtml(tag.Selection)
	if err != nil {
		return err
	}
	lines := strings.Split(htmlStr, "\n")
	var slc []string
	for _, line := range lines {
		col := ""
		foundChar := false
		for _, ch := range line {
			char := string(ch)
			if char == " " {
				if !foundChar {
					continue
				}
			} else {
				foundChar = true
			}
			if foundChar {
				col += char
				continue
			}
		}
		slc = append(slc, col)
	}
	oneLineHtml := strings.Join(slc, "")
	tag.Html = oneLineHtml
	fmt.Println(tag.Html)
	return nil
}
