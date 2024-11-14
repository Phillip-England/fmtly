package gtml_test

import (
	"fmt"
	"gtml"
	"strings"
	"testing"
)

func TestGuestList(t *testing.T) {

	elm, err := gtml.NewGtmlElementFromStr(`
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
			<div>
				<p _if="loggedIn">Im logged in</p>
				<p _else="loggedIn">In not logged in</p>
			</div>
		</div>
	`)
	if err != nil {
		panic(err)
	}

	clay := elm.GetHtml()
	for {

		gtml.WalkUpGtmlBranches(elm, func(child gtml.GtmlElement) error {
			childHtml := child.GetHtml()
			writeStringCall, _ := child.GetWriteStringCall()
			clay = strings.Replace(clay, childHtml, writeStringCall, 1)
			return nil
		})

		elm, err = gtml.NewGtmlElementFromStr(clay)
		if err != nil {
			panic(err)
		}

		if len(elm.GetChildren()) == 0 {
			break
		}

	}

	fmt.Println(clay)

}
