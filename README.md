# gtml
Let's be reel 🎣, html in Go sucks. But it doesn't *have* to. Introducing... 🥁 gtml.
![huh](./static/huh.webp)

## _component
_component elements are used to register a new component. Gtml will convert all _component elements into their own corrosponding Go function.

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