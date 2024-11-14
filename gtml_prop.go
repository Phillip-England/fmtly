package gtml

import "gtml/internal/purse"

type GtmlProp struct {
	Raw   string
	Value string
}

func NewPropsFromStr(str string) []GtmlProp {
	props := purse.ScanBetweenSubStrs(str, "{{", "}}")
	out := make([]GtmlProp, 0)
	for _, prop := range props {
		val := purse.RemoveAllSubStr(prop, "{{", "}}")
		val = purse.Flatten(val)
		val = purse.Squeeze(val)
		finalProps := GtmlProp{
			Raw:   prop,
			Value: val,
		}
		out = append(out, finalProps)
	}
	return out
}
