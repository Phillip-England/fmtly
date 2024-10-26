package fmtly

import (
	"strings"
	"testing"
)

func FormShell(title string, children ...string) string {
	return `
		<form fmt="FormShell">
		    <div>
		        <h2>` + title + `</h2>
		    </div>
		    <ul>
		        ` + strings.Join(children, "") + `
		    </ul>
		</form>
	`
}

func NavItem(href string, text string) string {
	return `
		<li fmt="NavItem">
		    <a href="` + href + `">` + text + `</a>
		</li>
	`
}

func LoginForm() string {
	return FormShell("Login", NavItem("/", "Home"), NavItem("/", "Home"), NavItem("/", "Home"))
}

type Customer struct {
	Name string
}

func CustomerList(customerSlice []Customer) string {
	var loop1Slice []string
	for i := 0; i < len(customerSlice); i++ {
		currentCustomer := customerSlice[i]
		loop1Slice = append(loop1Slice, `
			<li>
				<p>`+currentCustomer.Name+`</p>
			</li>
		`)
	}
	return `
		<ul fmt"CustomerList">
			` + strings.Join(loop1Slice, "\n") + `
		</ul>
	`
}

func TestMain(m *testing.M) {

	targetDir := "./components"
	compStr, err := ExtractCompStr(targetDir)
	if err != nil {
		panic(err)
	}
	comps, err := ExtractFmtlyComponents(compStr)
	if err != nil {
		panic(err)
	}
	SetFmtlyTemplateStatementDetails(comps)
	CreateFmtlyGoOutputTemplates(comps)
	SetFmtlyGoOutputName(comps)
	err = SortFmtlyTemplateStatements(comps)
	if err != nil {
		panic(err)
	}
	err = WriteOutStatementParams(comps)
	if err != nil {
		panic(err)
	}
	err = WriteFuncBodyForTemplates(comps)
	if err != nil {
		panic(err)
	}
}
