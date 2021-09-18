package entry

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

type Entry interface {
	fmt.Stringer

	Path() string
	IsDir() bool
	Size() int
	Depth() int
	Walk(fn func(e Entry) error) error
}

type File struct {
	name  string
	path  string
	depth int
}

func (f *File) String() string {
	return f.name
}

func (f *File) Path() string {
	return f.path
}

func (f *File) IsDir() bool {
	return false
}

func (f *File) Size() int {
	return 1
}

func (f *File) Depth() int {
	return f.depth
}

func (f *File) Walk(fn func(e Entry) error) error {
	if err := fn(f); err != nil {
		return err
	}
	return nil
}

type Dir struct {
	name       string
	path       string
	children   []Entry
	depth      int
	open       bool
	showHidden bool
}

func (d *Dir) String() string {
	return d.name
}

func (d *Dir) Path() string {
	return d.path
}

func (d *Dir) IsDir() bool {
	return true
}

func (d *Dir) Size() int {
	if !d.IsOpen() {
		return 1
	}
	var size int
	_ = d.Walk(func(e Entry) error { // never return error
		size += 1
		return nil
	})
	return size
}

func (d *Dir) Depth() int {
	return d.depth
}

func (d *Dir) Open() error {
	if err := d.read(); err != nil {
		return err
	}
	d.open = true
	return nil
}

func (d *Dir) Close() {
	d.open = false
}

func (d *Dir) Toggle() error {
	var err error
	if d.IsOpen() {
		d.Close()
	} else {
		err = d.Open()
	}
	return err
}

func (d *Dir) IsOpen() bool {
	return d.open
}

func (d *Dir) Walk(fn func(e Entry) error) error {
	if err := fn(d); err != nil {
		return err
	}
	if !d.IsOpen() {
		return nil
	}
	for _, c := range d.children {
		if err := c.Walk(fn); err != nil {
			return err
		}
	}
	return nil
}

func (d *Dir) Get(i int) (Entry, error) {
	if i < 0 {
		return nil, errors.New("not found")
	}

	entries := []Entry{}
	_ = d.Walk(func(e Entry) error { // never return error
		entries = append(entries, e)
		return nil
	})
	if i >= len(entries) {
		return nil, errors.New("not found")
	}
	return entries[i], nil
}

func (d *Dir) read() error {
	files, err := os.ReadDir(d.path)
	if err != nil {
		return err
	}

	children := []Entry{}
	for _, f := range files {
		if !d.showHidden && strings.HasPrefix(f.Name(), ".") {
			continue
		}
		children = append(children, d.newEntry(f))
	}

	// sort
	sort.Slice(children, func(i, j int) bool {
		if children[i].IsDir() && !children[j].IsDir() {
			return true
		}
		if !children[i].IsDir() && children[j].IsDir() {
			return false
		}
		return children[i].String() < children[j].String()
	})

	d.children = children
	return nil
}

func (d *Dir) newEntry(f os.DirEntry) Entry {
	var c Entry
	if f.IsDir() {
		c = &Dir{
			name:       f.Name(),
			path:       filepath.Join(d.path, f.Name()),
			depth:      d.depth + 1,
			showHidden: d.showHidden,
		}
	} else {
		c = &File{
			name:  f.Name(),
			path:  filepath.Join(d.path, f.Name()),
			depth: d.depth + 1,
		}
	}
	return c
}
