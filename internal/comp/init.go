package comp

import (
	"fmt"
	"fmtly/internal/parsley"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReadDir(targetDir string) ([]*Comp, error) {
	str, err := readDirToStr(targetDir)
	if err != nil {
		return nil, err
	}
	strComps := filterStrComps(str)
	docs, err := makeQueryDocs(strComps)
	if err != nil {
		return nil, err
	}
	comps, err := makeRawComps(docs)
	if err != nil {
		return nil, err
	}
	// general comp setup
	for _, comp := range comps {
		if err := setCompHtml(comp); err != nil {
			return nil, err
		}
		if err := setCompName(comp); err != nil {
			return nil, err
		}
		// for tag setup
		if err := setRawForTags(comp); err != nil {
			return nil, err
		}
		for _, tag := range comp.ForTags {
			if err := setForTagAttrs(tag); err != nil {
				return nil, err
			}
			if err := setForTagHtml(tag); err != nil {
				return nil, err
			}
			if err := setForTagAttrStr(tag); err != nil {
				return nil, err
			}
		}
		// fmt tag setup
		if err := setRawFmtTag(comp); err != nil {
			return nil, err
		}
		if err := setFmtTagAttrs(comp.FmtTag); err != nil {
			return nil, err
		}
		if err := setFmtTagHtml(comp.FmtTag); err != nil {
			return nil, err
		}
		if err := setFmtTagAttrStr(comp.FmtTag); err != nil {
			return nil, err
		}
		// if tag setup
		if err := setRawIfTags(comp); err != nil {
			return nil, err
		}
		for _, tag := range comp.IfTags {
			if err := setIfTagAttrs(tag); err != nil {
				return nil, err
			}
			if err := setIfTagHtml(tag); err != nil {
				return nil, err
			}
			if err := setIfElseTag(tag); err != nil {
				return nil, err
			}
		}
		if err := collectProps(comp); err != nil {
			return nil, err
		}

	}
	return comps, nil
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
	str, err := GetHtmlFromSelection(tag.Selection)
	if err != nil {
		return err
	}
	tag.Html = str
	return nil
}

func setForTagAttrs(tag *ForTag) error {
	inAttr, _ := tag.Selection.Attr("in")
	typeAttr, _ := tag.Selection.Attr("type")
	tagAttr, _ := tag.Selection.Attr("tag")
	asAttr, _ := tag.Selection.Attr("as")
	tag.AttrIn = inAttr
	tag.AttrTag = tagAttr
	tag.AttrType = typeAttr
	tag.AttrAs = asAttr
	return nil
}

func setForTagAttrStr(tag *ForTag) error {
	if node := tag.Selection.Get(0); node != nil {
		for _, attr := range node.Attr {
			if parsley.EqualsOneof(attr.Key, "in", "tag", "as", "type") {
				continue
			}
			tag.AttrStr = tag.AttrStr + attr.Key + "=\"" + attr.Val + "\" "
		}
	}
	return nil
}

func setRawFmtTag(comp *Comp) error {
	s := comp.Doc.Find("fmt")
	if s == nil {
		return fmt.Errorf("malformed <fmt> component:\n\n%s", comp.Html)
	}
	comp.FmtTag = &FmtTag{
		Selection: s,
	}
	return nil
}

func setFmtTagAttrs(tag *FmtTag) error {
	nameAttr, _ := tag.Selection.Attr("name")
	tagAttr, _ := tag.Selection.Attr("tag")
	tag.AttrName = nameAttr
	tag.AttrTag = tagAttr
	return nil
}

func setFmtTagHtml(tag *FmtTag) error {
	str, err := GetHtmlFromSelection(tag.Selection)
	if err != nil {
		return err
	}
	tag.Html = str
	return nil
}

func setFmtTagAttrStr(tag *FmtTag) error {
	if node := tag.Selection.Get(0); node != nil {
		for _, attr := range node.Attr {
			if parsley.EqualsOneof(attr.Key, "name", "tag") {
				continue
			}
			tag.AttrStr = tag.AttrStr + attr.Key + "=\"" + attr.Val + "\" "
		}
	}
	return nil
}

func setRawIfTags(comp *Comp) error {
	comp.Doc.Find("if").Each(func(i int, s *goquery.Selection) {
		parentCount := 0
		s.Parents().Each(func(i2 int, s2 *goquery.Selection) {
			tagName := goquery.NodeName(s2)
			if tagName == "if" {
				parentCount++
			}
		})
		comp.IfTags = append(comp.IfTags, &IfTag{
			Selection: s,
			Depth:     parentCount,
		})
	})
	return nil
}

func setIfTagHtml(tag *IfTag) error {
	str, err := GetHtmlFromSelection(tag.Selection)
	if err != nil {
		return err
	}
	tag.Html = str
	return nil
}

func setIfTagAttrs(tag *IfTag) error {
	conditionAttt, _ := tag.Selection.Attr("condition")
	tagAttr, _ := tag.Selection.Attr("tag")
	tag.AttrCondition = conditionAttt
	tag.AttrTag = tagAttr
	return nil
}

func setIfElseTag(tag *IfTag) error {
	elseCount := 0
	tag.Selection.Find("else").Each(func(i int, s *goquery.Selection) {
		elseCount++
	})
	if elseCount > 1 {
		return fmt.Errorf("<if> has more than one <else>:\n\n%s", tag.Html)
	}
	elseSelection := tag.Selection.Find("else")
	htmlStr, err := GetHtmlFromSelection(elseSelection)
	if err != nil {
		return err
	}
	tag.ElseTag = &ElseTag{
		Selection: elseSelection,
		Html:      htmlStr,
	}
	return nil
}

func collectProps(comp *Comp) error {
	inProp := false
	prop := ""
	for i, ch := range comp.Html {
		if i+1 > len(comp.Html)-1 {
			continue
		}
		char := string(ch)
		nextChar := string(comp.Html[i+1])
		search := char + nextChar
		if search == "{{" {
			inProp = true
		}
		if inProp {
			prop += char
			if search == "}}" {
				inProp = false
				prop += nextChar
				sq := parsley.Squeeze(prop)
				value := strings.ReplaceAll(sq, "{{", "")
				value = strings.ReplaceAll(value, "}}", "")
				if strings.Contains(value, ".") {
					comp.ForProps = append(comp.ForProps, &ForProp{
						Raw:   prop,
						Value: value,
					})
				} else {
					comp.Props = append(comp.Props, &Prop{
						Raw:   prop,
						Value: value,
					})
				}
				prop = ""
			}

		}
	}
	return nil
}
