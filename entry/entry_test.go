package entry

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
testdata
├── a
│   ├── b
│   │   ├── sample4.txt
│   │   └── sample5.txt
│   ├── c
│   │   └── sample6.txt
│   └── sample3.txt
├── sample1.txt
└── sample2.txt
*/

func TestFile(t *testing.T) {
	root, err := NewRoot("testdata", false)
	require.NoError(t, err)

	require.Len(t, root.children, 3)
	f, ok := root.children[1].(*File)
	require.True(t, ok)

	assert.Equal(t, "sample1.txt", f.String())
	assert.Equal(t, filepath.Join(root.Path(), f.String()), f.Path())
	assert.False(t, f.IsDir())
	assert.Equal(t, 1, f.Size())
	assert.Equal(t, 1, f.Depth())
}

func TestDir(t *testing.T) {
	root, err := NewRoot("testdata", false)
	require.NoError(t, err)

	require.Len(t, root.children, 3)
	d, ok := root.children[0].(*Dir)
	require.True(t, ok)

	assert.Equal(t, "a", d.String())
	assert.Equal(t, filepath.Join(root.Path(), d.String()), d.Path())
	assert.True(t, d.IsDir())
	assert.Equal(t, 1, d.Size())
	assert.Equal(t, 1, d.Depth())
}

func TestDirOpenAndClose(t *testing.T) {
	root, err := NewRoot("testdata", false)
	require.NoError(t, err)

	require.Len(t, root.children, 3)
	da, ok := root.children[0].(*Dir)
	require.True(t, ok)
	assert.Equal(t, 1, da.Depth())
	assert.False(t, da.IsOpen())

	// open ./testdata/a
	err = da.Open()
	require.NoError(t, err)
	assert.True(t, da.IsOpen())
	assert.Equal(t, 4, da.Size())

	require.Len(t, da.children, 3)
	db, ok := da.children[0].(*Dir)
	require.True(t, ok)
	assert.Equal(t, 2, db.Depth())
	assert.False(t, db.IsOpen())

	// open ./testdata/a/b
	err = db.Open()
	require.NoError(t, err)
	assert.True(t, db.IsOpen())
	assert.Equal(t, 3, db.Size())
	assert.Equal(t, 6, da.Size())

	require.Len(t, db.children, 2)
	assert.Equal(t, 3, db.children[0].Depth())
	assert.Equal(t, 3, db.children[1].Depth())

	// close ./testdata/a/b
	db.Close()
	assert.False(t, db.IsOpen())
	assert.Equal(t, 1, db.Size())
	assert.Equal(t, 4, da.Size())

	// close ./testdata/a
	da.Close()
	assert.False(t, da.IsOpen())
	assert.Equal(t, 1, da.Size())
}

func TestDirGet(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		root, err := NewRoot("testdata", false)
		require.NoError(t, err)

		e, err := root.Get(3)
		require.NoError(t, err)
		require.NotNil(t, e)
		assert.Equal(t, "sample2.txt", e.String())
	})

	t.Run("not found", func(t *testing.T) {
		root, err := NewRoot("testdata", false)
		require.NoError(t, err)

		_, err = root.Get(4)
		require.Error(t, err)
	})

	t.Run("open", func(t *testing.T) {
		root, err := NewRoot("testdata", false)
		require.NoError(t, err)

		e, err := root.Get(1)
		require.NoError(t, err)
		d, ok := e.(*Dir)
		require.True(t, ok)

		// open ./testdata/a
		err = d.Open()
		require.NoError(t, err)

		e, err = root.Get(3)
		require.NoError(t, err)
		require.NotNil(t, e)
		assert.Equal(t, "c", e.String())

		e, err = root.Get(4)
		require.NoError(t, err)
		require.NotNil(t, e)
		assert.Equal(t, "sample3.txt", e.String())
	})
}
