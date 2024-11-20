package gtml

import "fmt"

// ##==================================================================
type Slot interface {
	Print()
}

// ##==================================================================
type SlotComponent struct {
	Name    string
	FoundAs string
	Element Element
}

func NewSlot(name string, foundAs string, elm Element) Slot {
	slot := &SlotComponent{
		Name:    name,
		FoundAs: foundAs,
		Element: elm,
	}
	return slot
}

func (slot *SlotComponent) Print() {
	fmt.Println(fmt.Sprintf(`Name: %s
FoundAs: %s
Element: %s`, slot.Name, slot.FoundAs, slot.Element.GetHtml()))
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
