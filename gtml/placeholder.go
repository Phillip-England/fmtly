package gtml

import "fmt"

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
	nameAttr := pointingTo.GetAttr()
	params, err := GetElementParams(pointingTo)
	if err != nil {
		return nil, err
	}
	fmt.Println(params)
	place := &PlaceholderComponent{
		Name:       nameAttr,
		FoundAs:    foundAsHtml,
		PointingTo: pointingTo,
	}
	return place, nil
}

// ##==================================================================
type PlaceholderComponent struct {
	Name       string
	FoundAs    string
	PointingTo Element
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
