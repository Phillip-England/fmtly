# fmtly

Take HTML files and convert them into easy-to-use Go functions.

## Order of Operations

Here, I outline a Component which includes all of the current `fmtly` features.
This should give you a solid overview of how fmtly works.

Take the following component:
```html
<define fmt="CustomerList" tag="ul">
    <h1>{{ listTitle }}</h1>
    <for in="customers" type="[]*Customer">
        <p>{{ customer.Name }}</p>
        <for in="customer.Friends">
            <p>{{ friend.Name }}</p>
            <p>{{ friend.Age }}</p>
        </for>
        <if condition='isSubExpired'>
            <p>Customer subscription is not expired</p>
        </if>
        <if condition="!isSubExpired">
            <p>Customer subscription is expired</p>
        </if>
    </for>
</define>
```
