package main

import (
	"fmtly/internal/transpile"
)

func main() {

	_, err := transpile.CompToGo("./components", "./output.go")
	if err != nil {
		panic(err)
	}

}
