# fmtly

Take HTML files and convert them into easy-to-use Go functions.

## Templating Syntax

Components are server-rendered using familiar directives.
This is a complete list of all the templating directives.

### For
```html
<define fmt="CustomerList" tag="ul">
    <for in="customers" type="[]*Customer" tag="li">
        <p>{{ customer.Name }}</p>
        <for in="customer.Friends" tag="div">
            <p>{{ friend.Name }}</p>
            <p>{{ friend.Age }}</p>
        </for>
    </for>
</define>
```


### If
```html
<define fmt="LightBulb">
    <if condition="isOn">
        <p>light turned on!</p>
        <else>
            <p>light turned off!</p>
        </else>
    </if>
</define>
```
