package transpile

type ForTag struct {
	Text       string
	InAttr     string
	TypeAttr   string
	TagAttr    string
	IsCloseTag bool
	LineNumber int
}

func NewForTag(line *CompLine, tag *CompTag, comp *Comp) (*ForTag, error) {
	inAttr, _ := tag.GetAttr("in")
	typeAttr, _ := tag.GetAttr("type")
	tagAttr, _ := tag.GetAttr("tag")

	forTag := &ForTag{
		Text:       tag.Html,
		InAttr:     inAttr,
		TagAttr:    tagAttr,
		TypeAttr:   typeAttr,
		IsCloseTag: tag.IsCloseTag,
		LineNumber: line.Number,
	}

	return forTag, nil
}
