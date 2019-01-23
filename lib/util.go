package atgen

import (
	"bufio"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func PackageName(fs afero.Fs, gopath string, path string) (string, error) {
	for _, gpath := range filepath.SplitList(gopath) {
		if strings.HasPrefix(path, gpath) {
			return filepath.Rel(filepath.Join(gpath, "src"), path)
		}
	}

	file, err := fs.Open(path)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer file.Close()

	absPath, err := filepath.Abs(file.Name())
	if err != nil {
		return "", errors.WithStack(err)
	}
	gomodPath, err := searchGoModfile(fs, absPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	modfile, err := fs.Open(gomodPath)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer modfile.Close()

	packageName := extractPackageNameFromGomodfile(modfile)
	if gomodPath == "" {
		return "", errors.New("out of GOPATH and not found go.mod")
	}

	packageRoot := filepath.Dir(gomodPath)
	rel, err := filepath.Rel(packageRoot, path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return filepath.Join(packageName, rel), nil
}

func searchGoModfile(fs afero.Fs, absPath string) (string, error) {
	file, err := fs.Open(absPath)
	if err != nil {
		return "", errors.WithStack(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return "", errors.WithStack(err)
	}

	if !fileInfo.IsDir() {
		return "", errors.Errorf("specified file is not dir (absPath=%s)", absPath)
	}

	names, err := file.Readdirnames(-1)
	if err != nil {
		return "", errors.WithStack(err)
	}

	for _, name := range names {
		if name == "go.mod" {
			return filepath.Join(file.Name(), name), nil
		}
	}

	parentPath := filepath.Dir(absPath)
	return searchGoModfile(fs, parentPath)
}

// TODO: Use modfile package after cmd/go/internal/modfile was made public. https://github.com/golang/go/issues/23966
func extractPackageNameFromGomodfile(reader io.Reader) string {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module"))
		}
	}
	return ""
}
