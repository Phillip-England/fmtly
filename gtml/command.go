package gtml

import (
	"fmt"
	"os"

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
 -------------------------------
 An HTML to Golang Transpiler 💦
 Version 0.1.0 (2024-11-26)
 https://github.com/phillip-england/gtml
 -------------------------------`

	return purse.RemoveFirstLine(art)
}

// ##==================================================================
type Command interface {
	Print()
	Execute() error
}

// ##==================================================================
func NewCommand() (Command, error) {
	args := os.Args
	if len(args) == 1 {
		cmd, err := NewCommandHelp(KeyCommandHelp)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
	rootArg := args[1]
	if rootArg == KeyCommandBuild {
		cmd, err := NewCommandBuild(KeyCommandBuild)
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}
	cmd, err := NewCommandHelp(KeyCommandHelp)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// ##==================================================================
type CommandBuild struct {
	Args  []string
	Arg   string
	Type  string
	Emoji string
}

func NewCommandBuild(arg string) (*CommandBuild, error) {
	cmd := &CommandBuild{
		Arg:   arg,
		Type:  KeyCommandBuild,
		Emoji: "⚡",
	}
	return cmd, nil
}

func (cmd *CommandBuild) Print() { fmt.Println(cmd.Arg) }
func (cmd *CommandBuild) Execute() error {
	return nil
}

// ##==================================================================
type CommandHelp struct {
	Args  []string
	Arg   string
	Type  string
	Emoji string
}

func NewCommandHelp(arg string) (*CommandHelp, error) {
	cmd := &CommandHelp{
		Arg:   arg,
		Type:  KeyCommandHelp,
		Emoji: "🆘",
	}
	return cmd, nil
}

func (cmd *CommandHelp) Print() { fmt.Println(cmd.Arg) }
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
