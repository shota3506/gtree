package entry

import (
	"path/filepath"
)

type Entry interface {
	Name() string
	Path() string
	IsDir() bool
	Size() int
	Depth() int
}

func NewRoot(path string, showHidden bool) (*Dir, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	base := filepath.Base(absPath)

	d := &Dir{
		name:       base,
		path:       absPath,
		showHidden: showHidden,
	}
	if err := d.Open(); err != nil {
		return nil, err
	}
	return d, nil
}
