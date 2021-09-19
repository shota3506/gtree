package state

import (
	"testing"

	"github.com/shota3506/gtree/entry"
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

func TestStateUpAndDown(t *testing.T) {
	root, err := entry.NewRoot("testdata", false)
	require.NoError(t, err)

	st := NewState(root, 100, 100)
	require.Equal(t, 4, st.Root().Size())
	require.Equal(t, 0, st.offset)
	require.Equal(t, 0, st.pos)

	st.Down()
	assert.Equal(t, 1, st.pos)
	st.Down()
	assert.Equal(t, 2, st.pos)
	st.Down()
	assert.Equal(t, 3, st.pos)
	st.Down()
	assert.Equal(t, 3, st.pos)

	st.Up()
	assert.Equal(t, 2, st.pos)
	st.Up()
	assert.Equal(t, 1, st.pos)
	st.Up()
	assert.Equal(t, 0, st.pos)
	st.Up()
	assert.Equal(t, 0, st.pos)
}

func TestStateUpAndDownWithOffsetAdjustment(t *testing.T) {
	root, err := entry.NewRoot("testdata", false)
	require.NoError(t, err)

	st := NewState(root, 100, 2)
	require.Equal(t, 4, st.Root().Size())
	require.Equal(t, 0, st.offset)
	require.Equal(t, 0, st.pos)

	st.Down()
	assert.Equal(t, 1, st.pos)
	assert.Equal(t, 0, st.offset)
	st.Down()
	assert.Equal(t, 2, st.pos)
	assert.Equal(t, 1, st.offset)
	st.Down()
	assert.Equal(t, 3, st.pos)
	assert.Equal(t, 2, st.offset)

	st.Up()
	assert.Equal(t, 2, st.pos)
	assert.Equal(t, 2, st.offset)
	st.Up()
	assert.Equal(t, 1, st.pos)
	assert.Equal(t, 1, st.offset)
	st.Up()
	assert.Equal(t, 0, st.pos)
	assert.Equal(t, 0, st.offset)
}

func TestStateSetSize(t *testing.T) {
	root, err := entry.NewRoot("testdata", false)
	require.NoError(t, err)

	st := NewState(root, 100, 200)
	require.Equal(t, 100, st.width)
	require.Equal(t, 200, st.height)

	st.SetSize(10, 20)
	require.Equal(t, 10, st.width)
	require.Equal(t, 20, st.height)
}

func TestStateSetSizeWithOffsetAdjustment(t *testing.T) {
	root, err := entry.NewRoot("testdata", false)
	require.NoError(t, err)

	st := NewState(root, 100, 200)
	st.Down()
	st.Down()
	st.Down()

	assert.Equal(t, 0, st.offset)
	st.SetSize(100, 2)
	assert.Equal(t, 2, st.offset)
}

func TestStateToggle(t *testing.T) {
	root, err := entry.NewRoot("testdata", false)
	require.NoError(t, err)

	st := NewState(root, 100, 100)
	require.Equal(t, 4, st.Root().Size())
	st.Down()
	require.Equal(t, 1, st.pos)

	e, err := st.Root().Get(st.pos)
	require.NoError(t, err)
	d, ok := e.(*entry.Dir)
	require.True(t, ok)

	err = st.Toggle() // open
	require.NoError(t, err)
	assert.True(t, d.IsOpen())

	err = st.Toggle() // close
	require.NoError(t, err)
	assert.False(t, d.IsOpen())
}
