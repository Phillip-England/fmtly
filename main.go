package main

import (
	"fmt"
	"gtml/internal/gtml"
)

func main() {

	elm, err := gtml.NewElement(`
		<div _component="GuestList">
			<div _for="guest of guests []Guest">
				<h1>{{ guest.Name }}</h1>
				<p>The guest has brought the following items:</p>
				<div _for="item of guest.Items []Item">
					<p>{{ item.Name }}</p>
					<p>{{ item.Price }}</p>
				</div>
			</div>
			<div _for="color of colors []string">
				<p>{{ color }}</p>
				<p>{{ color }}</p>
			</div>
		</div>
	`, nil)
	if err != nil {
		panic(err)
	}

	for _, child := range elm.GetChildren() {
		child.DeleteSelf()
	}

	fmt.Println(len(elm.GetChildren()))

}
