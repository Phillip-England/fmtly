package gtml

import (
	"fmt"
	"gtml/gtml"
	"os"
	"testing"

	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

func runTestByNameDirName(t *testing.T, testDir string) error {
	inputPath := fmt.Sprintf("./tests/%s/input.html", testDir)
	expectPath := fmt.Sprintf("./tests/%s/expect.txt", testDir)
	f, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	fStr := string(f)
	sel, err := gqpp.NewSelectionFromStr(fStr)
	if err != nil {
		return err
	}
	elm, err := gtml.NewElement(sel)
	if err != nil {
		return err
	}
	fn, err := gtml.NewFunc(elm)
	if err != nil {
		return err
	}
	out := purse.Squeeze(purse.Flatten(fn.GetData()))

	f, err = os.ReadFile(expectPath)
	if err != nil {
		return err
	}
	fStr = string(f)
	expect := purse.Squeeze(purse.Flatten(fStr))

	if out != expect {
		t.Errorf("output does not meet expectations:\n\nexpected:\n\n%s\n\ngot:\n\n%s", fStr, fn.GetData())
	}
	return nil
}

func TestAll(t *testing.T) {
	err := fungi.Process(
		func() error { return runTestByNameDirName(t, "mesh") },
		func() error { return runTestByNameDirName(t, "if") },
		func() error { return runTestByNameDirName(t, "for") },
	)
	if err != nil {
		panic(err)
	}
}

func TestOne(t *testing.T) {

	template := `
        <div _component="UserGate">
            <div _else="isBlocked">
                <p>you are not blocked!</p>
            </div>
        </div>
    `

	sel, err := gqpp.NewSelectionFromStr(template)
	if err != nil {
		panic(err)
	}

	elm, err := gtml.NewElement(sel)
	if err != nil {
		panic(err)
	}

	_, err = gtml.NewFunc(elm)
	if err != nil {
		panic(err)
	}

}
