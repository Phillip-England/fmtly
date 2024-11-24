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

# Features for v0.1.0
1. For Loops ✅
2. Conditionals ✅
3. Slots ✅
4. Internal Placeholders ✅
5. Root Placeholders ✅
6. Automatic Attribute Organizing ✅
7. Basic Props ✅
8. Command Line Tool For Generation
9. A --watch Command
10. Type Generation
11. Var Name Randomization
12. Remove _ Attributes in Output
13. Solid README.md
14. Managing Imports and Package Names in Output File
15. Tests For Many Multilayered Components
