package fmtly

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

//=======================================
// CONSTS
//=======================================

const (
	phFuncName   = "FUNCNAME"
	phFuncParams = "FUNCPARAMS"
	phFuncBody   = "FUNCBODY"
)

const (
	stmtProp     = "PROP"
	stmtTypeProp = "TYPEPROP"
	stmtFor      = "FOR"
	stmtEnd      = "END"
	stmtElse     = "ELSE"
	stmtIf       = "IF"
	stmtIfElse   = "IFELSE"
	stmtChildren = "CHILDREN"
)

//=======================================
// COMPONENT
//=======================================

type Component struct {
	HTML                      string
	Name                      string
	TemplateStatementsDetails []*TemplateStatementDetails
	TemplateStatements        []TemplateStatement
	GoOutput                  string
}

//=======================================
// STATEMENTS
//=======================================

type TemplateStatementDetails struct {
	FoundAs string
	RawText string
	Value   string
	Order   int
}

type TemplateStatement interface {
	GetGoFuncParam() string
	GetStatementType() string
	WriteGoFuncBody(comp *Component) error
}

//=======================================
// PROP
//=======================================

type PropStatement struct {
	Details TemplateStatementDetails
}

func NewPropStatement(details TemplateStatementDetails) PropStatement {
	propStatement := PropStatement{
		Details: details,
	}
	return propStatement
}

func (f PropStatement) GetGoFuncParam() string {
	return f.Details.Value + " string,"
}

func (f PropStatement) GetStatementType() string {
	return stmtProp
}

func (f PropStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// TYPEPROP
//=======================================

type TypePropStatement struct {
	Details  TemplateStatementDetails
	TypeName string
}

func NewTypePropStatement(details TemplateStatementDetails) TypePropStatement {
	typePropStatement := TypePropStatement{
		Details: details,
	}
	return typePropStatement
}

func (f TypePropStatement) GetGoFuncParam() string {
	return ""
}

func (f TypePropStatement) GetStatementType() string {
	return stmtTypeProp
}

func (f TypePropStatement) GetParamName() string {
	str := strings.Replace(f.Details.RawText, "{{", "", 1)
	str = strings.Replace(str, "}}", "", 1)
	return strings.Split(str, ".")[0]
}

func (f TypePropStatement) GetGoProperty() string {
	str := strings.Replace(f.Details.RawText, "{{", "", 1)
	str = strings.Replace(str, "}}", "", 1)
	return strings.Split(str, ".")[1]
}

func (f TypePropStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// FOR
//=======================================

type ForStatement struct {
	Details TemplateStatementDetails
	Parts   []string
}

func NewForStatement(details TemplateStatementDetails) ForStatement {
	forStatement := ForStatement{
		Details: details,
	}
	return forStatement
}

func (f ForStatement) GetGoFuncParam() string {
	forParam := f.GetForParamName()
	forType := f.GetForType()
	return forParam + " " + forType + ","
}

func (f ForStatement) GetStatementType() string {
	return stmtFor
}

func (f ForStatement) GetForParamName() string {
	parts := strings.Split(f.Details.Value, " ")
	return parts[1]
}

func (f ForStatement) GetForType() string {
	parts := strings.Split(f.Details.Value, " ")
	return parts[3]
}

func (f ForStatement) GetForOutputVarName() string {
	return fmt.Sprintf("forOutput%d", f.Details.Order)
}

func (f ForStatement) ExtractForHTMLFromComponent(comp *Component) (string, error) {
	indexOfStatment := strings.Index(comp.HTML, f.Details.FoundAs)
	choppedHTML := comp.HTML[indexOfStatment:]
	indexOfFirstClose := strings.Index(choppedHTML, "}}")
	choppedHTML = choppedHTML[indexOfFirstClose:]
	choppedHTML = strings.Replace(choppedHTML, "}}", "", 1)
	var goodLines []string
	foundEnd := false
	lines := strings.Split(choppedHTML, "\n")
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if foundEnd {
			break
		}
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") && strings.Contains(line, "end") {
			foundEnd = true
			line = strings.Replace(line, "{{", "", 1)
			line = strings.Replace(line, "}}", "", 1)
			line = strings.Replace(line, "end", "", 1)
			if len(line) == 0 {
				continue
			}
		}
		goodLines = append(goodLines, line)
	}
	return strings.Join(goodLines, ""), nil
}

