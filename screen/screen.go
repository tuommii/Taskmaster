// Package screen provides cursor and screen
// manipulation operations for VT-100 terminals
package screen

import "fmt"

// Terminal represents terminal window
type Terminal struct {
	// Cursor x position
	X      int
	Width  int
	Height int
}

// New creates new screen instance
func New() *Terminal {
	term := &Terminal{}
	return term
}

// ResetAll custom styles
func (s *Terminal) ResetAll() {
	fmt.Printf("\033[0m")
}

// ResetColor to default color
func (s *Terminal) ResetColor() {
	fmt.Printf("\033[32m")
}

// ResetLine clears line and put cursor at beginning of it
func (s *Terminal) ResetLine() {
	fmt.Printf("\r\033[K")
}

// MoveCursor to given position
func (s *Terminal) MoveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// MoveCursorRight steps
func (s *Terminal) MoveCursorRight(step int) {
	fmt.Printf("\033[%dC", step)
}

// MoveCursorLeft ...
func (s *Terminal) MoveCursorLeft(step int) {
	fmt.Printf("\033[%dD", step)
}

// MoveCursorDown steps
func (s *Terminal) MoveCursorDown(step int) {
	fmt.Printf("\033[%dB", step)
}

// MoveCursorUp steps
func (s *Terminal) MoveCursorUp(step int) {
	fmt.Printf("\033[%dA", step)
}

// Clear Terminal
func (s *Terminal) Clear() {
	fmt.Printf("\033[2J")
}
