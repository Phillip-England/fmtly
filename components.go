package main

func CustomerList(title string, customers []*Customer, animals []string, isLoggedIn bool) string {
	return `<ul name="CustomerList" class="p-4 bg-black" directive="fmt"><h2>` + title + `</h2>` + collectStr(customers, func(i int, customer *Customer) string {
		return `<li><p>` + customer.Name + `</p>` + collectStr(customer.Friends, func(i int, friend *Friend) string {
			return `<div><p>` + friend.Name + `</p><p>` + friend.Age + `</p></div>`
		}) + `</li>`
	}) + `` + ifElse(isLoggedIn, `<div>`+collectStr(animals, func(i int, animal string) string { return `<div><p>` + animal + `</p></div>` })+`<p>logged in</p></div>`, `<div><p>not logged in</p></div>`) + `</ul>`
}
