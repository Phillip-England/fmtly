package main

import (
	"fmt"
	"gtml/src/cli"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

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
