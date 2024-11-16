package gtml

import (
	"gtml/gtml"
	"testing"

	"github.com/phillip-england/gqpp"
)

func TestGtml(t *testing.T) {

	sel, err := gqpp.NewSelectionFromStr(`
		<div _component="GuestList">
			<div _for="guest of guests []Guest">
				<h1>{{ guest.Name }}</h1>
				<p>The guest has brought the following items:</p>
				<div _for="item of guest.Items []Item">
					<p>{{ item.Name }}</p>
					<p>{{ item.Price }}</p>
					<div _for="color of item.Colors []Color">
						<p>{{ color.Shade }}</p>
						<p>{{ color.Hue }}</p>
					</div>
				</div>
			</div>
			<div _for="color of colors []string">
				<p>{{ color }}</p>
				<p>{{ color }}</p>
			</div>
		</div>
	`)
	if err != nil {
		panic(err)
	}

	root, err := gtml.NewElement(sel)
	if err != nil {
		panic(err)
	}

	fn, err := gtml.NewGoFunc(root)
	if err != nil {
		panic(err)
	}

	gtml.PrintGoFunc(fn)

}
