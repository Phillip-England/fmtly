# gtml
Let's be reel 🎣, html in Go sucks. But it doesn't *have* to. Introducing... 🥁 gtml.

## Html Attributes
gtml allows us to use html attributes to determine the structure of our components. For example, to define a for loop we may do something like this:

```html
<div _component="GuestList">
    <ul _for='guest of guests []Guest'>
        <li>{{ guest.Name }}</li>
    </ul>
</div>
```

gtml will take the component above, and generate the following Go output:

```go
func GuestList(guests []Guest) string {
	var builder strings.Builder
	guestFor := gtmlFor(guests, func(i int, guest Guest) string {
		var guestBuilder strings.Builder
		guestBuilder.WriteString(`<ul _for="guest of guests []Guest"><li>`)
		guestBuilder.WriteString(guest.Name)
		guestBuilder.WriteString(`</li></ul>`)
		return guestBuilder.String()
	})
	builder.WriteString(`<div _component="GuestList">`)
	builder.WriteString(guestFor)
	builder.WriteString(`</div>`)
	return builder.String()
}
```

When using gtml, we simply place our components in a `.html` file, and gtml will take care of generating 


## _component
_component elements are used to register a new component. gtml will convert all _component elements into their own corrosponding Go function.

### input
```html
<div _component='GreetingCard'>
    <h1>Hello, {{ name }}!</h1>
</div>
```

### output
```go
func GreetingCard(name string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component='GreetingCard'><h1>Hello,`)
    builder.WriteString(name)
    builder.WriteString(`!</h1>`)
    return builder.String()
}
```

# Development Notes
1. What if two elements have the same param name?