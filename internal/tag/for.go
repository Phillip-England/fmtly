package tag

import (
	"fmt"
	"fmtly/internal/parsley"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ForTag struct {
	Selection   *goquery.Selection
	Html        string
	AttrStr     string
	HtmlOutput  string
	ParamOutput string
	InAttr      string
	AsAttr      string
	TypeAttr    string
	TagAttr     string
}

func NewForTagFromSelection(selection *goquery.Selection) (*ForTag, error) {
	t := &ForTag{
		Selection: selection,
	}
	if err := t.setHtml(); err != nil {
		return nil, err
	}
	if err := t.setAttrStr(); err != nil {
		return nil, err
	}
	if err := t.setAttrs(); err != nil {
		return nil, err
	}
	if err := t.setHtmlOutput(); err != nil {
		return nil, err
	}
	if err := t.setParamOutput(); err != nil {
		return nil, err
	}
	if err := t.wrapHtmlOutputInGo(); err != nil {
		return nil, err
	}
	return t, nil
}

func (t *ForTag) setHtml() error {
	htmlStr, err := goquery.OuterHtml(t.Selection)
	if err != nil {
		return err
	}
	flatHtml := parsley.FlattenStr(htmlStr)
	t.Html = flatHtml
	return nil
}

func (t *ForTag) setAttrStr() error {
	attrStr := ""
	for i, node := range t.Selection.Nodes {
		if i == 0 {
			for _, attr := range node.Attr {
				if parsley.EqualsOneof(attr.Key, "in", "as", "tag", "type") {
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

func (t *ForTag) setAttrs() error {
	inAttr, exists := t.Selection.Attr("in")
	if !exists {
		return fmt.Errorf("<for> is missing 'in' attribute:\n\n%s", t.Html)
	}
	asAttr, exists := t.Selection.Attr("as")
	if !exists {
		return fmt.Errorf("<for> is missing 'as' attribute:\n\n%s", t.Html)
	}
	typeAttr, exists := t.Selection.Attr("type")
	if !exists {
		typeAttr = "any"
	}
	tagAttr, exists := t.Selection.Attr("tag")
	if !exists {
		return fmt.Errorf("<for> is missing 'tag' attribute:\n\n%s", t.Html)
	}
	t.InAttr = inAttr
	t.AsAttr = asAttr
	t.TypeAttr = typeAttr
	t.TagAttr = tagAttr
	return nil
}

func (t *ForTag) setHtmlOutput() error {
	htmlBody, err := t.Selection.Html()
	if err != nil {
		return err
	}
	t.HtmlOutput = fmt.Sprintf("<%s %s>%s</%s>", t.TagAttr, t.AttrStr, parsley.FlattenStr(htmlBody), t.TagAttr)
	return nil
}

func (t *ForTag) setParamOutput() error {
	if strings.Contains(t.InAttr, ".") {
		return nil
	}
	t.ParamOutput = t.InAttr
	return nil
}

func (t *ForTag) wrapHtmlOutputInGo() error {
	out := fmt.Sprintf(parsley.GetTick() + " + collectStr(" + t.InAttr + " , func(i int, " + t.AsAttr + " " + t.TypeAttr + ") string { return `" + t.HtmlOutput + "` }) + " + parsley.GetTick())
	t.HtmlOutput = out
	return nil
}
