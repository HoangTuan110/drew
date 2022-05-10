/*
	Main file of `drew`, a little small drawing app made with Tcell.
*/
package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	colorKeyMap = map[rune]tcell.Color{
		'1': tcell.ColorWhite,
		'2': tcell.NewRGBColor(26, 28, 44), // Light black color
		'3': tcell.ColorRed,
		'4': tcell.ColorGreen,
		'5': tcell.ColorYellow,
		'6': tcell.ColorBlue,
		'7': tcell.ColorDarkMagenta,
		'8': tcell.ColorBrown,
	}
)

// Draw texts.
func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for _, c := range text {
		s.SetCell(x, y, style, c)
		x++
	}
}

// Draw boxes.
func drawBox(s tcell.Screen, x, y, w, h int, style tcell.Style) {
	// Fill the inside.
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			s.SetCell(x+i, y+j, style, ' ')
		}
	}

	// Draw borders.
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			if lx == 0 || lx == w-1 || ly == 0 || ly == h-1 {
				s.SetCell(x+lx, y+ly, style, ' ')
			}
		}
	}
}

func main() {	
	// Initialize some stuff.
	x, y := 0, 0 // Current position of the pencil.
	w, h := 1, 1 // Width and height (aka size) of the pencil.
	eraseMode := false // Erase mode.
	primaryColor := colorKeyMap['1'] // Primary color (by default is Black).
	secondaryColor := colorKeyMap['1'] // Secondary color (by default is Black).
	currentColor := "primary" // Currently used color (by default is Black).

	// Default style.
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)
	pencilStyle := tcell.StyleDefault.Background(primaryColor).Foreground(tcell.ColorBlack)

	// Initialize the screen.
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	if e = s.Init(); e != nil {
		panic(e)
	}
	
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// Quit function.
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	// This will ensure that if there is an error, it will not affect the terminal.
	defer s.Fini()

	// Event loop.
	for {
    // Show the screen.
		s.Show()

		// Poll event.
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if _, ok := colorKeyMap[ev.Rune()]; ok {
				// Check the color.
				if currentColor == "primary" {
					primaryColor = colorKeyMap[ev.Rune()]
					pencilStyle = pencilStyle.Background(primaryColor)
				} else {
					secondaryColor = colorKeyMap[ev.Rune()]
					pencilStyle = pencilStyle.Background(secondaryColor)
				}
			}
			if ev.Rune() == 'x' || ev.Rune() == 'X' {
				if currentColor == "primary" {
					currentColor = "secondary"
					pencilStyle = pencilStyle.Background(secondaryColor)
				} else {
					currentColor = "primary"
					pencilStyle = pencilStyle.Background(primaryColor)
				}
			}
			if ev.Rune() == 'q' || ev.Rune() == 'Q' {
				quit()
			}
			if ev.Rune() == 'c' || ev.Rune() == 'C' {
				s.Clear()
			}
			if ev.Rune() == 'e' {
				eraseMode = !eraseMode
			}
			if ev.Rune() == ']' {
				// Increase the size of the pencil.
				w++
				h++
			}
			if ev.Rune() == '[' {
				// Decrease the size of the pencil.
				if w > 1 && h > 1 {
					w--
					h--
				}
			}
		case *tcell.EventMouse:
			x, y = ev.Position()
			button := ev.Buttons()

			// Only mouse event, not wheel event.
			button &= tcell.ButtonMask(0xff)

			if button == tcell.Button1 {
				if eraseMode {
					drawBox(s, x, y, w, h, defStyle) // Reset to default style.
				} else {
					drawBox(s, x, y, w, h, pencilStyle)
				}
			}
		}

		// Draw help.
		drawText(s, 0, 0, defStyle, "q - quit | x - switch color | c - clear | e - erase | Left Click - draw | 1-9 - colors | ? - help (TODO)")
		// Draw some information.
		// If you are wondering why there are extra spaces, it's because
		// Tcell (and its ancestor, termbox-go) uses a list of cells to set stuff.
		// Because of this, there will be some unused cells, which can be confusing.
		// The easiest way to fix this is to add some extra spaces.
		drawText(s, 0, 1, defStyle, fmt.Sprintf("%d, %d | %d, %d   ", x, y, w, h))
		if currentColor == "primary" {
			drawText(s, 0, 2, pencilStyle, " ")
		} else {
			drawText(s, 0, 2, pencilStyle.Background(primaryColor), " ")
		}
		drawText(s, 1, 2, pencilStyle.Background(secondaryColor) , " ")
		drawText(s, 0, 3, defStyle, fmt.Sprintf("%s | %t   ", currentColor, eraseMode))
	}
}
