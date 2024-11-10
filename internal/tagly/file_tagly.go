package tagly

type FileTagly struct {
	Path    string
	TagFmts []TagFmt
}

func NewFileTaglyFromFilePath(path string) (FileTagly, error) {
	fmtTags, err := NewTagFmtsFromFilePath(path)
	if err != nil {
		return FileTagly{}, err
	}
	return FileTagly{
		TagFmts: fmtTags,
	}, nil
}
