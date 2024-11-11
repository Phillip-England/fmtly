package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"tagly/internal/gqpp"
	"tagly/internal/parsley"
	"tagly/internal/tagly"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	fTagly, err := tagly.NewFileTaglyFromFilePath("./components/index.t.html")
	if err != nil {
		panic(err)
	}

	_, err = tagly.NewFileComponentFromFileTagly(fTagly)
	if err != nil {
		panic(err)
	}

	// taglyFile, err := filetype.NewTaglyFileFromFilePath("./components/index.t.html")
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = filetype.NewTemplateFileFromTaglyFile(taglyFile)
	// if err != nil {
	// 	panic(err)
	// }

	// err := emptyFile("./components.go")
	// if err != nil {
	// 	panic(err)
	// }

	// err = appendFile("./components.go", "package main"+"\n\n")
	// if err != nil {
	// 	panic(err)
	// }

	// str, err := dirToStr("./components")
	// if err != nil {
	// 	panic(err)
	// }

	// fmtTags, err := fmtTagsFromStr(str)
	// if err != nil {
	// 	panic(err)
	// }

	// for i, tag := range fmtTags {

	// 	tag, err = renameFmtTag(tag)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	tag, err = renameTagsBySelector(tag, "for", "if")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	tag, err = setTagDepthByDirectiveAttrs(tag, "for", "if")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	props, err := extractFmtTagProps(tag)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	paramStr, err := makeFmtTagParamStr(tag, props)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	tag, err = outputFmtTagProps(tag, props)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	tag, err = outputInnerTags(tag)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	tag, err = wrapFmtTagInGoFunc(tag, paramStr)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// err = formatOutput("./components.go")
	// if err != nil {
	// 	panic(err)
	// }

	// 	err = appendFile("./components.go", tag+"\n\n")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmtTags[i] = tag

	// }

}

func emptyFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func dirToStr(dir string) (string, error) {
	out := ""
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fStr := string(f)
		out += fStr + "\n"
		return nil
	})
	if err != nil {
		return "", err
	}
	return out, nil
}

func fmtTagsFromStr(str string) ([]string, error) {
	fmtTags := make([]string, 0)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	var potErr error
	potErr = nil
	doc.Find("fmt").Each(func(i int, s *goquery.Selection) {
		fmtTagStr, err := goquery.OuterHtml(s)
		if err != nil {
			potErr = err
			return
		}
		ignoreAttr, _ := s.Attr("ignore")
		if ignoreAttr == "true" {
			return
		}
		fmtTags = append(fmtTags, fmtTagStr)
	})
	if potErr != nil {
		return nil, potErr
	}
	return fmtTags, nil
}

func renameFmtTag(fmtTags string) (string, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(fmtTags)
	if err != nil {
		return "", err
	}
	tagAttr, _ := s.Attr("tag")
	s, err = gqpp.ChangeSelectionTagName(s, tagAttr)
	if err != nil {
		return "", err
	}
	s.SetAttr("directive", "fmt")
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return "", err
	}
	return htmlStr, nil
}

func renameTagsBySelector(fmtTag string, selectors ...string) (string, error) {
	for _, selector := range selectors {
		s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
		if err != nil {
			return "", err
		}
		var potErr error
		potErr = nil
		s.Find(selector).Each(func(i int, tag *goquery.Selection) {
			tagHtml, err := gqpp.GetHtmlFromSelection(tag)
			if err != nil {
				potErr = err
				return
			}
			tagAttr, _ := tag.Attr("tag")
			newTag, err := gqpp.ChangeSelectionTagName(tag, tagAttr)
			if err != nil {
				potErr = err
				return
			}
			newTag.SetAttr("directive", selector)
			newTagHtml, err := gqpp.GetHtmlFromSelection(newTag)
			if err != nil {
				potErr = err
				return
			}
			fmtTag = strings.Replace(fmtTag, tagHtml, newTagHtml, 1)
		})
		if potErr != nil {
			return "", potErr
		}
	}
	return fmtTag, nil
}

func setTagDepthByDirectiveAttrs(fmtTag string, targetDirectiveNames ...string) (string, error) {
	for _, targetDirectiveName := range targetDirectiveNames {
		s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
		if err != nil {
			return "", err
		}

		s.Find("*[directive='" + targetDirectiveName + "']").Each(func(i int, s *goquery.Selection) {
			directiveCount := 0
			gqpp.ClimbTreeUntil(s, func(parent *goquery.Selection) bool {
				dirAttr, _ := parent.Attr("directive")
				if dirAttr == "fmt" {
					return true
				}
				if dirAttr == "for" || dirAttr == "if" {
					directiveCount++
				}
				return false
			})

			s.SetAttr("depth", fmt.Sprintf("%d", directiveCount))
		})

		// Update fmtDir with the modified HTML
		htmlStr, err := gqpp.GetHtmlFromSelection(s)
		if err != nil {
			return "", err
		}
		fmtTag = htmlStr
	}
	return fmtTag, nil
}

