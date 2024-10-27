package fmtly

import (
	"fmt"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestMain(t *testing.T) {

	fStr, err := readDirToStr("./components")
	if err != nil {
		panic(err)
	}

	doc, err := getDoc(fStr)
	if err != nil {
		panic(err)
	}

	var comps []*HTMLComponent
	doc.Find("define").Each(func(i int, s *goquery.Selection) {
		compStr, err := s.Parent().Html()
		if err != nil {
			panic(err)
		}
		comp, err := NewHTMLComponent(compStr, s)
		if err != nil {
			panic(err)
		}
		comps = append(comps, comp)
	})

	for _, comp := range comps {
		fmt.Println(comp)
	}

}
