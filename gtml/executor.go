package gtml

import (
	"fmt"
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
 Convert HTML to Golang 💦
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

	printIntro := func() error {
		intro := purse.Fmt(`
building %s 💦`, ex.OutputFile)
		fmt.Println(intro)
		return nil
	}

	funcs := make([]Func, 0)
	walkInputDir := func() error {
		err := filepath.Walk(ex.InputDir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil // skip all dirs
			}
			if !strings.HasSuffix(path, ".html") {
				return nil // skip all non .html files
			}
			// extract the html _components from the file
			compNames, err := ReadComponentElementNamesFromFile(path)
			if err != nil {
				return err
			}
			compSels, err := ReadComponentSelectionsFromFile(path)
			if err != nil {
				return err
			}
			for _, sel := range compSels {
				err := MarkSelectionPlaceholders(sel, compNames)
				if err != nil {
					return err
				}
			}
			MarkSelectionsAsUnique(compSels)
			compElms, err := ConvertSelectionsIntoElements(compSels, compNames)
			if err != nil {
				return err
			}
			for _, elm := range compElms {
				fn, err := NewFunc(elm, compElms)
				if err != nil {
					return err
				}
				funcs = append(funcs, fn)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}

	writeOutputToFile := func() error {
		file, err := os.OpenFile(ex.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return fmt.Errorf("failed to create or clear output file: %w", err)
		}
		defer file.Close()
		_, err = file.WriteString("package " + ex.PackageName + "\n\n")
		if err != nil {
			return fmt.Errorf("failed to write package declaration: %w", err)
		}
		_, err = file.WriteString(ex.ImportBlock + "\n\n")
		if err != nil {
			return fmt.Errorf("failed to write import block: %w", err)
		}
		for _, fn := range funcs {
			_, err = file.WriteString(fn.GetData() + "\n")
			if err != nil {
				return fmt.Errorf("failed to write function data: %w", err)
			}
		}
		return nil
	}

	process := func() error {
		err := fungi.Process(
			func() error { return printIntro() },
			func() error { return walkInputDir() },
			func() error { return writeOutputToFile() },
		)
		if err != nil {
			return err
		}
		funcs = make([]Func, 0)
		return nil
	}

	for _, opt := range ex.Command.GetOptions() {
		process = opt.Inject(ex, process)
	}

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
  gtml [OPTIONS]... [INPUT DIR] [OUTPUT FILE]

Example: 
  gtml --watch ./components output.go

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