func getFmtTagMaxDepth(fmtTag string) (int, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
	if err != nil {
		return -1, err
	}
	var potErr error
	potErr = nil
	depthHeight := 0
	s.Find("*[directive]").Each(func(i int, s *goquery.Selection) {
		depth, _ := s.Attr("depth")
		d, err := strconv.Atoi(depth)
		if err != nil {
			potErr = err
			return
		}
		if d > depthHeight {
			depthHeight = d
		}
	})
	if potErr != nil {
		return -1, err
	}
	return depthHeight, nil
}

func getInnerTags(fmtTag string) ([]string, error) {
	tags := make([]string, 0)
	s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
	if err != nil {
		return nil, err
	}
	dirNames := []string{"for", "if"}
	for _, name := range dirNames {
		var potErr error
		potErr = nil
		s.Find("*[directive='" + name + "']").Each(func(i int, sel *goquery.Selection) {
			htmlStr, err := gqpp.GetHtmlFromSelection(sel)
			if err != nil {
				potErr = err
				return
			}
			tags = append(tags, htmlStr)
		})
		if potErr != nil {
			return nil, potErr
		}
	}
	return tags, nil
}

func sortInnerTagsByDepth(tags []string, maxDepth int) ([]string, error) {
	out := make([]string, 0)
	currentDepth := maxDepth
	for {
		if currentDepth == -1 {
			break
		}
		foundTagAtDepth := false
		for _, tag := range tags {
			d, _, err := gqpp.AttrFromStr(tag, "depth")
			if err != nil {
				return nil, err
			}
			dInt, err := strconv.Atoi(d)
			if err != nil {
				return nil, err
			}
			if !parsley.SliceContains(out, tag) && dInt == currentDepth {
				foundTagAtDepth = true
				out = append(out, tag)
			}
		}
		if foundTagAtDepth {
			continue
		}
		currentDepth--
	}
	return out, nil
}

func getInnerDirectivesByDepth(fmtTag string) ([]string, error) {
	tags, err := getInnerTags(fmtTag)
	if err != nil {
		return nil, err
	}
	maxDepth, err := getFmtTagMaxDepth(fmtTag)
	if err != nil {
		return nil, err
	}
	sorted, err := sortInnerTagsByDepth(tags, maxDepth)
	if err != nil {
		return nil, err
	}
	return sorted, nil
}

func extractFmtTagProps(fmtTag string) ([]string, error) {
	props := make([]string, 0)
	inProp := false
	prop := ""
	for i, ch := range fmtTag {
		if i+1 > len(fmtTag)-1 {
			break
		}
		char := string(ch)
		nextChar := string(fmtTag[i+1])
		combined := char + nextChar
		if combined == "{{" {
			inProp = true
		}
		if inProp {
			prop += char
			if combined == "}}" {
				inProp = false
				prop += nextChar
				props = append(props, prop)
				prop = ""
			}
		}
	}
	return props, nil
}

func outputFmtTagProps(fmtTag string, props []string) (string, error) {
	for _, prop := range props {
		out := strings.Replace(prop, "{{", "", 1)
		out = strings.Replace(out, "}}", "", 1)
		out = parsley.Squeeze(out)
		out = "`+" + out + "+`"
		fmtTag = strings.Replace(fmtTag, prop, out, 1)
	}
	return fmtTag, nil
}

func makeFmtTagParamStr(fmtTag string, props []string) (string, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
	if err != nil {
		return "", err
	}
	paramStr := ""
	for _, prop := range props {
		sq := parsley.Squeeze(prop)
		sq = strings.Replace(sq, "{{", "", 1)
		sq = strings.Replace(sq, "}}", "", 1)
		if strings.Contains(sq, ".") {
			continue
		}
		// filtering out duplicates
		if strings.Contains(paramStr, sq+" string, ") {
			continue
		}
		// filtering out <for>'s with type="string"
		shouldContinue := false
		s.Find("*[directive='for']").Each(func(i int, s *goquery.Selection) {
			asAttr, _ := s.Attr("as")
			typeAttr, _ := s.Attr("type")
			if typeAttr == "string" {
				if asAttr == sq {
					shouldContinue = true
				}
			}
		})
		if shouldContinue {
			continue
		}
		paramStr += sq + " string, "
	}
	s.Find("*[directive='for']").Each(func(i int, s *goquery.Selection) {
		inAttr, _ := s.Attr("in")
		typeAttr, _ := s.Attr("type")
		if !strings.Contains(inAttr, ".") {
			paramStr += inAttr + " []" + typeAttr + ", "
		}
	})
	s.Find("*[directive='if']").Each(func(i int, s *goquery.Selection) {
		condAttr, _ := s.Attr("condition")
		paramStr += condAttr + " " + "bool" + ", "
	})
	paramStr = paramStr[0 : len(paramStr)-2]
	return paramStr, nil
}

