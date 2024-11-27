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
INSTALL INSTRUCTIONS HERE

## Inspirations
Before you dive in an check out the features of gtml, I want to take a moment to give thanks to the technologies which have directly inspired this project.

### Templ
Most notable, my project is inspired by [Templ](https://templ.guide). Templ really challenged me to start looking into code generation. I don't think gtml would be what it is without Templ.

### Svelte
When I first dived into gtml, I was using `{{}}` everywhere for my syntax. When thinking of clean alternatives, the rune system from [Svelte](https://svelte.dev/) came to mind. The syntax used in my rune system is inspired by Svelte.

### HTMX
When I was thinking about how I wanted my templating directives to be declared, I thought of [HTMX](https://htmx.org/) and the way it uses attributes to define behaviour. I really liked the idea of using attributes as a way to tell gtml about how the component should be structured. I got this idea from HTMX.

