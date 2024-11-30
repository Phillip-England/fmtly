package cli

import (
	"fmt"
	"gtml/src/parser/element"
	"gtml/src/parser/gtmlfunc"
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
 HTML Components in Go Made Easy 💦
 Version 0.1.0 (2024-11-26)
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
	ImportBlock      string
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
		func() error { return ex.initImportBlock() },
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

func (ex *ExecutorBuild) initImportBlock() error {
	ex.ImportBlock = purse.Fmt(`
import "strings"`)
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

	_, err = file.WriteString("// Code generated by gtml; DO NOT EDIT.\n// +build ignore\n// v0.1.0 | you may see errors with types, you'll need to manage your own imports\n// type support coming soon!" + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write ignore declaration: %w", err)
	}

	// Write package declaration
	_, err = file.WriteString("package " + ex.PackageName + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write package declaration: %w", err)
	}

	// Write import block
	_, err = file.WriteString(ex.ImportBlock + "\n\n")
	if err != nil {
		return fmt.Errorf("failed to write import block: %w", err)
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

`))
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
