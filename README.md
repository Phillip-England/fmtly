   _____ _______ __  __ _      
  / ____|__   __|  \/  | |     
 | |  __   | |  | \  / | |     
 | | |_ |  | |  | |\/| | |     
 | |__| |  | |  | |  | | |____ 
  \_____|  |_|  |_|  |_|______|
----------------------------------
Html in Go Doesn't *have* to Suck
----------------------------------

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

### _for
A _for element allows us to iterate over a slice of type [T]. The value of the `_for` attribute must follow this schema:

```html
<div _for="ITEM of ITEMS []TYPE">...</div>
```

Here is a for element. Take note of how we access the underlying type's data using `{{}}` along with `item.Property` syntax:
```html
<div _component="GuestList">
    <p>{{ name }}</p>
    <ul _for='guest of guests []Guest'>
        <li>{{ guest.Name }}</li>
    </ul>
</div>
```

The output:
```go
func GuestList(name string, guests []Guest) string {
    var builder strings.Builder
    guestFor := gtml.For(guests, func(i int, guest Guest) string {
        var guestBuilder strings.Builder
        guestBuilder.WriteString(`<ul _for="guest of guests []Guest"><li>`)
        guestBuilder.WriteString(guest.Name)
        guestBuilder.WriteString(`</li></ul>`)
        return guestBuilder.String()
    })
    builder.WriteString(`<div _component="GuestList"><p>`)
    builder.WriteString(name)
    builder.WriteString(`</p>`)
    builder.WriteString(guestFor)
    builder.WriteString(`</div>`)
    return builder.String()
}
```

You can also choose to iterate over a `[]string` slice. Take note of how we access the strings value using `{{ .strname }}` syntax. Doing this will prevent `{{ color }}` from appearing in the output functions parameter definition.
```html
<div _component="FavoriteColors">
    <ul _for="color of colors []string">
        <li>{{ .color }}</li>
    </ul>
</div>
```

The output:
```go
func FavoriteColors(colors []string) string {
    var builder strings.Builder
    colorFor := gtml.For(colors, func(i int, color string) string {
        var colorBuilder strings.Builder
        colorBuilder.WriteString(`<ul _for="color of colors []string"><li>`)
        colorBuilder.WriteString(color)
        colorBuilder.WriteString(`</li></ul>`)
        return colorBuilder.String()
    })
    builder.WriteString(`<div _component="FavoriteColors">`)
    builder.WriteString(colorFor)
    builder.WriteString(`</div>`)
    return builder.String()
}
```

### _if & _else
_if and _else Elements are opposites. _if will allow a block of html to be rendered if a boolean is true. _else will allow a block to be rendered if a boolean is false.

Here an example:
```html
<div _component="WelcomeBanner">
    <div _if="loggedIn">
        <p>Welcome, User!</p>
    </div>
    <div _else="loggedIn">
        <p>Welcome, Guest!</p>
    </div>
</div>
```

The output:
```go
func WelcomeBanner(loggedIn bool) string {
    var builder strings.Builder
    loggedInIf := gtml.If(loggedIn, func() string {
        var loggedInBuilder strings.Builder
        loggedInBuilder.WriteString(`<div _if="loggedIn"><p>Welcome, User!</p></div>`)
        if loggedIn {
                return loggedInBuilder.String()
        }
        return ""
    })
    loggedInElse := gtml.Else(loggedIn, func() string {
        var loggedInBuilder strings.Builder
        loggedInBuilder.WriteString(`<div _else="loggedIn"><p>Welcome, Guest!</p></div>`)
        if !loggedIn {
                return loggedInBuilder.String()
        }
        return ""
    })
    builder.WriteString(`<div _component="WelcomeBanner">`)
    builder.WriteString(loggedInIf)
    builder.WriteString(loggedInElse)
    builder.WriteString(`</div>`)
    return builder.String()
}
```

