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



# Development Notes
1. What if two elements have the same param name?

<div _component="GreetingCard">
    <h1>{{ name }}</h1>
    <Greeting name="{{ firstGuestName }}" age="20">
        <div _slot="test">
            <ul _for="color of colors []string">
                <li>{{ color }}</li>
            </ul>
            <p>testin!</p>
        </div>
    </Greeting>
</div>

<GreetingCard _component="GreetJackAndBob" firstGuestName="bob" secondGuestName="jack"></GreetingCard>

messageSlot := gtmlSlot(func() string {
    var messageBuilder strings.Builder
    messageBuilder.WriteString(`<div _slot="message"><p>testin!</p></div>`)
    return messageBuilder.String()
})

messageSlot := gtmlSlot(func() string {
    var messageBuilder strings.Builder
    messageBuilder.WriteString(`<div _slot="message"><p>testin!</p></div>`)
    return messageBuilder.String()
})  

loopSlot := gtmlSlot(func() string {
        var loopBuilder strings.Builder
        colorFor := gtmlFor(colors, func(i int, color string) string {
                var colorBuilder strings.Builder
                colorBuilder.WriteString(`<ul _for="color of colors []string"><li>`)
                colorBuilder.WriteString(color)
                colorBuilder.WriteString(`</li></ul>`)
                return colorBuilder.String()
        })
        loopBuilder.WriteString(`<div _slot="loop">`)
        loopBuilder.WriteString(colorFor)
        loopBuilder.WriteString(`</div>`)
        return loopBuilder.String()
})

 loopSlot := gtmlSlot(func() string {
        var loopBuilder strings.Builder
        colorFor := gtmlFor(colors, func(i int, color string) string {
                var colorBuilder strings.Builder
                colorBuilder.WriteString(`<ul _for="color of colors []string"><li>`)
                colorBuilder.WriteString(color)
                colorBuilder.WriteString(`</li></ul>`)
                return colorBuilder.String()
        })
        loopBuilder.WriteString(`<div _slot="loop">`)
        loopBuilder.WriteString(colorFor)
        loopBuilder.WriteString(`</div>`)
        return loopBuilder.String()
})