func outputInnerTags(fmtTag string) (string, error) {
	innerDirs, err := getInnerDirectivesByDepth(fmtTag)
	if err != nil {
		panic(err)
	}
	foundDir := false
	for _, dir := range innerDirs {
		if strings.Contains(fmtTag, dir) {
			dirType, _, err := gqpp.AttrFromStr(dir, "directive")
			if err != nil {
				return "", err
			}
			if dirType == "for" {
				foundDir = true
				out, err := outputForTags(dir)
				if err != nil {
					return "", err
				}
				fmtTag = strings.Replace(fmtTag, dir, out, 1)
				break
			}
			if dirType == "if" {
				foundDir = true
				out, err := outputIfTags(dir)
				if err != nil {
					return "", err
				}
				fmtTag = strings.Replace(fmtTag, dir, out, 1)
				break
			}
		}
	}
	if foundDir {
		return outputInnerTags(fmtTag)
	}
	return fmtTag, nil
}

func outputForTags(forTag string) (string, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(forTag)
	if err != nil {
		return "", err
	}
	inAttr, _ := s.Attr("in")
	asAttr, _ := s.Attr("as")
	typeAttr, _ := s.Attr("type")
	tagAttr, _ := s.Attr("tag")
	attrStr := gqpp.GetAttrStr(s, "in", "as", "type", "directive")
	htmlStr, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return "", err
	}
	openTag := ""
	if len(attrStr) > 0 {
		openTag = fmt.Sprintf("<%s>", tagAttr)
	} else {
		openTag = fmt.Sprintf("<%s %s>", tagAttr, attrStr)
	}
	htmlStr = parsley.ReplaceFirstLine(htmlStr, openTag)
	out := fmt.Sprintf("`+collectStr(%s, func(i int, %s %s) string {return `%s`})+`", inAttr, asAttr, typeAttr, parsley.FlattenStr(htmlStr))
	return out, nil
}

func outputIfTags(ifTag string) (string, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(ifTag)
	if err != nil {
		return "", err
	}
	condAttr, _ := s.Attr("condition")
	tagAttr, _ := s.Attr("tag")
	attrStr := gqpp.GetAttrStr(s, "condition", "tag", "directive")
	elseSel := s.Find("else")
	elseHtml, err := goquery.OuterHtml(elseSel)
	if err != nil {
		return "", err
	}
	elseSel.Remove()
	ifHtml, err := gqpp.GetHtmlFromSelection(s)
	if err != nil {
		return "", err
	}
	openTag := ""
	if len(attrStr) > 0 {
		openTag = fmt.Sprintf("<%s>", tagAttr)
	} else {
		openTag = fmt.Sprintf("<%s %s>", tagAttr, attrStr)
	}
	ifHtml = parsley.ReplaceFirstLine(ifHtml, openTag)
	elseHtml = parsley.ReplaceFirstLine(elseHtml, parsley.GetFirstLine(ifHtml))
	elseHtml = parsley.ReplaceLastLine(elseHtml, parsley.GetLastLine(ifHtml))
	out := fmt.Sprintf("`+ ifElse(%s, `%s`, `%s`) +`", condAttr, parsley.FlattenStr(ifHtml), parsley.FlattenStr(elseHtml))
	return out, nil
}

func wrapFmtTagInGoFunc(fmtTag string, paramStr string) (string, error) {
	s, err := gqpp.NewSelectionFromHtmlStr(fmtTag)
	if err != nil {
		return "", err
	}
	nameAttr, _ := s.Attr("name")
	tagAttr, _ := s.Attr("tag")
	attrStr := gqpp.GetAttrStr(s, "in", "tag")
	openTag := ""
	if len(attrStr) > 0 {
		openTag = fmt.Sprintf("<%s %s>", tagAttr, attrStr)
	} else {
		openTag = fmt.Sprintf("<%s>", tagAttr)
	}
	fmtTag = parsley.ReplaceFirstLine(fmtTag, openTag)
	out := fmt.Sprintf("func %s(%s) string {\n\treturn`%s`\t}", nameAttr, paramStr, fmtTag)
	return parsley.FlattenStr(out), nil
}

func appendFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Append a newline character to the content
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func formatOutput(path string) error {
	// Check if gofmt is installed
	if _, err := exec.LookPath("gofmt"); err != nil {
		fmt.Println("Warning: <fmt> components cannot be formatted due to 'gofmt' not being installed. Please install it to enable formatting.")
		return nil
	}

	// Run gofmt if installed
	cmd := exec.Command("gofmt", "-w", path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to format file %s: %w", path, err)
	}

	return nil
}

func collectStr[T any](slice []T, mapper func(i int, t T) string) string {
	var builder strings.Builder
	for i, t := range slice {
		builder.WriteString(mapper(i, t))
	}
	return builder.String()
}

func ifElse(cond bool, ifTrue string, ifFalse string) string {
	if cond {
		return ifTrue
	}
	return ifFalse
}

// names := []string{"Alice", "Bob", "Charlie"}
// result := parsley.CollectStr(names, func(i int, name string) string {
// 	return fmt.Sprintf("<li>%s</li>", name)
// })
