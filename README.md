# Gtml
Convert HTML to Golang 💦

## Hello, World
Turn this:
```html
<div _component='Greeting'>
    <h1>Hello, $prop("name")</h1>
</div>
```

Into this:
```go
func Greeting(name string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component='Greeting'><h1>Hello, `)
    builder.WriteString(name)
    builder.WriteString(`!</h1>`)
    return builder.String()
}
```

## Installation
To install, simply clone the repo and build the binary on your system.

```bash
git clone https://github.com/phillip-england/gtml
cd gtml
go build gtml ## go version 1.22.3 or later
```

Then you'll be left with a binary you can move onto your PATH.

## Inspirations
Before you dive in an check out the features of gtml, I want to take a moment to give thanks to the technologies which have directly inspired this project.

### Templ
Most notable, my project is inspired by [Templ](https://templ.guide). Templ really challenged me to start looking into code generation. I don't think gtml would be what it is without Templ.

### Svelte
When I first dived into gtml, I was using `{{}}` everywhere for my syntax. When thinking of clean alternatives, the rune system from [Svelte](https://svelte.dev/) came to mind. The syntax used in my rune system is inspired by Svelte.

### HTMX
When I was thinking about how I wanted my templating directives to be declared, I thought of [HTMX](https://htmx.org/) and the way it uses attributes to define behaviour. I really liked the idea of using attributes as a way to tell gtml about how the component should be structured. I got this idea from HTMX.

### Tailwind
[Tailwind](https://tailwindcss.com/docs/installation) really challenged me on what was possible from simple markup. Tailwind made me ask the question, "What if there is more to html?" Without Tailwind, I may very well still be apprehensive about adding more to my html.

## Attributes Define Structure
In gtml, we make use of html attributes to determine a components structure. Here is a quick list of the available attributes:

- _component 
- _for
- _if
- _else
- _slot

### _component


### _for

### _if

### _else

### _slot



# Runes

## $prop()
`$prop()` is used to define a `prop` within our `_component`. The value passed into `$prop()` will end up in the function arguments of our output component.

🚨: `$prop()` only accepts strings: `$prop("someStr")`

Once a `$prop()` has been defined, it can used in elsewhere in the same component using `$val()`

Also, you can pipe the value of a `$prop()` into a child `_component` using `$pipe()`


## $val()
`$val()` is used to access the value of another prop within the component.

For example, in a `_for` element, `$val()` may access the value of the slice we are looping over using `$val(item.Property)`.

If a value has been defined using `$prop()`, you may access the value of the `$prop()` using `$val(propName)`





## Dev Notes
This section contains notes related to the ongoing development of gtml.

# Feature Wish List (v0.2.0)
- JSX <SingleTag/> support
- camelCase Supported in Attributes
- Type Generation (each component to have it's own ComponentNameProps type to match)
- Output Cleanup
- Solid Error Handling
- _component validations ran prior to building
- allow the command line tool to take in a single file instead of a dir as well


# Error Handling Todos
- If two components have the same name, throw an error
- What if we place an invalid rune into one of our attributes?

