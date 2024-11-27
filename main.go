package main

import (
	"fmt"
	"gtml/gtml"
	"os"
)

func main() {

	cmd, err := gtml.NewCommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	if cmd == nil {
		return
	}

	ex, err := gtml.NewExecutor(cmd)
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
