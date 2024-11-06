package tag

import "github.com/PuerkitoBio/goquery"

func CountStepsUntilParent(count int, child *goquery.Selection, parent *goquery.Selection) (int, bool, error) {
	nextParent := child.Parent()
	if nextParent.Length() == 0 {
		return 0, false, nil
	}
	equal := true
	parent.Each(func(i int, s1 *goquery.Selection) {
		s2 := nextParent.Eq(i)
		if s1.Nodes[0] != s2.Nodes[0] {
			equal = false
		}
	})
	if equal {
		return count, true, nil
	}
	count++
	return CountStepsUntilParent(count, nextParent, parent)
}
