package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func findModuleRoot(start string) (root string, modulePath string, err error) {
	dir := start
	for {
		fp := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(fp); err == nil {
			data, err := os.ReadFile(fp)
			if err != nil {
				return "", "", err
			}

			f, err := modfile.Parse(fp, data, nil)
			if err != nil {
				return "", "", err
			}

			if f.Module == nil {
				return "", "", fmt.Errorf("no module directive in %s", fp)
			}

			return dir, f.Module.Mod.Path, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}

func packageImportPath(moduleRoot, modulePath, filePath string) (string, error) {
	dir := filepath.Dir(filePath)
	rel, err := filepath.Rel(moduleRoot, dir)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(rel, "..") {
		return "", fmt.Errorf("file %s is outside module root %s", filePath, moduleRoot)
	}

	rel = filepath.ToSlash(rel)

	if rel == "." {
		return modulePath, nil
	}

	return modulePath + "/" + rel, nil
}
