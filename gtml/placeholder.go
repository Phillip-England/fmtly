package gtml

import (
	"fmt"

	"github.com/phillip-england/fungi"
)

// ##==================================================================
type Placeholder interface {
	Print()
	GetFoundAs() string
	GetPointingTo() Element
}

func NewPlaceholder(foundAsHtml string, pointingTo Element) (Placeholder, error) {
	if pointingTo.GetType() != KeyElementComponent {
		return nil, fmt.Errorf("a placeholder must point to a valid _component element: %s", pointingTo.GetHtml())
	}
	place, err := NewPlaceholderComponent(foundAsHtml, pointingTo)
	if err != nil {
		return nil, err
	}
	return place, nil
}

// ##==================================================================
type PlaceholderComponent struct {
	Name       string
	FoundAs    string
	PointingTo Element
}

func NewPlaceholderComponent(foundAsHtml string, pointingTo Element) (*PlaceholderComponent, error) {
	place := &PlaceholderComponent{
		FoundAs:    foundAsHtml,
		PointingTo: pointingTo,
	}
	err := fungi.Process(
		func() error { return place.initName() },
	)
	if err != nil {
		return nil, err
	}
	params, err := GetElementParams(pointingTo)
	if err != nil {
		return nil, err
	}
	for _, param := range params {
		param.Print()
	}
	place.Print()
	return place, nil
}

func (place *PlaceholderComponent) initName() error {
	nameAttr := place.PointingTo.GetAttr()
	place.Name = nameAttr
	return nil
}

func (place *PlaceholderComponent) Print() {
	fmt.Println("Name: " + place.Name)
	fmt.Println("FoundAs: " + place.FoundAs)
	fmt.Print("PointingTo: ")
	place.PointingTo.Print()
}
func (place *PlaceholderComponent) GetFoundAs() string     { return place.FoundAs }
func (place *PlaceholderComponent) GetPointingTo() Element { return place.PointingTo }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
