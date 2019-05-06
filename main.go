package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bcicen/tcolors/logging"
	"github.com/gdamore/tcell"
	"github.com/teacat/noire"
)

var (
	log   = logging.Init()
	boxW  = 120
	boxH  = 80
	disp  *Display
	blkSt = tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlack)
	indicatorSt = tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(110, 110, 110)).
			Background(tcell.ColorBlack)
	hiIndicatorSt = tcell.StyleDefault.
			Foreground(tcell.NewRGBColor(255, 255, 255)).
			Background(tcell.ColorBlack)
	stepBasis int
)

const (
	littleStep = 1
	bigStep    = 10
)

func draw(s tcell.Screen) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	ly := 1
	st := tcell.StyleDefault

	if stepBasis == bigStep {
		s.SetCell(1, 0, st, '⏩')
	} else {
		s.SetCell(1, 0, st, ' ')
	}

	x := paddingX
	if disp.width == maxWidth {
		x = (w - maxWidth) / 2
	}
	ly += disp.Draw(x, ly, s)

	s.Show()
}

func main() {
	defer log.Exit()
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()
	w, _ := s.Size()
	disp = NewDisplay()
	disp.Resize(w)
	disp.build()

	quit := make(chan struct{})
	//go func() {
	//time.Sleep(1 * time.Second)
	//disp.SetColor(tcell.NewRGBColor(207, 064, 138))
	//draw(s)
	//}()
	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				stepBasis = littleStep
				if ev.Modifiers()&tcell.ModShift == tcell.ModShift {
					stepBasis = bigStep
				}
				switch ev.Key() {
				case tcell.KeyRune:
					switch ev.Rune() {
					case 'r':
						disp.Reset()
						draw(s)
					case 'h':
						if ok := disp.ValueDown(stepBasis); ok {
							draw(s)
						}
					case 'k':
						if ok := disp.SectionUp(); ok {
							draw(s)
						}
					case 'j':
						if ok := disp.SectionDown(); ok {
							draw(s)
						}
					case 'l':
						if ok := disp.ValueUp(stepBasis); ok {
							draw(s)
						}
					case 'q':
						close(quit)
						return
					default:
						log.Debugf("ignoring key [%s]", string(ev.Rune()))
					}
				case tcell.KeyRight:
					if ok := disp.ValueUp(stepBasis); ok {
						draw(s)
					}
				case tcell.KeyLeft:
					if ok := disp.ValueDown(stepBasis); ok {
						draw(s)
					}
				case tcell.KeyUp:
					if ok := disp.SectionUp(); ok {
						draw(s)
					}
				case tcell.KeyDown:
					if ok := disp.SectionDown(); ok {
						draw(s)
					}
				case tcell.KeyEscape, tcell.KeyCtrlC:
					close(quit)
					return
				case tcell.KeyCtrlL:
					s.Sync()
				}
			case *tcell.EventResize:
				w, h := s.Size()
				log.Debugf("handling resize: w=%04d h=%04d", w, h)
				disp.Resize(w)
				s.Clear()
				draw(s)
				s.Sync()
			}
		}
	}()

	draw(s)

loop:
	for {
		select {
		case <-quit:
			break loop
		case <-time.After(time.Millisecond * 50):
		}
	}

	w, h := s.Size()

	s.Clear()
	//lx := 1
	//ly := 1
	//for i := 0.0; i < 360.5; i += 0.5 {
	//c := noire.NewHSL(i, 100, 50)
	//st := tcell.StyleDefault.Background(toTColor(c))
	//s.SetCell(lx, ly, st, []rune(fmt.Sprintf("%0.2f ", i))...)
	//lx += 2
	//if lx >= w {
	//ly++
	//lx = 0
	//}
	//}
	//s.Show()
	//s.Sync()
	//time.Sleep(5 * time.Second)

	s.Fini()
	fmt.Printf("w=%d h=%d hues=%d \n", w, h, len(disp.HueNav.items))
	for n, x := range disp.xHues {
		h, s, l := x.HSL()
		r, g, b := x.RGB()
		fmt.Printf("[%d] %+0.2f %+0.2f %+0.2f [%0.3f %0.3f %0.3f]\n", n, h, s, l, r, g, b)
	}
	for i := 0; i < 1; i++ {
		x := noire.NewRGB(207, 64, 138)
		h, s, l := x.HSL()
		r, g, b := x.RGB()
		fmt.Printf("[%d] %+0.2f %+0.2f %+0.2f [%0.3f %0.3f %0.3f]\n", 0, h, s, l, r, g, b)
	}
	disp.state.Save()
}
