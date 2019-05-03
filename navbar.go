package main

import (
	"github.com/gdamore/tcell"
)

const scrollAhead = 3

type NavBar struct {
	items  []tcell.Color // navigation colors
	label  string
	pos    int
	offset int
	width  int
	pst    tcell.Style // pointer style
	state  *State
}

func NewNavBar(s *State, length int) *NavBar {
	return &NavBar{
		items: make([]tcell.Color, length),
		state: s,
	}
}

// Draw redraws bar at given coordinates and screen, returning the number
// of rows occupied
func (bar *NavBar) Draw(x, y int, s tcell.Screen) int {
	var st tcell.Style

	n := bar.offset
	col := 0

	// border bars
	s.SetCell(x-1, y+1, bar.pst, '│')
	s.SetCell(x-1, y+2, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+1, bar.pst, '│')
	s.SetCell(bar.width+x+1, y+2, bar.pst, '│')

	for col <= bar.width && n < len(bar.items) {
		st = st.Background(bar.items[n])
		s.SetCell(col+x, y, blkSt, '█')
		s.SetCell(col+x, y+1, st, ' ')
		s.SetCell(col+x, y+2, st, ' ')

		col++
		n++
	}

	ix := (bar.pos - bar.offset) + x
	s.SetCell(ix, y, bar.pst, '▾')
	s.SetCell(bar.width/2, y+3, bar.pst, []rune(bar.label)...)

	return 4
}

func (bar *NavBar) SetLabel(s string) { bar.label = s }

func (bar *NavBar) SetPos(idx int) {
	switch {
	case idx > bar.pos:
		bar.up(idx - bar.pos)
	case idx < bar.pos:
		bar.down(bar.pos - idx)
	}
}

func (bar *NavBar) Resize(w int) {
	bar.width = w
	bar.up(0)
	bar.down(0)
}

// NavBar implements Section
func (bar *NavBar) Up(int)                         {}
func (bar *NavBar) Down(int)                       {}
func (bar *NavBar) Handle(StateChange)             {}
func (bar *NavBar) Width() int                     { return bar.width }
func (bar *NavBar) SetPointerStyle(st tcell.Style) { bar.pst = st }

func (bar *NavBar) up(step int) {
	max := len(bar.items) - 1
	maxOffset := max - bar.width

	switch {
	case step <= 0:
	case bar.pos == max:
		return
	case bar.pos+step > max:
		bar.pos = max
	default:
		bar.pos += step
	}

	if (bar.pos - bar.offset) > bar.width-scrollAhead {
		bar.offset = (bar.pos - bar.width) + scrollAhead
	}
	if bar.offset >= maxOffset {
		bar.offset = maxOffset
	}
}

func (bar *NavBar) down(step int) {
	switch {
	case step <= 0:
	case bar.pos == 0:
		return
	case bar.pos-step < 0:
		bar.pos = 0
	default:
		bar.pos -= step
	}

	if bar.pos-bar.offset < scrollAhead {
		bar.offset = bar.pos - scrollAhead
	}
	if bar.offset < 0 {
		bar.offset = 0
	}
}