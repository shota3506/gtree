package commands

import (
	"github.com/shota3506/gtree/state"
)

type Command interface {
	Do(state state.State) (state.State, error)
}

type CommandUp struct{}

func (c CommandUp) Do(st state.State) (state.State, error) {
	st.Up()
	return st, nil
}

type CommandDown struct{}

func (c CommandDown) Do(st state.State) (state.State, error) {
	st.Down()
	return st, nil
}

type CommandResize struct {
	Width  int
	Height int
}

func (c CommandResize) Do(st state.State) (state.State, error) {
	st.SetSize(c.Width, c.Height)
	return st, nil
}

type CommandSelect struct{}

func (c CommandSelect) Do(st state.State) (state.State, error) {
	err := st.Toggle()
	if err != nil {
		return nil, err
	}
	return st, nil
}
