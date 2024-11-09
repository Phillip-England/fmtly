package filetype

import "tagly/internal/tag"

type TaglyFile struct {
	Path    string
	FmtTags []tag.FmtTag
}

func NewTaglyFileFromFilePath(path string) (TaglyFile, error) {
	fmtTags, err := tag.NewFmtTagsFromFilePath(path)
	if err != nil {
		return TaglyFile{}, err
	}
	return TaglyFile{
		FmtTags: fmtTags,
	}, nil
}
