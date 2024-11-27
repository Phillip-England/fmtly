# Gtml
Convert HTML to Golang 💦

## Hello, World
Turn this:
```html
<div _component='Greeting'>
    <h1>Hello, {{ name }}</h1>
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
Gtml was build with go version 1.23.3. 

## Inspirations
Before you dive in an check out the features of gtml, I want to take a moment to give thanks to the technologies which have directly inspired this project.

### Templ
Most notable, my project is inspired by [Templ](https://templ.guide). Templ really challenged me to start looking into code generation. I don't think gtml would be what it is without Templ.

### Svelte
When I first dived into gtml, I was using `{{}}` everywhere for my syntax. When thinking of clean alternatives, the rune system from [Svelte](https://svelte.dev/) came to mind. The syntax used in my rune system is inspired by Svelte.

### HTMX
When I was thinking about how I wanted my templating directives to be declared, I thought of [HTMX](https://htmx.org/) and the way it uses attributes to define behaviour. I really liked the idea of using attributes as a way to tell gtml about how the component should be structured. I got this idea from HTMX.






# Rules Noted
- all HTML attribute names must be written in kebab casing while attribute values may be camel case
- when declaring a prop using {{ propName }} syntax, you must use camel casing to define the name
- use @ to pipe props into an child Elements
- use @ to pipe props into a child Element
- use {{ propName }} within an attribute to define a prop as well

# Feature Wish List (v0.2.0)
- JSX <SingleTag/> support
- camelCase Supported in Attributes
- Type Generation (each component to have it's own ComponentNameProps type to match)
- Output Cleanup
- Solid Error Handling
- _component validations ran prior to building
- Implement the Rune Idea from Below

# Vision Of Changes (v0.2.0)
Usage of `{{}}` and it's multiple use-cases is odd, for example, within an attribute, using `{{}}` will define a param in the output func. It also represents a kind of placeholder because the param name itself will be in the output, just without the `{{}}`. BUT if we use `{{}}` in the context of a _for loop, then it has a different meaning. It is also used to define slots. It just feels like {{}} is wearing a lot of different hats. I like the idea of splitting these functionalities across a series of runes (similar to Svelte but with different implementation). For example, we could introduce the `$prop("propName")` rune which is very direct in it's desire to define a prop and also represents the prop as a placeholder. Maybe if we wanted to take things further we could say, `$pipe("propName")` to grab `propName` from the current context and pipe it into another component. This would also eliminate the weird `@` syntax. We could do `$slot("slotName")` to eliminate our issue with slots. In _for loops, we could do `$val("guest.Name")` or `$val("color")` and it will be easy to tell which input is a string of a specific type. This way we get rid of `{{}}` completely and then we will be in the position to consider which runes can be used to inject client side interactivity and a bunch of other things 🦄
