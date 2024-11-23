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

	path := "./tests/multiple/placeholder/input.html"

	compNames, err := gtml.ReadComponentElementNamesFromFile(path)
	if err != nil {
		panic(err)
	}

	compElms, err := gtml.ReadComponentElementsFromFile(path, compNames)
	if err != nil {
		panic(err)
	}

	for _, elm := range compElms {

		elm, err = gtml.MarkElementPlaceholders(elm)
		if err != nil {
			panic(err)
		}
		// _, err = gtml.NewFunc(elm)
		// if err != nil {
		// 	panic(err)
		// }

		fn, err := gtml.NewFunc(elm)
		if err != nil {
			panic(err)
		}
		fn.Print()
	}

}

// func GreetingCard(name string, colors []string) string {
// 	var builder strings.Builder
// 	greetingPlaceholder := func() string {
// 		messageSlot := gtmlSlot(func() string {
// 			var messageBuilder strings.Builder
// 			messageBuilder.WriteString(`<div _slot="message"><p>testin!</p></div>`)
// 			return messageBuilder.String()
// 		})
// 		loopSlot := gtmlSlot(func() string {
// 			var loopBuilder strings.Builder
// 			colorFor := gtmlFor(colors, func(i int, color string) string {
// 				var colorBuilder strings.Builder
// 				colorBuilder.WriteString(`<ul _for="color of colors []string"><li>`)
// 				colorBuilder.WriteString(color)
// 				colorBuilder.WriteString(`</li></ul>`)
// 				return colorBuilder.String()
// 			})
// 			loopBuilder.WriteString(`<div _slot="loop">`)
// 			loopBuilder.WriteString(colorFor)
// 			loopBuilder.WriteString(`</div>`)
// 			return loopBuilder.String()
// 		})
// 		return Greeting(firstGuestName, 20, messageSlot, loopSlot)
// 	}
// 	builder.WriteString(`<div _component="GreetingCard"><h1>`)
// 	builder.WriteString(name)
// 	builder.WriteString(`</h1>`)
// 	builder.WriteString(greetingPlaceholder())
// 	builder.WriteString(`</div>`)
// 	return builder.String()
// }
