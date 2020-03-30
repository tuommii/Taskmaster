package tty

import "fmt"

// MoveCursor to given position
func (t *Terminal) MoveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// MoveCursorRight steps
func (t *Terminal) MoveCursorRight(step int) {
	fmt.Printf("\033[%dC", step)
}

// MoveCursorLeft ...
func (t *Terminal) MoveCursorLeft(step int) {
	fmt.Printf("\033[%dD", step)
}

// MoveCursorDown steps
func (t *Terminal) MoveCursorDown(step int) {
	fmt.Printf("\033[%dB", step)
}

// MoveCursorUp steps
func (t *Terminal) MoveCursorUp(step int) {
	fmt.Printf("\033[%dA", step)
}
