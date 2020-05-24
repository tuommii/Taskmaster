package tty

import (
	"fmt"
	"os"
)

// State represents terminal state
type State struct {
	// Cursor x position
	pos int
	// Width of user input
	inputLen int
	// Key presses
	key int
	// $>
	prompt    string
	promptLen int
	// Width
	cols           int
	buf            []byte
	history        []string
	historyCount   int
	historyPos     int
	proposer       Proposer
	proposerPos    int
	proposerStatus bool
	suggestions    []string
	// Multiline support:
	// Rows      int
	// LinesUsed int
	// OldPos    int
}

// New returns new State
func New(maxLen int) *State {
	s := &State{
		cols:      80,
		prompt:    "$>",
		promptLen: 2,
		buf:       make([]byte, maxLen),
	}
	return s
}

// ReadKey reads one byte at time
func (s *State) ReadKey(ch chan os.Signal) string {
	var code int
	s.clearBuffer()
	s.printPrompt()
	for {
		code = keyPressed()
		s.key = code
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
		case code == Tab:
			s.handleTab()
		}
	}
}

func (s *State) handlePrintable() {
	if s.pos == s.inputLen {
		s.buf = append(s.buf, byte(s.key))
		s.pos++
		s.inputLen++
		s.clearLine()
		s.printPrompt()
		s.printBuffer()
		// s.pos++
		// s.inputLen++
	} else {
		// make space for a new char
		s.buf = append(s.buf, '0')
		// shift
		copy(s.buf[s.pos+1:], s.buf[s.pos:])
		s.buf[s.pos] = byte(s.key)

		s.clearLine()
		s.printPrompt()
		s.printBuffer()

		s.pos++
		s.inputLen++

		// Move cursor to right place
		fmt.Print("\r")
		fmt.Printf("\033[%dC", s.pos+s.promptLen)
	}
	s.proposerPos = 0
}

func (s *State) handleEnter() string {
	input := string(s.buf)
	s.historyAdd(input)
	s.clearBuffer()
	fmt.Print("\n\r")
	s.proposerPos = 0
	// s.printPrompt()
	return input
}

func remove(slice []byte, s int) []byte {
	return append(slice[:s], slice[s+1:]...)
}

func (s *State) handleBackspace() {
	if s.pos == 0 {
		return
	}
	s.buf = remove(s.buf, s.pos-1)
	s.pos--
	s.inputLen--
	fmt.Print("\r\033[K")
	fmt.Print(s.prompt)
	fmt.Print(string(s.buf))
	fmt.Print("\r")
	fmt.Printf("\033[%dC", s.pos+s.promptLen)
	s.proposerPos = 0
}

func (s *State) handleLeft() {
	if s.pos > 0 {
		s.pos--
		fmt.Print("\033[1D")
	}
}

func (s *State) handleRight() {
	if s.pos < s.inputLen {
		s.pos++
		fmt.Print("\033[1C")
	}
}

func (s *State) handleUp() {
	if s.historyCount == 0 {
		return
	}
	if s.historyPos < 0 {
		s.clearLine()
		s.clearBuffer()
		fmt.Print(s.prompt)
		s.inputLen = 0
		s.pos = 0
		s.historyPos = s.historyCount - 1
		return
	}
	s.buf = []byte(s.history[s.historyPos])
	s.pos = len(s.history[s.historyPos])
	s.inputLen = s.pos
	s.historyPos--
	fmt.Print("\r\033[K")
	fmt.Print(s.prompt)
	fmt.Print(string(s.buf))

}

func (s *State) handleDown() {
	// if s.historyCount == 0 {
	// 	return
	// }
	// if s.historyPos < 0 {
	// 	s.historyPos = 0
	// }
	// if s.historyPos == s.historyCount-1 {
	// 	s.clearLine()
	// 	s.clearBuffer()
	// 	fmt.Print(s.Prompt)
	// 	s.inputLen = 0
	// 	s.pos = 0
	// 	// s.historyPos = s.historyCount - 1
	// 	return
	// }
	// s.historyPos++
	// s.buf = []byte(s.History[s.historyPos])
	// s.pos = len(s.History[s.historyPos])
	// s.inputLen = s.pos
	// fmt.Print("\r\033[K")
	// fmt.Print(s.Prompt)
	// fmt.Print(string(s.buf))
}

func (s *State) handleTab() {
	if s.proposer == nil {
		return
	}
	if s.proposerPos == 0 {
		s.suggestions = s.proposer(string(s.buf))
	}
	if len(s.suggestions) == 0 {
		return
	}
	if s.proposerPos >= len(s.suggestions) {
		s.proposerPos = 0
	}
	if len(s.suggestions[s.proposerPos]) > 0 {
		s.clearLine()
		s.clearBuffer()
		s.buf = []byte(s.suggestions[s.proposerPos])
		fmt.Print(s.prompt)
		s.printBuffer()
		s.inputLen = len(s.buf)
		s.pos = len(s.buf)
		s.proposerPos++
	}
}

func (s *State) clearBuffer() {
	s.buf = s.buf[:0]
	s.inputLen = 0
	s.pos = 0
}

// clearLine clears current line
func (s *State) clearLine() {
	fmt.Print("\r\033[K")
}

// printPrompt prints prompt
func (s *State) printPrompt() {
	fmt.Print(s.prompt)
}

// printBuffer prints buffer
func (s *State) printBuffer() {
	fmt.Print(string(s.buf))
}
