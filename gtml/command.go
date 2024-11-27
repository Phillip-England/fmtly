package gtml

import (
	"fmt"
	"os"
	"strings"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyCommandBuild = "build"
	KeyCommandHelp  = "help"
)

// ##==================================================================

func getCommandList() []string {
	return []string{KeyCommandBuild, KeyCommandHelp}
}

func errHelp() string {
	return "Run 'gtml help' for usage."
}

// ##==================================================================
type Command interface {
	Print()
	GetType() string
	GetFilteredArgs() []string
	GetOptions() []Option
}

// ##==================================================================
func NewCommand() (Command, error) {
	args := os.Args
	if len(args) == 1 {
		fmt.Println(errHelp())
		return nil, nil
	}
	args = args[1:]
	opts := make([]Option, 0)
	for _, arg := range args {
		match := purse.FindMatchInStrSlice(getCommandList(), arg)
		if match == "" {
			opt, err := NewOption(arg)
			if err != nil {
				return nil, err
			}
			opts = append(opts, opt)
			continue
		}
		if match == KeyCommandHelp {
			cmd, err := NewCommandHelp()
			if err != nil {
				return nil, err
			}
			return cmd, nil
		}
		if match == KeyCommandBuild {
			cmd, err := NewCommandBuild(opts)
			if err != nil {
				return nil, err
			}
			return cmd, nil
		}
	}
	fmt.Println(errHelp())
	return nil, nil
}

// ##==================================================================
type CommandBuild struct {
	Type         string
	Options      []Option
	FilteredArgs []string
	Whitelist    []string
}

func NewCommandBuild(opts []Option) (*CommandBuild, error) {
	cmd := &CommandBuild{
		Type:    KeyCommandBuild,
		Options: opts,
	}
	err := fungi.Process(
		func() error { return cmd.initFilteredArgs() },
		func() error { return cmd.initInputWhitelist() },
		func() error { return cmd.initValidateInputDir() },
		func() error { return cmd.initValidateOutputFile() },
		func() error { return cmd.initValidatePackageName() },
	)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (cmd *CommandBuild) Print()                    { fmt.Println(cmd.Type) }
func (cmd *CommandBuild) GetType() string           { return cmd.Type }
func (cmd *CommandBuild) GetFilteredArgs() []string { return cmd.FilteredArgs }
func (cmd *CommandBuild) GetOptions() []Option      { return cmd.Options }

func (cmd *CommandBuild) initFilteredArgs() error {
	filtered := make([]string, 0)
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		isOpt := false
		for _, opt := range cmd.Options {
			if arg == opt.GetType() {
				isOpt = true
				continue
			}
		}
		if isOpt {
			continue
		}
		if arg == cmd.Type {
			continue
		}
		filtered = append(filtered, arg)
	}
	if len(filtered) != 3 {
		msg := purse.Fmt(`
gtml build has 3 required args
gtml build [INPUT DIR] [OUTPUT FILE] [PACKAGE NAME]
%s`, errHelp())
		return fmt.Errorf(msg)
	}
	cmd.FilteredArgs = filtered

	return nil
}

func (cmd *CommandBuild) initInputWhitelist() error {
	whitelist := make([]string, 0)
	whitelist = append(whitelist, purse.GetAllLetters()...)
	whitelist = append(whitelist, purse.GetAllNumbers()...)
	whitelist = append(whitelist, "-")
	whitelist = append(whitelist, "_")
	whitelist = append(whitelist, "/")
	whitelist = append(whitelist, ".")
	cmd.Whitelist = whitelist
	return nil
}

