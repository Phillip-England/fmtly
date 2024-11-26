package main

import "gtml/gtml"

func main() {
	cmd, err := gtml.NewCommand()
	if err != nil {
		panic(err)
	}
	err = cmd.Execute()
	if err != nil {
		panic(err)
	}
}
