package transpile

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func ExtractComps(targetDir string) ([]*Comp, error) {

	compStr, err := getCompStr(targetDir)
	if err != nil {
		return nil, err
	}

	compStrs, err := makeCompStrs(compStr)
	if err != nil {
		return nil, err
	}

	comps, err := makeComps(compStrs)
	if err != nil {
		return nil, err
	}

	return comps, nil

}

func getCompStr(targetDir string) (string, error) {
	var ss []string
	err := filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		f, err := os.ReadFile(path)
		fStr := string(f)
		ss = append(ss, fStr)
		return nil
	})
	if err != nil {
		return "", err
	}
	return strings.Join(ss, "\n"), nil
}

func makeCompStrs(compStr string) ([]string, error) {
	var comps []string
	lines := strings.Split(compStr, "\n")
	for i, line := range lines {
		sq := strings.ReplaceAll(line, " ", "")
		inComp := false
		if strings.HasPrefix(sq, "<fmt") {
			inComp = true
		}
		if inComp {
			var comp []string
			for i2, line2 := range lines {
				if i2 >= i {
					sq2 := strings.ReplaceAll(line2, " ", "")
					comp = append(comp, line2)
					if strings.HasPrefix(sq2, "</fmt>") {
						comps = append(comps, strings.Join(comp, "\n"))
						break
					}
				}
			}
		}
	}
	return comps, nil
}

func makeComps(compStrs []string) ([]*Comp, error) {
	var comps []*Comp
	for _, str := range compStrs {
		comp, err := NewFmtComp(str)
		if err != nil {
			return comps, err
		}
		comps = append(comps, comp)
	}
	return comps, nil
}
