package main

import (
	"fmtly/internal/comp"
)

func main() {

	_, err := comp.ReadDir("./components")
	if err != nil {
		panic(err)
	}

}

// func CustomerList(title string, customers []*Customer) string {
// 	return `
//         <fmt name="CustomerList" tag="ul">
//             <h2>{{ title }}</h2>
//             <for in="customers" type="Customer" tag="li">
//                 <p>{{ customer.Name }}</p>
//                 <for in="customer.Friends" tag="div">
//                     <p>{{ friend.Name }}</p>
//                     <p>{{ friend.Age }}</p>
//                 </for>
//             </for>
//             <if condition="isLoggedIn" tag="div">
//                 <p>logged in</p>
//                 <else>
//                     <p>not logged in</p>
//                 </else>
//             </if>
//         </fmt>
//     `
// }
