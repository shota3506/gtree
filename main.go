package main

import (
	"flag"
	"log"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/gdamore/tcell/v2/views"
	"github.com/shota3506/gtree/commands"
	"github.com/shota3506/gtree/entry"
	"github.com/shota3506/gtree/state"
)

var showHidden bool

func init() {
	const (
		showHiddenUsage = "show hidden files or directories"
	)
	flag.BoolVar(&showHidden, "hidden", false, showHiddenUsage)
	flag.BoolVar(&showHidden, "h", false, showHiddenUsage+" (shorthand)")
}

func main() {
	flag.Parse()

	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	path := "."
	if args := flag.Args(); len(args) == 1 {
		path = args[0]
	}

	root, err := entry.NewRoot(path, showHidden)
	if err != nil {
		return err
	}

	s, err := initScreen()
	if err != nil {
		return err
	}
	defer s.Fini()

	commandChan := make(chan commands.Command)
	stateChan := make(chan state.State)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		handleCommand(s, commandChan)
		wg.Done()
	}()
	go func() {
		handleChange(s, stateChan)
		wg.Done()
	}()
	go func() {
		start(s, root, commandChan, stateChan)
		wg.Done()
	}()
	wg.Wait()

	return nil
}

func initScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err = s.Init(); err != nil {
		return nil, err
	}
	defaultStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defaultStyle)
	s.DisableMouse()
	s.DisablePaste()
	s.Clear()
	return s, nil
}

func handleCommand(s tcell.Screen, commandChan chan commands.Command) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(commandChan)
				return
			case tcell.KeyUp:
				commandChan <- commands.CommandUp{}
			case tcell.KeyDown:
				commandChan <- commands.CommandDown{}
			case tcell.KeyEnter:
				commandChan <- commands.CommandSelect{}
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'q':
					close(commandChan)
					return
				case 'k':
					commandChan <- commands.CommandUp{}
				case 'j':
					commandChan <- commands.CommandDown{}
				case ' ':
					commandChan <- commands.CommandSelect{}
				}
			}
		case *tcell.EventResize:
			width, height := s.Size()
			commandChan <- commands.CommandResize{Width: width, Height: height - 1}
		}
	}
}

func handleChange(s tcell.Screen, stateChan chan state.State) {
	for {
		state, ok := <-stateChan
		if !ok {
			return
		}
		render(s, state)
	}
}

func start(s tcell.Screen, root *entry.Dir, commandChan chan commands.Command, stateChan chan state.State) {
	width, height := s.Size()
	state := state.NewState(root, width, height-1)
	stateChan <- state

	for {
		command, ok := <-commandChan
		if !ok {
			close(stateChan)
			return
		}
		if nextState, err := command.Do(state); err == nil {
			state = nextState
			stateChan <- state
		}
	}
}

func render(s tcell.Screen, st state.State) {
	s.Clear()

	layout := views.NewBoxLayout(views.Vertical)

	v := views.NewCellView()
	v.SetModel(st.View())

	statusBar := views.NewSimpleStyledTextBar()
	statusBar.SetLeft(st.Root().Path())

	layout.SetView(s)
	layout.AddWidget(v, 1.0)
	layout.AddWidget(statusBar, 0.0)
	layout.Draw()

	s.Show()
}