### _slot
_slot elements are used to create reusable layout components. In the following example, we have one _component utilizing another _component as a `placeholder` (more on placeholders later):

```html
<div _component="UsingSlots">
    <Sandwich>
        <div _slot="top">
            <p>I'm on top!</p>
        </div>
        <div _slot="bottom">
            <p>I'm on bottom</p>
        </div>
    </Sandwich>
</div>


<div _component="Sandwich">
    {{ slot top }}
    <p>🥪</p>
    {{ slot bottom }}
</div>
```

The output:
```go
func UsingSlots() string {
    var builder strings.Builder
    sandwichPlaceholder := func() string {
        topSlot := gtml.Slot(func() string {
            var topBuilder strings.Builder
            topBuilder.WriteString(`<div _slot="top"><p>I&#39;m on top!</p></div>`)
            return topBuilder.String()
        })
        bottomSlot := gtml.Slot(func() string {
            var bottomBuilder strings.Builder
            bottomBuilder.WriteString(`<div _slot="bottom"><p>I&#39;m on bottom</p></div>`)
            return bottomBuilder.String()
        })
        return Sandwich(topSlot, bottomSlot)
    }
    builder.WriteString(`<div _component="UsingSlots">`)
    builder.WriteString(sandwichPlaceholder())
    builder.WriteString(`</div>`)
    return builder.String()
}
func Sandwich(top string, bottom string) string {
    var builder strings.Builder
    builder.WriteString(`<div _component="Sandwich">`)
    builder.WriteString(top)
    builder.WriteString(`<p>🥪</p>`)
    builder.WriteString(bottom)
    builder.WriteString(`</div>`)
    return builder.String()
}
```
## Props
Props are used to inform gtml about dynamic data within our components, along with a few other things. They are defined using `{{}}` syntax.

This area of the project is still under work. I like the `$rune()` syntax Svelte uses, and I might adopt it to create a distinction between `{{ props }}` and `$runes()`.

🚨 All `prop` names **must** be camelCase.

Here are the different types of props which can be used in your components:

### String Props
String props are the most common type of prop you'll encounter. They are used to define a piece of dynamic data within our component. Here is an example:

```html
<div _component='Greeting'>
    <h1>Hello, {{ name }}</h1>
</div>
```

In the above example, name will appear in our output function like so:

```go
func Greeting(name string) {...}
```

### For Type Props
For Type props are used when we are iterating over a slice using `_for`. A For Type Prop is identified by it's use of a `.` like `{{ user.Email }}`.

These types of props will not appear in our output function. Rather, they are used to let gtml know we are trying to access an types internal data.

Here is an example:
```html
<div _component="GuestList">
    <ul _for='guest of guests []Guest'>
        <li>{{ guest.Name }}</li>
    </ul>
</div>
```

### For Str Props
For Str Props are very similar to For Type Props, except they are used when iterating over a []string slice. A For Str Prop is identified by it's use of `.` at the start of the prop like `{{ .color }}`.

These types of props must match the name provided by the `_for` attribute they are refferencing:

```html
<div _component="FavoriteColors">
    <ul _for="color of colors []string">
        <li>{{ .color }}</li>
    </ul>
</div>
```

### Slot Props
A Slot Prop is used to tell gtml where to place children within a `placeholder`. Here is an example:

```html
<div _component="Sandwich">
    {{ slot top }}
    <p>🥪</p>
    {{ slot bottom }}
</div>
```

Then, we can make use of these slots when we use `Sandwich` as a `placeholder`
```html
<div _component="UsingSlots">
    <Sandwich>
        <div _slot="top">
            <p>I'm on top!</p>
        </div>
        <div _slot="bottom">
            <p>I'm on bottom</p>
        </div>
    </Sandwich>
</div>
```

## Placeholders
When a `_component` is used inside of another `_component`, it is reffered to as a `placeholder`. Here is an example of us making use of a simple placeholder:

```html
<form _component="LoginForm">
    <h1>Login</h1>
    <input type="text" name="username"/>
    <SubmitButton></SubmitButton>