func (cmd *CommandBuild) initValidateInputDir() error {
	inputDir := cmd.FilteredArgs[0]
	if len(inputDir) == 0 {
		return fmt.Errorf("gtml build requires an input directory.\n" + errHelp())
	}
	inputDirFirstChar := string(inputDir[0])
	if inputDirFirstChar != "." {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid input directory provided: %s
input directory must start with  '.' like: './path/to/component_dir'
%s`, inputDir, errHelp()))
		return fmt.Errorf(msg)
	}
	if len(inputDir) > 1 {
		inputDirSecondChar := string(inputDir[1])
		if inputDirSecondChar != "/" {
			msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid input directory provided: %s
valid directory example: './path/to/component_dir'
if you are on windows, you may not use conventional windows directory paths
%s`, inputDir, errHelp()))
			return fmt.Errorf(msg)
		}
	}
	if strings.Contains(inputDir, "//") {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid input directory provided: %s
your input directory contains two '//' characters in a row
%s`, inputDir, errHelp()))
		return fmt.Errorf(msg)
	}
	if strings.Count(inputDir, ".") > 1 {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid input directory receieved: %s
your input directory contains more than one '.'
only characters, underscores, or hyphens are valid
%s`, inputDir, errHelp()))
		return fmt.Errorf(msg)
	}
	if !purse.EnforeWhitelist(inputDir, cmd.Whitelist) {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid input directory receieved: %s
only characters, underscores, or hyphens are valid
%s`, inputDir, errHelp()))
		return fmt.Errorf(msg)
	}
	return nil
}

func (cmd *CommandBuild) initValidateOutputFile() error {
	outputFile := cmd.FilteredArgs[1]
	if len(outputFile) == 0 {
		return fmt.Errorf("gtml build requires an output file.\n" + errHelp())
	}

	if !strings.HasSuffix(outputFile, ".go") {
		msg := purse.Fmt(`
invalid output file provided: %s
output file must end in '.go'
%s`, outputFile, errHelp())
		return fmt.Errorf(msg)
	}
	if !strings.HasPrefix(outputFile, "./") {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid output file provided: %s
output file must start with  '.' like: './output.go'
%s`, outputFile, errHelp()))
		return fmt.Errorf(msg)
	}
	if strings.Contains(outputFile, "//") {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid output file provided: %s
youre output file contains two '//' characters in a row
%s`, outputFile, errHelp()))
		return fmt.Errorf(msg)
	}
	if strings.Count(outputFile, ".") > 2 {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid output file provided: %s
your output file contains more than two '.' characters
only characters, underscores, or hyphens are valid
%s`, outputFile, errHelp()))
		return fmt.Errorf(msg)
	}
	if !purse.EnforeWhitelist(outputFile, cmd.Whitelist) {
		msg := purse.RemoveFirstLine(fmt.Sprintf(`
invalid output file provided: %s
only characters, underscores, or hyphens are valid
%s`, outputFile, errHelp()))
		return fmt.Errorf(msg)
	}
	return nil
}

func (cmd *CommandBuild) initValidatePackageName() error {
	packageName := cmd.FilteredArgs[2]
	whitelist := purse.GetAllLowerCaseLetters()
	whitelist = append(whitelist, []string{"-", "_"}...)
	pass := purse.EnforeWhitelist(packageName, whitelist)
	if !pass {
		msg := purse.Fmt(`
invalid package name provided: %s
package name may only have lowercase characters, underscores, and hyphens
%s`, packageName, errHelp())
		return fmt.Errorf(msg)
	}
	return nil
}

// ##==================================================================
type CommandHelp struct {
	Type         string
	FilteredArgs []string
	Options      []Option
}

func NewCommandHelp() (*CommandHelp, error) {
	cmd := &CommandHelp{
		Type: KeyCommandHelp,
	}
	return cmd, nil
}

func (cmd *CommandHelp) Print()                    { fmt.Println(cmd.Type) }
func (cmd *CommandHelp) GetType() string           { return cmd.Type }
func (cmd *CommandHelp) GetFilteredArgs() []string { return cmd.FilteredArgs }
func (cmd *CommandHelp) GetOptions() []Option      { return cmd.Options }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
