package cli

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlfunc"
	"gtml/src/parser/gtmlvar"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================

func getGtmlArt() string {

	art := `
   _____ _______ __  __ _
  / ____|__   __|  \/  | |
 | |  __   | |  | \  / | |
 | | |_ |  | |  | |\/| | |
 | |__| |  | |  | |  | | |____
  \_____|  |_|  |_|  |_|______|
 ---------------------------------------
 Make Writing HTML in Go a Breeze 🍃
 Version 0.1.9 (2024-12-6)
 https://github.com/phillip-england/gtml
 ---------------------------------------`
	return purse.RemoveFirstLine(art)
}

// ##==================================================================
type Executor interface {
	Run() error
	GetCommand() Command
}

func NewExecutor(cmd Command) (Executor, error) {
	if cmd.GetType() == KeyCommandBuild {
		ex, err := NewExecutorBuild(cmd)
		if err != nil {
			return nil, err
		}
		return ex, nil
	}
	if cmd.GetType() == KeyCommandHelp {
		ex, err := NewExecutorHelp(cmd)
		if err != nil {
			return nil, err
		}
		return ex, nil
	}
	msg := purse.Fmt(`
not executor available for command of type: %s
	`, cmd.GetType())
	return nil, fmt.Errorf(msg)
}

// ##==================================================================
type ExecutorBuild struct {
	Command          Command
	InputDir         string
	OutputFile       string
	PackageName      string
	OutputFileExists bool
}

func NewExecutorBuild(cmd Command) (*ExecutorBuild, error) {
	ex := &ExecutorBuild{
		Command: cmd,
	}
	err := fungi.Process(
		func() error { return ex.initInputDir() },
		func() error { return ex.initOutputFile() },
		func() error { return ex.initPackageName() },
		func() error { return ex.initOutputFileExists() },
	)
	if err != nil {
		return nil, err
	}
	return ex, nil
}

func (ex *ExecutorBuild) GetCommand() Command { return ex.Command }
func (ex *ExecutorBuild) Run() error {

	process := func() error {
		err := ex.printIntro()
		if err != nil {
			return err
		}
		funcs, err := ex.buildComponentFuncs()
		if err != nil {
			return err
		}
		err = ex.writeComponentFuncs(funcs)
		if err != nil {
			return err
		}
		return nil
	}

	for _, opt := range ex.Command.GetOptions() {
		process = opt.Inject(ex, process)
	}

	fmt.Println(getGtmlArt())

	err := process()
	if err != nil {
		return err
	}

	return nil

}

func (ex *ExecutorBuild) initInputDir() error {
	ex.InputDir = ex.Command.GetFilteredArgs()[0]
	return nil
}

func (ex *ExecutorBuild) initOutputFile() error {
	ex.OutputFile = ex.Command.GetFilteredArgs()[1]
	return nil
}

func (ex *ExecutorBuild) initPackageName() error {
	ex.PackageName = ex.Command.GetFilteredArgs()[2]
	return nil
}

func (ex *ExecutorBuild) initOutputFileExists() error {
	_, err := os.Stat(ex.OutputFile)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return err
		}
		ex.OutputFileExists = false
		return nil
	}
	ex.OutputFileExists = true
	return nil
}

func (ex *ExecutorBuild) printIntro() error {
	intro := purse.Fmt(`
building %s 💦`, ex.OutputFile)
	fmt.Println(intro)
	return nil
}

