package tag

type TaglyFile struct {
	Path    string
	FmtTags []FmtTag
}

func NewTaglyFileFromFilePath(path string) (TaglyFile, error) {
	_, err := NewFmtTagsFromFilePath(path)
	if err != nil {
		return TaglyFile{}, err
	}
	return TaglyFile{}, nil
}
