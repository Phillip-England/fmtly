package tokenize

import "fmtly/internal/comp"

func Components(targetDir string) error {

	err := comp.ReadDir(targetDir)
	if err != nil {
		return err
	}

	return nil
}