func (ex *ExecutorBuild) buildComponentFuncs() ([]gtmlfunc.Func, error) {
	funcs := make([]gtmlfunc.Func, 0)
	err := filepath.Walk(ex.InputDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil // skip all dirs
		}
		if !strings.HasSuffix(path, ".html") {
			return nil // skip all non .html files
		}

		// extract the html _components from the file
		compNames, err := element.ReadComponentElementNamesFromFile(path)
		if err != nil {
			return err
		}
		compSels, err := element.ReadComponentSelectionsFromFile(path)
		if err != nil {
			return err
		}
		for _, sel := range compSels {
			err := element.MarkSelectionPlaceholders(sel, compNames)
			if err != nil {
				return err
			}
		}
		element.MarkSelectionsAsUnique(compSels)
		compElms, err := element.ConvertSelectionsIntoElements(compSels, compNames)
		if err != nil {
			return err
		}
		for _, elm := range compElms {
			fn, err := gtmlfunc.NewFunc(elm, compElms)
			if err != nil {
				return err
			}
			funcs = append(funcs, fn)
		}
		return nil
	})
	if err != nil {
		return funcs, err
	}
	return funcs, nil
}

func (ex *ExecutorBuild) writeComponentFuncs(funcs []gtmlfunc.Func) error {
	// Ensure the directory exists
	outputDir := filepath.Dir(ex.OutputFile)
	err := os.MkdirAll(outputDir, 0755) // Create directories if they don't exist
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Open the file for writing, creating it if it doesn't exist
	file, err := os.OpenFile(ex.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to create or clear output file: %w", err)
	}
	defer file.Close()
	// buildIgnore := "// +build ignore\n"
	buildIgnore := ""
	_, err = file.WriteString("// Code generated by gtml; DO NOT EDIT.\n" + buildIgnore + "\n// v0.1.0 | you may see errors with types, you'll need to manage your own imports\n// type support coming soon!" + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write ignore declaration: %w", err)
	}

	// Write package declaration
	_, err = file.WriteString("package " + ex.PackageName + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write package declaration: %w", err)
	}

	// Write import block
	imports := []string{"\"strings\""}
	foundMd := false
	for _, fn := range funcs {
		vars := fn.GetVars()
		for _, v := range vars {
			if v.GetType() == gtmlvar.KeyVarGoMd {
				foundMd = true
			}
		}
	}
	if foundMd {
		imports = append(imports, "\t"+`chromahtml "github.com/alecthomas/chroma/v2/formatters/html"`)
		imports = append(imports, "\t"+`highlighting "github.com/yuin/goldmark-highlighting/v2"`)
		imports = append(imports, "\t"+`goldmarkhtml "github.com/yuin/goldmark/renderer/html"`)
		imports = append(imports, "\t"+`"github.com/yuin/goldmark"`)
		imports = append(imports, "\t"+`"github.com/yuin/goldmark/parser"`)
		imports = append(imports, "\t"+`"bytes"`)
		imports = append(imports, "\t"+`"os"`)
		imports = append(imports, "\t"+`"github.com/PuerkitoBio/goquery"`)
	}
	var importBlock string
	if len(imports) == 1 {
		importBlock = "import \"strings\""
	} else {
		importStr := strings.Join(imports, "\n")
		importBlock = purse.Fmt(`
import (
	%s
)
		`, importStr)
	}
	_, err = file.WriteString(importBlock + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write import block: %w", err)
	}

	// setting up gtmlMd
	var gtmlMd string
	if foundMd {
		gtmlMd = purse.RemoveFirstLine(`
func gtmlMd(mdPath string, theme string) string {
	mdFileContent, _ := os.ReadFile(mdPath)
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle(theme),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			goldmarkhtml.WithHardWraps(),
			goldmarkhtml.WithXHTML(),
			goldmarkhtml.WithUnsafe(),
		),
	)
	var buf bytes.Buffer
	_ = md.Convert([]byte(mdFileContent), &buf)
	str := buf.String()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))

	doc.Find("*").Each(func(i int, inner *goquery.Selection) {
		nodeName := goquery.NodeName(inner)
		currentStyle, _ := inner.Attr("style")
		switch nodeName {
		case "pre":
			inner.SetAttr("style", currentStyle+"padding: 1rem; font-size: 0.875rem; overflow-x: auto; border-radius: 0.25rem; margin-bottom: 1rem;")
		case "h1":
			inner.SetAttr("style", currentStyle+"font-weight: bold; font-size: 1.875rem; padding-bottom: 1rem;")
		case "h2":
			inner.SetAttr("style", currentStyle+"font-size: 1.5rem; font-weight: bold; padding-bottom: 1rem; padding-top: 0.5rem; border-top-width: 1px; border-top-style: solid; border-color: #1f2937; padding-top: 1rem;")
		case "h3":
			inner.SetAttr("style", currentStyle+"font-size: 1.25rem; font-weight: bold; margin-top: 1.5rem; margin-bottom: 1rem;")
		case "p":
			inner.SetAttr("style", currentStyle+"font-size: 0.875rem; line-height: 1.5; margin-bottom: 1rem;")
		case "ul":
			inner.SetAttr("style", currentStyle+"padding-left: 1.5rem; margin-bottom: 1rem; list-style-type: disc;")
		case "ol":
			inner.SetAttr("style", currentStyle+"padding-left: 1.5rem; margin-bottom: 1rem; list-style-type: decimal;")
		case "li":
			inner.SetAttr("style", currentStyle+"margin-bottom: 0.5rem;")
		case "blockquote":
			inner.SetAttr("style", currentStyle+"margin-left: 1rem; padding-left: 1rem; border-left: 4px solid #ccc; font-style: italic; color: #555;")
		case "code":
			parent := inner.Parent()
			if goquery.NodeName(parent) == "pre" {
				return
			}
			inner.SetAttr("style", currentStyle+"font-family: monospace; background-color: #1f2937; padding: 0.25rem 0.5rem; border-radius: 0.25rem;")
		case "hr":
			inner.SetAttr("style", currentStyle+"border: none; border-top: 1px solid #ccc; margin: 2rem 0;")
		case "a":
			inner.SetAttr("style", currentStyle+"color: #007BFF; text-decoration: none;")
		case "img":
			inner.SetAttr("style", currentStyle+"max-width: 100%; height: auto; border-radius: 0.25rem; margin: 1rem 0;")
		}
	})
	modifiedHTML, _ := doc.Html()
	return modifiedHTML
}
`)
	}

	// Write helper functions
	_, err = file.WriteString(purse.Fmt(`
func gtmlFor[T any](slice []T, callback func(i int, item T) string) string {
	var builder strings.Builder
	for i, item := range slice {
		builder.WriteString(callback(i, item))
	}
	return builder.String()
}

func gtmlIf(condition bool, fn func() string) string {
if condition {
	return fn()
}
	return ""
}

func gtmlElse(condition bool, fn func() string) string {
	if !condition {
		return fn()
	}
	return ""
}

func gtmlSlot(contentFunc func() string) string {
	return contentFunc()
}

func gtmlEscape(input string) string {
	return input
}

%s
`, gtmlMd))
	if err != nil {
		return fmt.Errorf("failed to write helper functions: %w", err)
	}

	// Write function data
	for _, fn := range funcs {
		_, err = file.WriteString(fn.GetData() + "\n\n")
		if err != nil {
			return fmt.Errorf("failed to write function data: %w", err)
		}
	}

	return nil
}

// ##==================================================================
type ExecutorHelp struct {
	Command Command
}

func NewExecutorHelp(cmd Command) (*ExecutorHelp, error) {
	ex := &ExecutorHelp{
		Command: cmd,
	}
	return ex, nil
}

func (ex *ExecutorHelp) GetCommand() Command { return ex.Command }
func (ex *ExecutorHelp) Run() error {
	message := fmt.Sprintf(`
%s
Usage:
  gtml [OPTIONS]... [INPUT DIR] [OUTPUT FILE] [GO PACKAGE NAME]

Example:
  gtml --watch build ./components output.go output

Options:
  --watch       rebuild when source files are modified
`, getGtmlArt())
	message = purse.RemoveFirstLine(message)
	fmt.Println(message)
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
