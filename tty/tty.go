package tty

import (
	"bytes"
	"fmt"
	"os"
)

// State represents terminal state
type State struct {
	// Cursor x position
	Pos       int
	Key       int
	Input     string
	Prompt    string
	Cols      int
	PromptLen int
	InputLen  int
	Buffer    *bytes.Buffer
	buf       []byte
	// If multiline support:
	// Rows      int
	// LinesUsed int
	// OldPos    int
}

// New returns new State
func New(maxLen int) *State {
	s := &State{
		Cols:      80,
		Prompt:    "$>",
		PromptLen: 2,
		buf:       make([]byte, maxLen),
	}
	return s
}

// ReadKey reads one byte at time
func (s *State) ReadKey(ch chan os.Signal) string {
	var code int
	s.ClearBuffer()
	s.PrintPrompt()
	for {
		code = keyPressed()
		s.Key = code
		switch {
		case code == 3:
			ch <- os.Interrupt
		case code == 4:
			ch <- os.Interrupt
			// return "exit"
		case code == Esc:
			return "exit"
		case IsPrintable(code):
			s.handlePrintable()
		case code == Enter:
			return s.handleEnter()
		case code == Backspace:
			s.handleBackspace()
		case code == Left:
			s.handleLeft()
		case code == Right:
			s.handleRight()
		case code == Up:
			s.handleUp()
		case code == Down:
			s.handleDown()
		}
	}
}

func (s *State) handlePrintable() {
	if s.Pos == s.InputLen {
		s.buf = append(s.buf, byte(s.Key))
		s.Pos++
		s.InputLen++
		s.ClearLine()
		s.PrintPrompt()
		s.PrintBuffer()
		// s.Pos++
		// s.InputLen++
	} else {
		// make space for a new char
		s.buf = append(s.buf, '0')
		// shift
		copy(s.buf[s.Pos+1:], s.buf[s.Pos:])
		s.buf[s.Pos] = byte(s.Key)

		s.ClearLine()
		s.PrintPrompt()
		s.PrintBuffer()

		s.Pos++
		s.InputLen++

		// Move cursor
		fmt.Print("\r")
		for i := 0; i < s.Pos+s.PromptLen; i++ {
			fmt.Print("\033[1C")
		}
	}
}

func (s *State) handleEnter() string {
	input := string(s.buf)
	s.ClearBuffer()
	fmt.Print("\n\r")
	// s.PrintPrompt()
	return input
}

func (s *State) handleBackspace() {

}

func (s *State) handleLeft() {
	if s.Pos > 0 {
		s.Pos--
		fmt.Print("\033[1D")
	}
}

func (s *State) handleRight() {
	if s.Pos < s.InputLen {
		s.Pos++
		fmt.Print("\033[1C")
	}
}

func (s *State) handleUp() {
	s.buf = s.buf[:0]
	s.buf = []byte("miikka")
	s.Pos = 6
	s.InputLen = 6
	fmt.Print("\r\033[K")
	fmt.Print(s.Prompt)
	fmt.Print(string(s.buf))

}

func (s *State) handleDown() {

}

// ClearBuffer clears buffer
func (s *State) ClearBuffer() {
	s.buf = s.buf[:0]
	s.InputLen = 0
	s.Pos = 0
}

// ClearLine clears current line
func (s *State) ClearLine() {
	fmt.Print("\r\033[K")
}

// PrintPrompt prints prompt
func (s *State) PrintPrompt() {
	fmt.Print(s.Prompt)
}

// PrintBuffer prints buffer
func (s *State) PrintBuffer() {
	fmt.Print(string(s.buf))
}
