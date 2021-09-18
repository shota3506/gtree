package state

import (
	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shota3506/gtree/entry"
)

type State interface {
	View() views.CellModel
	Root() *entry.Dir
	Up()
	Down()
	SetSize(width, height int)
	Toggle() error
}

type state struct {
	root   *entry.Dir
	pos    int
	offset int
	width  int
	height int
}

func NewState(root *entry.Dir, width, height int) State {
	return &state{
		root:   root,
		pos:    0,
		offset: 0,
		width:  width,
		height: height,
	}
}

func (s *state) Root() *entry.Dir {
	return s.root
}

func (st *state) Up() {
	if st.pos > 0 {
		st.pos -= 1
	}

	st.adjust()
}

func (st *state) Down() {
	if st.pos+1 < st.Root().Size() {
		st.pos += 1
	}

	st.adjust()
}

func (st *state) SetSize(width, height int) {
	st.width = width
	st.height = height

	st.adjust()
}

func (st *state) Toggle() error {
	e, err := st.Root().Get(st.pos)
	if err != nil {
		return err
	}

	if d, ok := e.(*entry.Dir); ok {
		if d.IsOpen() {
			d.Close()
		} else {
			if err = d.Open(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (st *state) adjust() {
	if st.pos < st.offset {
		st.offset = st.pos
	}
	if st.pos-st.offset >= st.height {
		st.offset = st.pos - st.height + 1
	}
}

func (st *state) View() views.CellModel {
	items := []*struct {
		Name  string
		Depth int
	}{}
	_ = st.Root().Walk(func(e entry.Entry) error {
		name := e.String()
		if e, ok := e.(*entry.Dir); ok {
			if e.IsOpen() {
				name = "▾ " + name
			} else {
				name = "▸ " + name
			}
		}

		items = append(items, &struct {
			Name  string
			Depth int
		}{
			Name:  name,
			Depth: e.Depth(),
		})
		return nil
	})

	m := map[int]interface{}{}
	lines := []*ViewLine{}
	for i := len(items) - 1; i >= 0; i-- {
		prefix := ""
		for j := 0; j < items[i].Depth-1; j++ {
			if _, ok := m[j+1]; ok {
				prefix += " │ "
			} else {
				prefix += "   "
			}
		}
		if items[i].Depth > 0 {
			if _, ok := m[items[i].Depth]; ok {
				prefix += " ├─ "
			} else {
				prefix += " └─ "
			}
		}

		style := tcell.StyleDefault
		if i == st.pos {
			style = style.Reverse(true)
		}

		lines = append([]*ViewLine{{prefix + items[i].Name, style}}, lines...)
		m[items[i].Depth] = struct{}{}
		if i < len(items)-1 {
			for depth := items[i].Depth + 1; depth <= items[i+1].Depth; depth++ {
				delete(m, depth)
			}
		}
	}

	return &viewState{
		lines:  lines[st.offset:],
		width:  st.width,
		height: st.height,
	}
}

type viewState struct {
	lines  []*ViewLine
	width  int
	height int
}

func (v *viewState) GetCell(x, y int) (rune, tcell.Style, []rune, int) {
	if y >= len(v.lines) || y >= v.height {
		return ' ', tcell.StyleDefault, nil, 1
	}

	style := v.lines[y].style
	rs := []rune(v.lines[y].String())
	if x < len(rs) {
		return rs[x], style, nil, 1
	}
	return ' ', style, nil, 1
}

func (v *viewState) GetBounds() (int, int) {
	return v.width, v.height
}

func (v *viewState) GetCursor() (int, int, bool, bool) {
	return 0, 0, false, false
}

func (v *viewState) SetCursor(int, int) {}

func (v *viewState) MoveCursor(offx, offy int) {}

type ViewLine struct {
	text  string
	style tcell.Style
}

func (l *ViewLine) String() string {
	return l.text
}

func (l *ViewLine) Style() tcell.Style {
	return l.style
}
