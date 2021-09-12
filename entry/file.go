package entry

type File struct {
	name  string
	path  string
	depth int
}

func (f *File) Name() string {
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
