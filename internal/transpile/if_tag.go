package transpile

type IfTag struct {
	Text          string
	ConditionAttr string
	TagAttr       string
	IsClosingTag  bool
	LineNunber    int
}

func NewIfTag(line *CompLine, tag *CompTag, comp *Comp) (*IfTag, error) {
	conditionAttr, _ := tag.GetAttr("condition")
	tagAttr, _ := tag.GetAttr("tag")
	ifTag := &IfTag{
		Text:          line.Html,
		ConditionAttr: conditionAttr,
		TagAttr:       tagAttr,
		IsClosingTag:  tag.IsCloseTag,
		LineNunber:    line.Number,
	}
	return ifTag, nil
}
