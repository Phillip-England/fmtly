package main

import (
	"fmt"
	"gtml/src/cli"
	"os"
)

func main() {

	cmd, err := cli.NewCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	if cmd == nil {
		return
	}

	ex, err := cli.NewExecutor(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	err = ex.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

}
