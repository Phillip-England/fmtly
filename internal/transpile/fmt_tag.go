package transpile

type FmtTag struct {
	Text     string
	NameAttr string
	TagAttr  string
}

func NewFmtTag(line *CompLine, tag *CompTag, comp *Comp) (*FmtTag, error) {
	nameAttr, _ := tag.GetAttr("name")
	tagAttr, _ := tag.GetAttr("tag")
	fmtTag := &FmtTag{
		Text:     line.Html,
		NameAttr: nameAttr,
		TagAttr:  tagAttr,
	}
	return fmtTag, nil
}
