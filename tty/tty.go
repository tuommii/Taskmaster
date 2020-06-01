package tty

import (
	"fmt"
	"os"
	"strings"
)

// State represents terminal state
type State struct {
	// Cursor x position
	pos int
	// Length of user input
	inputLen int
	// Key pressed
	key int
	// Prompt string e.g "$>""
	prompt    string
	promptLen int
	// Width
	cols int
	buf  []byte
	hist
	autocomplete
	// Multiline support:
	// Rows      int
	// LinesUsed int
	// OldPos    int
}

// New returns new State
func New(maxLen int) *State {
	s := &State{
		cols:      80,
		prompt:    "taskmaster$>",
		promptLen: 12,
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
		case code == CtrlC:
			ch <- os.Interrupt
		case code == CtrlD:
			ch <- os.Interrupt
		case code == Esc:
			return "exit"
		case isPrintable(code):
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
		// append
		s.buf = append(s.buf, byte(s.key))
		s.pos++
		s.inputLen++
		s.clearLine()
		s.printPrompt()
		s.printBuffer()
	} else {
		// insert
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
		// fmt.Print("\r")
		// fmt.Printf("\033[%dC", s.pos+s.promptLen)
		s.restoreCursor()
	}
	// reset autocomplete suggestions
	s.resetSuggestions()
}

func (s *State) handleEnter() string {
	input := string(s.buf)
	s.historyAdd(input)
	s.clearBuffer()
	fmt.Print("\n\r")
	s.resetSuggestions()
	return input
}

func (s *State) handleBackspace() {
	if s.pos == 0 {
		return
	}
	s.buf = remove(s.buf, s.pos-1)
	s.pos--
	s.inputLen--
	s.clearLine()
	s.printPrompt()
	s.printBuffer()
	// fmt.Print("\r")
	// fmt.Printf("\033[%dC", s.pos+s.promptLen)
	s.restoreCursor()
	s.resetSuggestions()
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
		s.printPrompt()
		s.inputLen = 0
		s.pos = 0
		s.historyPos = s.historyCount - 1
		return
	}
	s.buf = []byte(s.history[s.historyPos])
	s.pos = len(s.history[s.historyPos])
	s.inputLen = s.pos
	s.historyPos--
	s.clearLine()
	s.printPrompt()
	s.printBuffer()
}

func (s *State) handleDown() {
	// TODO: Figure this out
}

func (s *State) handleTab() {
	if s.proposer == nil {
		return
	}
	if s.proposerPos == 0 {
		s.suggestions = s.proposer(string(s.buf), s.jobNames)
	}
	if len(s.suggestions) == 0 {
		return
	}
	if s.proposerPos >= len(s.suggestions) {
		s.resetSuggestions()
	}
	if len(s.suggestions[s.proposerPos]) > 0 {
		s.clearLine()
		var name bool
		if strings.Contains(string(s.buf), " ") {
			name = true
		}
		cpy := s.buf
		splitted := strings.SplitN(string(cpy), " ", 2)
		s.clearBuffer()
		if name {
			s.buf = []byte(splitted[0] + " " + s.suggestions[s.proposerPos])
		} else {
			s.buf = []byte(s.suggestions[s.proposerPos])
		}
		width := len(s.buf)
		s.inputLen = width
		s.pos = width
		s.printPrompt()
		s.printBuffer()
		s.proposerPos++
	}
}

func (s *State) clearBuffer() {
	s.buf = s.buf[:0]
	s.inputLen = 0
	s.pos = 0
}

func (s *State) restoreCursor() {
	fmt.Printf("\r\033[%dC", s.pos+s.promptLen)
}

// clearLine clears current line
func (s *State) clearLine() {
	fmt.Print("\r\033[K")
}

// printPrompt prints prompt
func (s *State) printPrompt() {
	fmt.Printf("\033[1;34m%s\033[0m", s.prompt)
}

// printBuffer prints buffer
func (s *State) printBuffer() {
	fmt.Print(string(s.buf))
}

func (s *State) resetSuggestions() {
	s.proposerPos = 0
}

func remove(slice []byte, s int) []byte {
	return append(slice[:s], slice[s+1:]...)
}
