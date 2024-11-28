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

Then you'll be left with a binary you can move onto your PATH.

```bash
git clone https://github.com/phillip-england/gtml
cd gtml
go build gtml ## go version 1.22.3 or later
```

## Usage
```bash
   _____ _______ __  __ _      
  / ____|__   __|  \/  | |     
 | |  __   | |  | \  / | |     
 | | |_ |  | |  | |\/| | |     
 | |__| |  | |  | |  | | |____ 
  \_____|  |_|  |_|  |_|______|
 ---------------------------------------
 Convert HTML to Golang 💦
 Version 0.1.0 (2024-11-26)
 https://github.com/phillip-england/gtml
 ---------------------------------------

Usage: 
  gtml [OPTIONS]... [INPUT DIR] [OUTPUT FILE]

Example: 
  gtml --watch build ./components output.go output

Options:
  --watch       rebuild when source files are modified

```



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
    $slot("content")
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
In gtml, we make use of `runes` to manage the way data flows throughout our components. Certain `runes` accept string values while others expect raw values. Here is a quick list of the available runes in gtml:

- $prop()
- $val()
- $slot()
- $pipe()


## $prop()
`$prop()` is used to define a `prop` within our `_component`. A `prop` is a value which is usable by sibling and child elements. The value passed into `$prop()` will end up in the arguments of our output function.

> 🚨: `$prop()` only accepts strings: `$prop("someStr")`

For example, we may define a `$prop()` like so:
```html
<div _component="RuneProp">
    <p>Hello, $prop("name")!</p>
</div>
```

The output:
```go
func RuneProp(name string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="RuneProp" _id="0"><p>Hello, `)
	builder.WriteString(name)
	builder.WriteString(`!</p></div>`)
	return builder.String()
}
```

Once a `$prop()` has been defined, it can used in elsewhere in the same component using `$val()`. Also, you can pipe the value of a `$prop()` into a child `_component` using `$pipe()`


## $val()
`$val()` is used to access the value of `prop`.

> 🚨: `$val()` only accepts raw values: `$val(rawValue)`

For example, here we make use of `$val()` to access a neighboring `prop`:
```html
<div _component="Echo">
    <p>$prop("message")</p>
    <p>$val(message)</p>
</div>
```

```go
func Echo(message string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="Echo" _id="0"><p>`)
	builder.WriteString(message)
	builder.WriteString(`</p><p>`)
	builder.WriteString(message)
	builder.WriteString(`</p></div>`)
	return builder.String()
}
```

`$val()` is also used to access the data of iteration items in `_for` elements.

For example:
```html
<div _component="GuestList">
    <ul _for='guest of Guests []Guest'>
        <p>$val(guest.Name)</p>
    </ul>
</div>
```

the output:
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

## $pipe()
`$pipe()` is used in situations where you want to `pipe` data from a parent `_component` into a child `placeholder`.

> 🚨: `$pipe()` only accepts raw values: `$pipe(rawValue)`

For example:
```html
<div _component="RunePipe">
    <p>Sally is $prop("age") years old</p>
    <Greeting age="$pipe(age)"></Greeting> <== piping in the age
</div>

<div _component="Greeting">
    <h1>This age was piped in!</h1> 
    <p>$prop("age")</p>
</div>
```

The output:
```go
func RunePipe(age string) string {
	var builder strings.Builder
	greetingPlaceholder1 := func() string {
		return Greeting(age)
	}
	builder.WriteString(`<div _component="RunePipe" _id="0"><p>Sally is `)
	builder.WriteString(age)
	builder.WriteString(` years old</p>`)
	builder.WriteString(greetingPlaceholder1())
	builder.WriteString(` &lt;== piping in the age</div>`)
	return builder.String()
}

func Greeting(age string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="Greeting" _id="0"><h1>This age was piped in!</h1> <p>`)
	builder.WriteString(age)
	builder.WriteString(`</p></div>`)
	return builder.String()
}
```

## Placeholders
When a `_component` is used within another `_component`, we refer to it as a `placeholder`. `placeholders` enable us to mix and match components with ease.

> 🚨: This is not JSX and you cannot do self-closing tags like `<SomeComponent/>`. All tags must consist of both an opening tag and closing tag. This can be supported in future versions.

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

### Placeholder Attributes
You may pass data into a `placeholder` using it's attributes. These attributes must corrospond to the target `_component`'s `props`. 

> 🚨: when passing values into `placeholder` attributes, you must refer to the attributes in kebab-casing. This will be patched in future version. For example, below we do `$prop("firstName")` to define `firstName`, but when we pass values into the `_element` we use `first-name` instead.

For example:
```html
<div _component="NameTag">
    <h1>$prop("firstName")</h1>
    <p>$prop("message")</p>
</div>

<NameTag _component="PlaceholderWithAttrs" message="is the best" first-name="gtml"></NameTag>
```

The output:
```go
func NameTag(firstName string, message string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="NameTag" _id="0"><h1>`)
	builder.WriteString(firstName)
	builder.WriteString(`</h1><p>`)
	builder.WriteString(message)
	builder.WriteString(`</p></div>`)
	return builder.String()
}

func PlaceholderWithAttrs() string {
	var builder strings.Builder
	nametagPlaceholder0 := func() string {
		return NameTag("gtml", "is the best")
	}
	builder.WriteString(nametagPlaceholder0())
	return builder.String()
}
```

### Placeholder Piping
If a `placeholder` needs to access a value from a parent `_component`, the value may be piped in using the `$pipe()` `rune`.

For example:
For example:
```html
<div _component="RunePipe">
    <p>Sally is $prop("age") years old</p>
    <Greeting age="$pipe(age)"></Greeting> <== piping in the age
</div>

<div _component="Greeting">
    <h1>This age was piped in!</h1> 
    <p>$prop("age")</p>
</div>
```

The output:
```go
func RunePipe(age string) string {
	var builder strings.Builder
	greetingPlaceholder1 := func() string {
		return Greeting(age)
	}
	builder.WriteString(`<div _component="RunePipe" _id="0"><p>Sally is `)
	builder.WriteString(age)
	builder.WriteString(` years old</p>`)
	builder.WriteString(greetingPlaceholder1())
	builder.WriteString(` &lt;== piping in the age</div>`)
	return builder.String()
}

func Greeting(age string) string {
	var builder strings.Builder
	builder.WriteString(`<div _component="Greeting" _id="0"><h1>This age was piped in!</h1> <p>`)
	builder.WriteString(age)
	builder.WriteString(`</p></div>`)
	return builder.String()
}
```

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
- take into consideration which funcs are private / public