func (f ForStatement) WriteGoFuncBody(comp *Component) error {
	varName := f.GetForOutputVarName()
	comp.GoOutput = strings.Replace(comp.GoOutput, phFuncBody,
		fmt.Sprintf("%s := \"\"", varName)+"\n\t\t"+phFuncBody, 1)
	comp.GoOutput = strings.Replace(comp.GoOutput, phFuncBody,
		fmt.Sprintf("for i := 0; i < len(%s); i++ {", f.GetForParamName())+"\n\t\t\t"+phFuncBody, 1)
	comp.GoOutput = strings.Replace(comp.GoOutput, phFuncBody,
		fmt.Sprintf("resource := %s[i]", f.GetForParamName())+"\n\t\t\t"+phFuncBody, 1)
	// forHTML, err := f.ExtractForHTMLFromComponent(comp)
	// if err != nil {
	// 	return err
	// }

	return nil
	// return output
}

//=======================================
// END
//=======================================

type EndStatement struct {
	Details TemplateStatementDetails
}

func NewEndStatement(details TemplateStatementDetails) EndStatement {
	endStatement := EndStatement{
		Details: details,
	}
	return endStatement
}

func (f EndStatement) GetGoFuncParam() string {
	return ""
}

func (f EndStatement) GetStatementType() string {
	return stmtEnd
}

