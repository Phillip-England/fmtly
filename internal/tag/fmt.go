package tag

import (
	"fmt"
	"fmtly/internal/parsley"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Fmt struct {
	Selection   *goquery.Selection
	Html        string
	HtmlOutput  string
	ParamOutput string
	AttrStr     string
	NameAttr    string
	TagAttr     string
	ForTags     []*ForTag
	IfTags      []*IfTag
	StrProps    []*Prop
	ForProps    []*Prop
}

type Prop struct {
	Raw   string
	Value string
}

func NewFmtTagFromStr(s string) (*Fmt, error) {
	t := &Fmt{}
	if err := t.setSelection(s); err != nil {
		return nil, err
	}
	if err := t.setHtml(); err != nil {
		return nil, err
	}
	if err := t.setForTags(); err != nil {
		return nil, err
	}
	if err := t.setForTagsHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setIfTags(); err != nil {
		return nil, err
	}
	if err := t.setIfTagsHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setAttrStr(); err != nil {
		return nil, err
	}
	if err := t.setAttrs(); err != nil {
		return nil, err
	}
	if err := t.setProps(); err != nil {
		return nil, err
	}
	if err := t.setParamOutput(); err != nil {
		return nil, err
	}
	if err := t.setHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setHtmlOutputProps(); err != nil {
		return nil, err
	}
	if err := t.wrapHtmlOutput(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Fmt) setSelection(s string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		return err
	}
	selection := doc.Find("fmt")
	t.Selection = selection
	return nil
}

func (t *Fmt) setHtml() error {
	htmlStr, err := goquery.OuterHtml(t.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	t.Html = flatHtml
	t.HtmlOutput = flatHtml
	return nil
}

func (t *Fmt) setIfTags() error {
	var potErr error
	potErr = nil
	t.Selection.Find("if").Each(func(i int, s *goquery.Selection) {
		ifTag, err := NewIfTagFromSelection(s)
		if err != nil {
			potErr = err
			return
		}
		t.IfTags = append(t.IfTags, ifTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (t *Fmt) setIfTagsHtmlOutput() error {
	for _, ft := range t.IfTags {
		fmt.Println(ft.Html)
		if strings.Contains(t.HtmlOutput, ft.Html) {
			fmt.Println("hit")
			t.HtmlOutput = strings.Replace(t.HtmlOutput, ft.Html, ft.HtmlOutput, 1)
		}
	}
	return nil
}

func (t *Fmt) setForTags() error {
	var potErr error
	potErr = nil
	t.Selection.Find("for").Each(func(i int, s *goquery.Selection) {
		forTag, err := NewForTagFromSelection(s)
		if err != nil {
			potErr = err
			return
		}
		t.ForTags = append(t.ForTags, forTag)
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func (t *Fmt) setForTagsHtmlOutput() error {
	for _, ft := range t.ForTags {
		if strings.Contains(t.HtmlOutput, ft.Html) {
			t.HtmlOutput = strings.Replace(t.HtmlOutput, ft.Html, ft.HtmlOutput, 1)
		}
	}
	return nil
}

func (t *Fmt) setAttrStr() error {
	attrStr := ""
	for i, node := range t.Selection.Nodes {
		if i == 0 {
			for _, attr := range node.Attr {
				if parsley.EqualsOneof(attr.Key, "name", "tag") {
					continue
				}
				attrStr += fmt.Sprintf("%s=\"%s\" ", attr.Key, attr.Val)
			}
		}
	}
	if len(attrStr) != 0 {
		attrStr = attrStr[:len(attrStr)-1]
	}
	t.AttrStr = attrStr
	return nil
}

func (t *Fmt) setAttrs() error {
	nameAttr, exists := t.Selection.Attr("name")
	if !exists {
		return fmt.Errorf("<fmt> is missing 'name' attribute:\n\n%s", t.Html)
	}
	tagAttr, exists := t.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<fmt> is missing 'tag' attribute:\n\n%s", t.Html)
	}
	t.NameAttr = nameAttr
	t.TagAttr = tagAttr
	return nil
}

func (t *Fmt) setProps() error {
	var props []string
	inProp := false
	prop := ""
	for i, ch := range t.Html {
		if i > len(t.Html)-2 {
			break
		}
		char := string(ch)
		nextChar := string(t.Html[i+1])
		search := char + nextChar
		if search == "{{" {
			inProp = true
		}
		if inProp {
			prop += char
			if search == "}}" {
				prop += nextChar
				props = append(props, prop)
				prop = ""
				inProp = false
			}
		}
	}
	for _, prop := range props {
		value := strings.Replace(prop, "{{", "", 1)
		value = strings.Replace(value, "}}", "", 1)
		value = parsley.Squeeze(value)
		if strings.Contains(value, ".") {
			t.ForProps = append(t.ForProps, &Prop{
				Raw:   prop,
				Value: value,
			})
		} else {
			t.StrProps = append(t.StrProps, &Prop{
				Raw:   prop,
				Value: value,
			})
		}
	}
	return nil
}

func (t *Fmt) setParamOutput() error {
	out := ""
	for _, prop := range t.StrProps {
		out = out + prop.Value + " string, "
	}
	for _, forTag := range t.ForTags {
		inAttr := forTag.InAttr
		if strings.Contains(inAttr, ".") {
			continue
		}
		out = out + inAttr + " []" + forTag.TypeAttr + ", "
	}
	out = out[:len(out)-2]
	t.ParamOutput = out
	return nil
}

func (t *Fmt) setHtmlOutput() error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(t.HtmlOutput))
	if err != nil {
		return err
	}
	s := doc.Find("fmt")
	innerHtml, err := s.Html()
	if err != nil {
		return err
	}
	out := fmt.Sprintf("<%s %s>%s</%s>", t.TagAttr, t.AttrStr, innerHtml, t.TagAttr)
	out = parsley.FlattenStr(out)
	t.HtmlOutput = out
	return nil
}

func (t *Fmt) setHtmlOutputProps() error {
	for _, p := range t.ForProps {
		t.HtmlOutput = strings.ReplaceAll(t.HtmlOutput, p.Raw, "` + "+p.Value+" + `")
	}
	for _, p := range t.StrProps {
		t.HtmlOutput = strings.ReplaceAll(t.HtmlOutput, p.Raw, "` + "+p.Value+" + `")
	}
	return nil
}

func (t *Fmt) wrapHtmlOutput() error {
	t.HtmlOutput = fmt.Sprintf("func " + t.NameAttr + "(" + t.ParamOutput + ") string { return `" + t.HtmlOutput + "` }")
	return nil
}