</form>

<button _component="SubmitButton">Submit</button>
```

### Not JSX
Take note, this is not JSX, we cannot do this (yet?):
```html
<form _component="LoginForm">
    <h1>Login</h1>
    <input type="text" name="username"/>
    <SubmitButton />
</form>

<button _component="SubmitButton">Submit</button>
```

### Kebab Casing on Attributes
Placeholders may need data to be rendered properly, like this one:
```html
<div _component="FancyPost">
    <FancyText some-text="I am Fancy!"></FancyText>
</div>

<div _component="FancyText">{{ someText }}</div>
```

🚨 When passing data into a `placeholder`, you must use kebab-casing. This can be patched in future releases.

### Piping Data into Placeholders using "@"
Placeholders may need to pipe data from a parent component into it's own context. To do so, we can use `@`:

```html
<div _component="SurvivalTeam">
    <ul _for="member of members []Memeber">
        <Survivalist member-name="@member.Name" member-age=@member.age></Survivalist>
    </ul>
</div>

<div _component="SurvivalList">
    <h1>{{ memberName }}</h1>
    <p>{{ memberAge }}</p>
</div>
```

Without the `@` operator, gtml would treat piped in values as strings instead.

### Placeholders and Props
What if we have a `placeholder`, and we want the value we are passing into it to be dynamic? Well, just use a prop:

```html
<div _component="GreetingCard">
    <Greeting name="{{ name }}" age="20"></Greeting>
</div>

<div _component="Greeting">
    <h1>Hello, {{ name }}</h1>
    <p>you are {{ age }} years old!</p>
</div>
```

Now, `GreetingCard()` will expect a name, which will be piped into the context of `Greeting()`

# Features for v0.1.0
- For Loops ✅
- Conditionals ✅
- Slots ✅
- Internal Placeholders ✅
- Root Placeholders ✅
- Automatic Attribute Organizing ✅
- Basic Props ✅
- Command Line Tool For Generation
- A --watch Command
- Solid README.md ✅
- Managing Imports and Package Names in Output File
- Tests For Many Multilayered Components
- Attributes can use props ✅
- Varname Randomization ✅

# Rules Noted
- all HTML attribute names must be written in kebab casing while attribute values may be camel case
- when declaring a prop using {{ propName }} syntax, you must use camel casing to define the name
- use @ to pipe props into a child Element
- use {{ propName }} within an attribute to define a prop as well


# Feature Wish List (v0.2.0)
- JSX <SingleTag/> support
- camelCase Supported in Attributes
- Type Generation (each component to have it's own ComponentNameProps type to match)
- Output Cleanup
- Solid Error Handling
- Implement the Rune Idea from Below

# Vision Of Changes (v0.2.0)
Usage of `{{}}` and it's multiple use-cases is odd, for example, within an attribute, using `{{}}` will define a param in the output func. It also represents a kind of placeholder because the param name itself will be in the output, just without the `{{}}`. BUT if we use `{{}}` in the context of a _for loop, then it has a different meaning. It is also used to define slots. It just feels like {{}} is wearing a lot of different hats. I like the idea of splitting these functionalities across a series of runes (similar to Svelte but with different implementation). For example, we could introduce the `$prop("propName")` rune which is very direct in it's desire to define a prop and also represents the prop as a placeholder. Maybe if we wanted to take things further we could say, `$pipe("propName")` to grab `propName` from the current context and pipe it into another component. This would also eliminate the weird `@` syntax. We could do `$slot("slotName")` to eliminate our issue with slots. In _for loops, we could do `$val("guest.Name")` or `$val("color")` and it will be easy to tell which input is a string of a specific type. This way we get rid of `{{}}` completely and then we will be in the position to consider which runes can be used to inject client side interactivity and a bunch of other things 🦄