func (f EndStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// ELSE
//=======================================

type ElseStatement struct {
	Details TemplateStatementDetails
}

func NewElseStatement(details TemplateStatementDetails) ElseStatement {
	elseStatement := ElseStatement{
		Details: details,
	}
	return elseStatement
}

func (f ElseStatement) GetGoFuncParam() string {
	return ""
}

func (f ElseStatement) GetStatementType() string {
	return stmtElse
}

func (f ElseStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// IF
//=======================================

type IfStatement struct {
	Details      TemplateStatementDetails
	HTMLEndIndex int
}

func NewIfStatement(details TemplateStatementDetails) IfStatement {
	ifStatement := IfStatement{
		Details: details,
	}
	return ifStatement
}

func (f IfStatement) GetGoFuncParam() string {
	str := f.GetIfValue()
	return str + " bool,"
}

func (f IfStatement) GetStatementType() string {
	return stmtIf
}

func (f IfStatement) GetIfValue() string {
	str := strings.ReplaceAll(f.Details.RawText, "{{if", "")
	str = strings.ReplaceAll(str, "}}", "")
	return str
}

func (f IfStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// IFELSE
//=======================================

type IfElseStatement struct {
	Details TemplateStatementDetails
}

func NewIfElseStatement(details TemplateStatementDetails) IfElseStatement {
	ifElseStatement := IfElseStatement{
		Details: details,
	}
	return ifElseStatement
}

func (f IfElseStatement) GetGoFuncParam() string {
	str := f.GetIfValue()
	return str + " bool,"
}

func (f IfElseStatement) GetStatementType() string {
	return stmtIfElse
}

func (f IfElseStatement) GetIfValue() string {
	str := strings.ReplaceAll(f.Details.RawText, "{{if", "")
	str = strings.ReplaceAll(str, "}}", "")
	return str
}

func (f IfElseStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// CHILDREN
//=======================================

type ChildrenStatement struct {
	Details TemplateStatementDetails
}

func (f ChildrenStatement) GetGoFuncParam() string {
	return f.Details.Value + " []string,"
}

func (f ChildrenStatement) GetStatementType() string {
	return stmtChildren
}

func NewChildrenStatement(details TemplateStatementDetails) ChildrenStatement {
	childrenStatement := ChildrenStatement{
		Details: details,
	}
	return childrenStatement
}

func (f ChildrenStatement) WriteGoFuncBody(comp *Component) error {
	return nil
}

//=======================================
// PROCEDURE FUNCS
//=======================================

func GetGoQueryDoc(htmlStr string) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func ExtractCompStr(targetDir string) (string, error) {
	var compStr string
	err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			f, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			fContent := string(f)
			compStr = compStr + fContent
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return compStr, nil
}

func ExtractFmtlyComponents(htmlStr string) ([]*Component, error) {
	var components []*Component
	doc, err := GetGoQueryDoc(htmlStr)
	if err != nil {
		return components, err
	}
	var errMessage error
	errMessage = nil
	doc.Find("define").Each(func(i int, s *goquery.Selection) {
		name, nameExists := s.Attr("name")
		if !nameExists {
			errMessage = fmt.Errorf("all <define> tags require a name")
		}
		compHTML, err := s.Html()
		if err != nil {
			errMessage = err
		}
		comp := Component{
			Name: name,
			HTML: compHTML,
		}
		components = append(components, &comp)
	})
	if errMessage != nil {
		return components, errMessage
	}
	return components, nil
}

func SetFmtlyTemplateStatementDetails(comps []*Component) {
	for _, comp := range comps {
		insideTemplateStatement := false
		statement := ""
		statementsFound := 0

		for i := 0; i < len(comp.HTML)-1; i++ { // Ensure we have bounds for nextChar
			strChar := string(comp.HTML[i])
			nextChar := string(comp.HTML[i+1])

			if !insideTemplateStatement && strChar == "{" && nextChar == "{" {
				insideTemplateStatement = true
				i++ // Skip nextChar to avoid including "{{"
				continue
			}

			if insideTemplateStatement {
				if strChar == "}" && nextChar == "}" {
					insideTemplateStatement = false
					value := statement
					statement += "}}"
					statement = "{{" + statement
					statementRawText := statement
					templateStatement := TemplateStatementDetails{
						FoundAs: statementRawText,
						RawText: strings.ReplaceAll(statementRawText, " ", ""),
						Value:   strings.TrimSuffix(strings.TrimPrefix(value, " "), " "),
						Order:   statementsFound,
					}
					statementsFound++
					comp.TemplateStatementsDetails = append(comp.TemplateStatementsDetails, &templateStatement)
					statement = ""
					i++ // Skip nextChar to avoid including "}}"
				} else {
					statement += strChar
				}
			}
		}
	}
}

func CreateFmtlyGoOutputTemplates(comps []*Component) {
	for _, comp := range comps {
		comp.GoOutput = fmt.Sprintf("func %s(%s) string {\n\t\t%s\n}", phFuncName, phFuncParams, phFuncBody)
	}
}

func SetFmtlyGoOutputName(comps []*Component) {
	for _, comp := range comps {
		comp.GoOutput = strings.Replace(comp.GoOutput, phFuncName, comp.Name, 1)
	}
}

func SortFmtlyTemplateStatements(comps []*Component) error {
	for _, comp := range comps {
		for i2, st := range comp.TemplateStatementsDetails {
			if strings.Contains(st.RawText, "{{if") {
				foundElse := false
				for i3, st2 := range comp.TemplateStatementsDetails {
					if i3 > i2 {
						if st2.RawText == "{{else}}" {
							foundElse = true
						}
					}
				}
				if foundElse == true {
					comp.TemplateStatements = append(comp.TemplateStatements, NewIfElseStatement(*st))

				} else {
					comp.TemplateStatements = append(comp.TemplateStatements, NewIfStatement(*st))
				}
				continue
			}
			if strings.Contains(st.RawText, "{{for") {
				comp.TemplateStatements = append(comp.TemplateStatements, NewForStatement(*st))
				continue
			}
			if st.RawText == "{{end}}" {
				comp.TemplateStatements = append(comp.TemplateStatements, NewEndStatement(*st))
				continue
			}
			if st.RawText == "{{else}}" {
				comp.TemplateStatements = append(comp.TemplateStatements, NewElseStatement(*st))
				continue
			}
			if strings.Contains(st.RawText, "{{...children}}") {
				comp.TemplateStatements = append(comp.TemplateStatements, NewChildrenStatement(*st))
				continue
			}
			if strings.Contains(st.RawText, ".") {
				comp.TemplateStatements = append(comp.TemplateStatements, NewTypePropStatement(*st))
				continue
			}
			comp.TemplateStatements = append(comp.TemplateStatements, NewPropStatement(*st))
		}
	}
	return nil
}

func WriteOutStatementParams(comps []*Component) error {
	for _, comp := range comps {
		for _, temp := range comp.TemplateStatements {
			paramStr := temp.GetGoFuncParam()
			if paramStr == "" {
				continue
			}
			comp.GoOutput = strings.Replace(comp.GoOutput, phFuncParams, paramStr+" "+phFuncParams, 1)
		}
		comp.GoOutput = strings.Replace(comp.GoOutput, ", "+phFuncParams, "", 1)
		comp.GoOutput = strings.Replace(comp.GoOutput, phFuncParams, "", 1)
	}
	return nil
}

func WriteFuncBodyForTemplates(comps []*Component) error {
	for _, comp := range comps {
		forCount := 0
		for _, temp := range comp.TemplateStatements {
			typeof := temp.GetStatementType()
			if typeof == stmtFor {
				forCount++
				err := temp.WriteGoFuncBody(comp)
				if err != nil {
					return err
				}
				fmt.Println(comp.GoOutput)
			}
		}
	}
	return nil
}
