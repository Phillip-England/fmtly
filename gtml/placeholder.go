package gtml

import "fmt"

// ##==================================================================
type Placeholder interface {
	Print()
	GetElement() Element
}

// ##==================================================================
type PlaceholderComponent struct {
	Name    string
	FoundAs string
	Element Element
}

func NewPlaceholder(name string, foundAs string, elm Element) Placeholder {
	place := &PlaceholderComponent{
		Name:    name,
		FoundAs: foundAs,
		Element: elm,
	}
	return place
}

func (place *PlaceholderComponent) Print() {
	fmt.Println(fmt.Sprintf(`Name: %s
FoundAs: %s
Element: %s`, place.Name, place.FoundAs, place.Element.GetHtml()))
}
func (place *PlaceholderComponent) GetElement() Element { return place.Element }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
