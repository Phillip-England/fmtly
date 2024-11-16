package gtml

import (
	"gtml/internal/gqpp"
	"strings"
	"testing"
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

	root, err := NewElement(sel)
	if err != nil {
		panic(err)
	}

	fn, err := NewGoFunc(root)
	if err != nil {
		panic(err)
	}

	PrintGoFunc(fn)

}

func GuestList(color string, guests []Guest, colors []string) string {
	var builder strings.Builder
	guestLoop := collect(guests, func(i int, guest Guest) string {
		var guestBuilder strings.Builder
		itemLoop := collect(guest.Items, func(i int, item Item) string {
			var itemBuilder strings.Builder
			colorLoop := collect(item.Colors, func(i int, color Color) string {
				var colorBuilder strings.Builder
				colorBuilder.WriteString(`<div _for="color of item.Colors []Color"><p>`)
				colorBuilder.WriteString(color.Shade)
				colorBuilder.WriteString(`</p><p>`)
				colorBuilder.WriteString(color.Hue)
				colorBuilder.WriteString(`</p></div>`)
				return colorBuilder.String()
			})
			itemBuilder.WriteString(`<div _for="item of guest.Items []Item"><p>`)
			itemBuilder.WriteString(item.Name)
			itemBuilder.WriteString(`</p><p>`)
			itemBuilder.WriteString(item.Price)
			itemBuilder.WriteString(`</p>`)
			itemBuilder.WriteString(colorLoop)
			itemBuilder.WriteString(`</div>`)
			return itemBuilder.String()
		})
		guestBuilder.WriteString(`<div _for="guest of guests []Guest"><h1>`)
		guestBuilder.WriteString(guest.Name)
		guestBuilder.WriteString(`</h1><p>The guest has brought the following items:</p>`)
		guestBuilder.WriteString(itemLoop)
		guestBuilder.WriteString(`</div>`)
		return guestBuilder.String()
	})
	colorLoop := collect(colors, func(i int, color string) string {
		var colorBuilder strings.Builder
		colorBuilder.WriteString(`<div _for="color of colors []string"><p>`)
		colorBuilder.WriteString(color)
		colorBuilder.WriteString(`</p><p>`)
		colorBuilder.WriteString(color)
		colorBuilder.WriteString(`</p></div>`)
		return colorBuilder.String()
	})
	builder.WriteString(`<div _component="GuestList">`)
	builder.WriteString(guestLoop)
	builder.WriteString(colorLoop)
	builder.WriteString(`</div>`)
	return builder.String()
}
