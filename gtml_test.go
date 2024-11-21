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

func testSingle(t *testing.T, testDir string) error {
	inputPath := fmt.Sprintf("./tests/single/%s/input.html", testDir)
	expectPath := fmt.Sprintf("./tests/single/%s/expect.txt", testDir)
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
		func() error { return testSingle(t, "mesh") },
		func() error { return testSingle(t, "if") },
		func() error { return testSingle(t, "for") },
		func() error { return testSingle(t, "else") },
	)
	if err != nil {
		panic(err)
	}
}

func TestOne(t *testing.T) {

	// put the placeholder into a prop {{ SubmitButton("Submit!") }}
	// then when we scan each element to get the props
	// the placeholders will show up, get sorted, and then they will
	// end up in our string builder series in the end

	compElms, err := gtml.ReadComponentElementsFromFile("./tests/multiple/placeholder/input.html")
	if err != nil {
		panic(err)
	}

	for _, elm := range compElms {

		placeholders, err := gtml.GetElementPlaceholders(elm, compElms)
		if err != nil {
			panic(err)
		}

		for _, place := range placeholders {
			elm.Print()
			place.Print()
		}

		// _, err = gtml.NewFunc(elm)
		// if err != nil {
		// 	panic(err)
		// }
	}

}
