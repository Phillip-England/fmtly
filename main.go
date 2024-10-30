package main

import (
	"fmtly/internal/transpile"
)

func main() {

	_, err := transpile.ExtractComps("./components")
	if err != nil {
		panic(err)
	}

}
