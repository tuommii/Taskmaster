package tty

import (
	"fmt"
	"os"
)

// State represents terminal state
type State struct {
	// Cursor x position
	Pos int
	// Width of user input
	InputLen int
	// Key presses
	Key int
	// $>
	Prompt    string
	PromptLen int
	// Width
	Cols         int
	buf          []byte
	History      []string
	HistoryCount int
	HistoryPos   int
	// Multiline support:
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
		// CTRL + C
		case code == 3:
			ch <- os.Interrupt
		// CTRL + D
		case code == 4:
			ch <- os.Interrupt
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

		// Move cursor to right place
		fmt.Print("\r")
		fmt.Printf("\033[%dC", s.Pos+s.PromptLen)
	}
}

func (s *State) handleEnter() string {
	input := string(s.buf)
	s.HistoryAdd(input)
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
	if s.HistoryCount == 0 {
		return
	}
	if s.HistoryPos < 0 {
		s.ClearLine()
		s.ClearBuffer()
		fmt.Print(s.Prompt)
		s.InputLen = 0
		s.Pos = 0
		s.HistoryPos = s.HistoryCount - 1
		return
	}
	s.buf = []byte(s.History[s.HistoryPos])
	s.Pos = len(s.History[s.HistoryPos])
	s.InputLen = s.Pos
	s.HistoryPos--
	fmt.Print("\r\033[K")
	fmt.Print(s.Prompt)
	fmt.Print(string(s.buf))

}

func (s *State) handleDown() {
	// if s.HistoryCount == 0 {
	// 	return
	// }
	// if s.HistoryPos < 0 {
	// 	s.HistoryPos = 0
	// }
	// if s.HistoryPos == s.HistoryCount-1 {
	// 	s.ClearLine()
	// 	s.ClearBuffer()
	// 	fmt.Print(s.Prompt)
	// 	s.InputLen = 0
	// 	s.Pos = 0
	// 	// s.HistoryPos = s.HistoryCount - 1
	// 	return
	// }
	// s.HistoryPos++
	// s.buf = []byte(s.History[s.HistoryPos])
	// s.Pos = len(s.History[s.HistoryPos])
	// s.InputLen = s.Pos
	// fmt.Print("\r\033[K")
	// fmt.Print(s.Prompt)
	// fmt.Print(string(s.buf))
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
