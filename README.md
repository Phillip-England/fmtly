# Gtml
Convert HTML to Golang 💦

## Hello, World
Turn this:
```html
<div _component="Greeting">
    <h1>Hello, $prop("name")</h1>
</div>
```

Into this:
```go
func Greeting(name string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component="Greeting" _id="0"><h1>Hello, `)
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
When gtml is scanning `.html` files, it is searching for `_component` elements. When it finds a `_component`, it will generate a function in go which will output the  `_component`'s html.

> 🚨 `_component` may not be defined within another `_component`. However, you can use a `_component` as a `placeholder` within another `_component`

When defining a `_component`, you must give it a name:
```html
<button _component="CustomButton">Click Me!</button>
```

The above `_component` will be converted into:
```go
func CustomButton() string {
	var builder strings.Builder
	builder.WriteString(`<button _component="CustomButton" _id="0">Click Me!</button>`)
	return builder.String()
}
```

### _for
`_for` elements are used to iterate over a slice. The slice may be a custom type or a string slice. 

`_for` elements require their attribute value to be structured in the following way:
```bash
_for="ITEM OF ITEMS []TYPE"
```

Such as:
```html
<div _component="ColorList">
    <ul _for='color of colors []string'>
        <p>$val(color)</p>
    </ul>
</div>
```

The above component will generate:
```go
func ColorList(colors []string) string {
	var builder strings.Builder
	colorFor1 := gtmlFor(colors, func(i int, color string) string {
		var colorBuilder strings.Builder
		colorBuilder.WriteString(`<ul _for="color of colors []string" _id="1"><p>`)
		colorBuilder.WriteString(color)
		colorBuilder.WriteString(`</p></ul>`)
		return colorBuilder.String()
	})
	builder.WriteString(`<div _component="ColorList" _id="0">`)
	builder.WriteString(colorFor1)
	builder.WriteString(`</div>`)
	return builder.String()
}
```

We can also do the same with a slice of custom types:
```html
<div _component="GuestList">
    <ul _for='guest of Guests []Guest'>
        <p>$val(guest.Name)</p>
    </ul>
</div>
```

Which outputs:
```go
func GuestList(Guests []Guest) string {
	var builder strings.Builder
	guestFor1 := gtmlFor(Guests, func(i int, guest Guest) string {
		var guestBuilder strings.Builder
		guestBuilder.WriteString(`<ul _for="guest of Guests []Guest" _id="1"><p>`)
		guestBuilder.WriteString(guest.Name)
		guestBuilder.WriteString(`</p></ul>`)
		return guestBuilder.String()
	})
	builder.WriteString(`<div _component="GuestList" _id="0">`)
	builder.WriteString(guestFor1)
	builder.WriteString(`</div>`)
	return builder.String()
}
```

> 🚨 bring your own types, gtml will not autogenerate them (yet..? 🦄)

### _if
`_if` elements are used to render a piece of html if a condition is met.

input:
```html
<div _component="AdminPage">
    <div _if="isLoggedIn">
        <p>you are logged in!</p>
    </div>
</div>
```

output:
```go
func AdminPage(isLoggedIn bool) string {
	var builder strings.Builder
	isLoggedInIf1 := gtmlIf(isLoggedIn, func() string {
		var isLoggedInBuilder strings.Builder
		isLoggedInBuilder.WriteString(`<div _if="isLoggedIn" _id="1"><p>you are logged in!</p></div>`)
		if isLoggedIn {
			return isLoggedInBuilder.String()
		}
		return ""
	})
	builder.WriteString(`<div _component="AdminPage" _id="0">`)
	builder.WriteString(isLoggedInIf1)
	builder.WriteString(`</div>`)
	return builder.String()
}
```

### _else
`_else` elements are used to render a piece of html if a condition is not met.

input:
```html
<div _component="AdminPage">
    <div _else="isLoggedIn">
        <p>you are not logged in!</p>
    </div>
</div>
```

output:
```go
func AdminPage(isLoggedIn bool) string {
	var builder strings.Builder
	isLoggedInElse1 := gtmlElse(isLoggedIn, func() string {
		var isLoggedInBuilder strings.Builder
		isLoggedInBuilder.WriteString(`<div _else="isLoggedIn" _id="1"><p>you are not logged in!</p></div>`)
		if !isLoggedIn {
			return isLoggedInBuilder.String()
		}
		return ""
	})
	builder.WriteString(`<div _component="AdminPage" _id="0">`)
	builder.WriteString(isLoggedInElse1)
	builder.WriteString(`</div>`)
	return builder.String()
}
```

### _slot
`_slot` elements are unique in the sense that they are not used within a `_component` itself, rather, they are used in it's `placeholder`.

We will discuss placeholders more in a bit, but for now, just know that a placeholder is what we refer to a `_component` as when it is being used *within another component*.

For example, this `LoginForm` uses `CustomButton` as a `placeholder`
```html
<form _component="LoginForm">
    ...
    <CustomButton></CustomButton>
</form>

<div _component="CustomButton">
    <button>Submit</button>
</div>
```

Now, imagine a scenario where you are going to use certain components on each page, like maybe each page of your website has the same navbar and footer?

That is where `_slot`'s are useful.

For example:
```html
<div _component="GuestLayout">
    <navbar>my navbar</navbar>
    $slot("content")  <============= this is a rune, will discuss them later
    <footer></footer>
<div>

<GuestLayout _component="HomePage">
    <div _slot="content">
        I will appear in the content section!
    </div>
</GuestLayout>
```

outputs:
```go
func GuestLayout(content string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="GuestLayout" _id="0"><navbar>my navbar</navbar>`)
	builder.WriteString(content)
	builder.WriteString(`<footer></footer><div><guestlayout _component="HomePage" _placeholder="GuestLayout" _id="0"><div _slot="content" _id="1">I will appear in the content section!</div></guestlayout></div></div>`)
	return builder.String()
}

func HomePage() string {
	var builder strings.Builder
	guestlayoutPlaceholder0 := func() string {
		contentSlot1 := gtmlSlot(func() string {
			var contentBuilder strings.Builder
			contentBuilder.WriteString(`<div _slot="content" _id="1">I will appear in the content section!</div>`)
			return contentBuilder.String()
		})
		return GuestLayout(contentSlot1)
	}
	builder.WriteString(guestlayoutPlaceholder0())
	return builder.String()
}
```

## Runes Define Data
In gtml, we make use of runes to manage the way data flows throughout our components. Here is a quick list of the available runes in gtml:

- $prop()
- $val()
- $slot()
- $pipe()

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
- make gtml force you to use a valid name for _components
- what if a component element has no name?

