# gtml
Html in Go doesn't *have* to suck. Introducing... 🥁 gtml.

## Hello, World

Turn this:
```html
<div _component='HelloWorld'>
    <h1>Hello, World!</h1>
</div>
```

Into this:
```go
func HelloWorld() string {
    var builder strings.Builder
    builder.WriteString(`<div _component='GreetingCard'><h1>Hello, World!</h1>`)
    return builder.String()
}
```

## Attribute Listing
gtml uses html attributes to determine the structure of our components. Here is a list of the available attributes:

### _component
A _component element is the root-level element for a gtml component. A _component element must have no parents and it must be named using PascalCasing. A _component element may not have any _component elements within it.

Here is a basic _component element making use of a prop (more on props later):
```html
<div _component='Greeting'>
    <h1>Hello, {{ name }}</h1>
</div>
```

gtml will scan all the `.html` files in a given directory, generating a Go function for each _component element. The above component will resolve to:
```go
func Greeting(name string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component='Greeting'><h1>Hello, `)
    builder.WriteString(name)
    builder.WriteString(`!</h1>`)
    return builder.String()
}
```

# Features for v0.1.0
- For Loops (bug where props inside loops wont register in param str)
- Conditionals ✅
- Slots ✅
- Internal Placeholders ✅
- Root Placeholders ✅
- Automatic Attribute Organizing ✅
- Basic Props ✅
- Command Line Tool For Generation
- A --watch Command
- Type Generation
- Solid README.md
- Managing Imports and Package Names in Output File
- Tests For Many Multilayered Components
- Attributes can use props ✅

# Rules Noted
- all HTML attribute names must be written in kebab casing while attribute values may be camel case
- when declaring a prop using {{ propName }} syntax, you must use camel casing to define the name
- use @ to pipe props into an child Elements.
