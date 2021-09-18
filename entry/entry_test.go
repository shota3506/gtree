package entry

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T) {
	root, err := NewRoot("testdata", false)
	require.NoError(t, err)

	require.Len(t, root.children, 2)
	c := root.children[1]

	f, ok := c.(*File)
	require.True(t, ok)

	assert.Equal(t, "sample1.txt", f.Name())
	assert.Equal(t, filepath.Join(root.Path(), f.Name()), f.Path())
	assert.False(t, f.IsDir())
	assert.Equal(t, 1, f.Size())
	assert.Equal(t, 1, f.Depth())
}

func TestDir(t *testing.T) {
	root, err := NewRoot("testdata", false)
	require.NoError(t, err)

	require.Len(t, root.children, 2)
	c := root.children[0]

	d, ok := c.(*Dir)
	require.True(t, ok)

	assert.Equal(t, "a", d.Name())
	assert.Equal(t, filepath.Join(root.Path(), d.Name()), d.Path())
	assert.True(t, d.IsDir())
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 1, d.Depth())
}
