package static

import (
	"embed"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed wp-content wp-includes
var FsEx embed.FS

type Fs struct {
	embed.FS
	Path string
}

func (f Fs) Open(path string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(path, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}
	fullName := strings.TrimLeft(path, "/")

	fullName = f.Path + "/" + fullName
	file, err := f.FS.Open(fullName)
	return file, err
}
