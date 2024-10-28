# fmtly

Use .html files to write .go files 🐈‍⬛

fmtly is an abstraction on top of html to write components in go.

## Component Formatting

Components are expected to follow certain convensions on formatting.
This ensures the output code is consistent across codebases.
It also ensures any output code is formatted consistently as well.

1. All fmtly tags (<fmt>, <if>, <for>, ect) must be on their own line.
2. Tags have requied attributes, as listed in their documentation.

## Fmtly Tags

Tags in your HTML determine serverside actions.
These tags are transpiled into serverside code.

Here is all the tags provided by fmtly:

### For
```html
<fmt name="CustomerList" tag="ul">
    <for in="customers" type="[]*Customer" tag="li">
        <p>{{ customer.Name }}</p>
        <for in="customer.Friends" tag="div">
            <p>{{ friend.Name }}</p>
            <p>{{ friend.Age }}</p>
        </for>
    </for>
</fmt>
```


### If
```html
<fmt name="LightBulb">
    <if condition="isOn">
        <p>light turned on!</p>
        <else>
            <p>light turned off!</p>
        </else>
    </if>
</fmt>
```
