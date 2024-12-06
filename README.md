# GTML
Make Writing HTML in Go a Breeze 🍃

## What is GTML?
GTML is a compiler which converts `.html` files into composable `.go` functions. Think of it like [JSX](https://react.dev/learn/writing-markup-with-jsx) for go.

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
With `go 1.22.3` or later, clone the repo and build the binary on your system.

```bash
git clone https://github.com/phillip-england/gtml;
cd gtml;
go build -o gtml main.go;
```

Then you'll be left with a binary you can move onto your PATH.
```bash
mv gtml ./some/dir/on/your/path
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
 Make Writing HTML in Go a Breeze 🍃
 Version 0.1.9 (2024-12-6)
 https://github.com/phillip-england/gtml
 ---------------------------------------

Usage: 
  gtml [OPTIONS]... [INPUT DIR] [OUTPUT FILE] [PACKAGE NAME]

Example: 
  gtml --watch build ./components output.go output

Options:
  --watch       rebuild when source files are modified

```

## Attributes Define Structure
In gtml, we make use of html attributes to determine a components structure. Here is a quick list of the available attributes:

- _component 
- _for
- _if
- _else
- _slot
- _md

## _component
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

## _for
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

We can also do the same with a slice of custom types:
```html
<div _component="GuestList">
    <ul _for='guest of Guests []Guest'>
        <p>$val(guest.Name)</p>
    </ul>
</div>
```

> 🚨 bring your own types, gtml will not autogenerate them (yet..? 🦄)

## _if
`_if` elements are used to render a piece of html if a condition is met.

input:
```html
<div _component="AdminPage">
    <div _if="isLoggedIn">
        <p>you are logged in!</p>
    </div>
</div>
```

## _else
`_else` elements are used to render a piece of html if a condition is not met.

input:
```html
<div _component="AdminPage">
    <div _else="isLoggedIn">
        <p>you are not logged in!</p>
    </div>
</div>
```

## _slot
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

## _md
`_md` elements are used to render a markdown file into html. You can also provide a theme in `_md-theme`. [Here](https://github.com/alecthomas/chroma/tree/master/styles) is a list of the available themes.

gtml uses [goldmark](https://github.com/yuin/goldmark-highlighting) under the hood to parse `.md` files. 

input:
```html
<div _component="BlogPost">
    <div _md="/path/to/file.md" _md-theme="dracula"></div>
</div>
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

`$val()` is also used to access the data of iteration items in `_for` elements.

For example:
```html
<div _component="GuestList">
    <ul _for='guest of Guests []Guest'>
        <p>$val(guest.Name)</p>
    </ul>
</div>
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

## Dev Notes
This section contains notes related to the ongoing development of gtml.

# Feature Wish List (v0.2.0)
- Solid Error Handling
- _component validations ran prior to building
- implement the $ctx() rune - stores a value in a global context which is made available to children and avoids the used of $pipe()
- implement the $var() rune - creates a local variable (meaning it cannot be used in $pipe())
- $md() rune support - enable the ability to inline markdown content into components
- _components cannot be named a traditional html tag name ✅
- required _components to have a name ✅
- _components cannot have the same name ✅

# Feature Wish List (v0.3.0)
- JSX <SingleTag/> support (preprocessing required)
- camelCase Supported in Attributes (preprocessing required)
- Type Generation (feels more like a luxery feature?)
- Output Cleanup (again, luxery?)
- allow the command line tool to take in a single file instead of a dir as well (not vital)
- in this [reddit convo](https://www.reddit.com/r/golang/comments/1h1yb4w/gtml_convert_html_to_golang/), I talk to someone about growing out the buffers in the output components, need to do this!

# Error Handling Todos
- What if we place an invalid rune into one of our attributes?
- take into consideration which funcs are private / public
- how to catch and report free-floating runes?
