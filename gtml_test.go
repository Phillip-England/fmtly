package gtml

import (
	"fmt"
	"gtml/gtml"
	"os"
	"testing"

	"github.com/phillip-england/fungi"
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
	elm, err := gtml.NewElement(fStr, []string{})
	if err != nil {
		return err
	}
	err = gtml.SaltElements(elm)
	if err != nil {
		return err
	}
	fn, err := gtml.NewFunc(elm, make([]gtml.Element, 0))
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

func testMultiple(t *testing.T, testDir string) error {
	path := "./tests/multiple/" + testDir + "/input.html"
	compNames, err := gtml.ReadComponentElementNamesFromFile(path)
	if err != nil {
		return err
	}
	compElms, err := gtml.ReadComponentElementsFromFile(path, compNames)
	if err != nil {
		return err
	}
	funcs := make([]gtml.Func, 0)
	for _, elm := range compElms {
		err = gtml.SaltElements(elm)
		if err != nil {
			return err
		}
		elm, err = gtml.MarkElementPlaceholders(elm)
		if err != nil {
			return err
		}
		fn, err := gtml.NewFunc(elm, compElms)
		if err != nil {
			return err
		}
		funcs = append(funcs, fn)
	}
	actual := ""
	for _, fn := range funcs {
		actual += fn.GetData() + "\n"
	}
	expectPath := "./tests/multiple/" + testDir + "/expect.txt"
	expectedF, err := os.ReadFile(expectPath)
	if err != nil {
		return err
	}
	expect := string(expectedF)
	sqActual := purse.Flatten(actual)
	sqExpect := purse.Flatten(expect)
	if sqActual != sqExpect {
		t.Errorf("actual output does not meet expected output:\n\nexpected:\n\n%s\n\ngot:\n\n%s", expect, actual)
	}
	return nil
}

func TestSingles(t *testing.T) {
	err := fungi.Process(
		// func() error { return testSingle(t, "mesh") },
		// func() error { return testSingle(t, "if") },
		// func() error { return testSingle(t, "for") },
		// func() error { return testSingle(t, "for_str") },
		// func() error { return testSingle(t, "else") },
		func() error { return testSingle(t, "if_else") },
	)
	if err != nil {
		panic(err)
	}
}

func TestMultiples(t *testing.T) {
	err := fungi.Process(
	// func() error { return testMultiple(t, "placeholder") },
	// func() error { return testMultiple(t, "placeholder_root") },
	// func() error { return testMultiple(t, "placeholder_root_slot") },
	// func() error { return testMultiple(t, "attribute_prop") },
	// func() error { return testMultiple(t, "loop_with_placeholders") },
	// func() error { return testMultiple(t, "slot") },
	// func() error { return testMultiple(t, "simple_placeholder") },
	// func() error { return testMultiple(t, "placeholder_kebab") },
	// func() error { return testMultiple(t, "placeholder_with_prop") },
	)
	if err != nil {
		panic(err)
	}
}
