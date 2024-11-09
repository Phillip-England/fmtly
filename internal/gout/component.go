package gout

import "tagly/internal/tag"

type ComponentFunc struct {
}

func NewComponentFuncFromFmtTag(t tag.FmtTag) (ComponentFunc, error) {
	cf := &ComponentFunc{}
	return *cf, nil
}
