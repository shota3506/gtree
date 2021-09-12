package entry

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Dir struct {
	name       string
	path       string
	children   []Entry
	depth      int
	open       bool
	showHidden bool
}

func (d *Dir) Name() string {
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
	// never return error
	_ = d.Walk(func(e Entry) error {
		size += 1
		return nil
	})
	return size
}

func (d *Dir) Depth() int {
	return d.depth
}

func (d *Dir) Children() []Entry {
	if !d.open {
		return nil
	}
	if d.children == nil {
		if err := d.Read(); err != nil {
			return nil
		}
	}
	return d.children
}

func (d *Dir) Open() error {
	if err := d.Read(); err != nil {
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

func (d *Dir) Flatten() []Entry {
	es := []Entry{}
	// never return error
	_ = d.Walk(func(e Entry) error {
		es = append(es, e)
		return nil
	})
	return es
}

func (d *Dir) Walk(fn func(e Entry) error) error {
	if err := fn(d); err != nil {
		return err
	}

	for _, c := range d.Children() {
		switch c := c.(type) {
		case *Dir:
			if err := c.Walk(fn); err != nil {
				return err
			}
		default:
			if err := fn(c); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Dir) Get(i int) (Entry, error) {
	if i < 0 {
		return nil, errors.New("not found")
	}

	entries := d.Flatten()
	if i >= len(entries) {
		return nil, errors.New("not found")
	}
	return entries[i], nil
}

func (d *Dir) Read() error {
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
		return children[i].Name() < children[j].Name()
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
