package gtml

import (
	"fmt"
	"os"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyCommandBuild = "build"
	KeyCommandHelp  = "help"
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

func getCommandList() []string {
	return []string{KeyCommandBuild, KeyCommandHelp}
}

// ##==================================================================
type Command interface {
	Print()
	GetType() string
	Execute() error
}

// ##==================================================================
func NewCommand() (Command, error) {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Run 'gtml help' for usage.")
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
	fmt.Println("Run 'gtml help' for usage.")
	return nil, nil
}

// ##==================================================================
type CommandBuild struct {
	Type         string
	InputDir     string
	OutputDir    string
	Options      []Option
	FilteredArgs []string
}

func NewCommandBuild(opts []Option) (*CommandBuild, error) {
	cmd := &CommandBuild{
		Type:    KeyCommandBuild,
		Options: opts,
	}
	err := fungi.Process(
		func() error { return cmd.initFilteredArgs() },
	)
	if err != nil {
		return nil, err
	}
	for _, arg := range cmd.FilteredArgs {
		fmt.Println(arg)
	}
	return cmd, nil
}

func (cmd *CommandBuild) Print()          { fmt.Println(cmd.Type) }
func (cmd *CommandBuild) GetType() string { return cmd.Type }
func (cmd *CommandBuild) Execute() error {
	return nil
}

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
		filtered = append(filtered, arg)
	}
	cmd.FilteredArgs = filtered
	return nil
}

// ##==================================================================
type CommandHelp struct {
	Type string
}

func NewCommandHelp() (*CommandHelp, error) {
	cmd := &CommandHelp{
		Type: KeyCommandHelp,
	}
	return cmd, nil
}

func (cmd *CommandHelp) Print()          { fmt.Println(cmd.Type) }
func (cmd *CommandHelp) GetType() string { return cmd.Type }
func (cmd *CommandHelp) Execute() error {
	message := fmt.Sprintf(`
%s

Usage: 
  gtml [OPTIONS]... [INPUT DIR] [OUTPUT DIR]

Example: 
  gtml --watch ./components ./internal

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
