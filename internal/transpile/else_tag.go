package transpile

type ElseTag struct {
	Text         string
	IsClosingTag bool
	LineNunber   int
}

func NewElseTag(line *CompLine, tag *CompTag, comp *Comp) (*ElseTag, error) {
	elseTag := &ElseTag{
		Text:         line.Html,
		IsClosingTag: tag.IsCloseTag,
		LineNunber:   line.Number,
	}
	return elseTag, nil
}
