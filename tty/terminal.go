// Package tty provides cursor and screen
// manipulation operations for VT-100 terminals
package tty

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
)

type state struct {
	history      []string
	historyMutex sync.RWMutex
	cols         int
	r            *bufio.Reader
	cursorRows   int
}

// Terminal represents terminal window
type Terminal struct {
	state
	// Cursor x position
	Pos       int
	KeyCode   int
	Input     string
	PromptStr string
	OldPos    int
	Cols      int
	Rows      int
	PromptLen int
	InputLen  int
	LinesUsed int
	Buffer    *bytes.Buffer
}

// Prompt sets prompt and prints it
func (t *Terminal) Prompt(prompt string) {
	t.PromptStr = prompt
	t.PromptLen = len(prompt)
	t.PrintPrompt()
}

// PrintPrompt prints prompt
func (t *Terminal) PrintPrompt() {
	t.Buffer.WriteString(t.PromptStr)
	fmt.Print(t.Buffer.String())
}

func (t *Terminal) Redraw() {
	rows := (t.PromptLen + t.InputLen + t.Cols - 1) / t.Cols
	t.ClearRows()
	fmt.Printf(t.Input)
	if t.Pos > 0 && t.Pos == t.InputLen && (t.Pos+t.PromptLen)%t.Cols == 0 {
		fmt.Print("\n\r")
		rows++
		if rows > t.LinesUsed {
			t.LinesUsed = rows
		}
	}
	pos := (t.PromptLen + t.Pos + t.Cols) / t.Cols
	steps := rows - pos
	if steps > 0 {
		t.MoveCursorUp(steps)
	}
	col := (t.PromptLen + t.Pos) % t.Cols
	if col > 0 {
		t.MoveCursorRight(col)
	} else {
		fmt.Printf("\r")
	}
	t.OldPos = t.Pos
}

// EraseInput ...
func (t *Terminal) EraseInput() {
	t.Pos = 0
	t.Input = ""
	t.InputLen = 0
	t.LinesUsed = 1
}

func (t *Terminal) ClearRows() {
	rows := (t.PromptLen + t.InputLen + t.Cols - 1) / t.Cols
	rpos := (t.PromptLen + t.OldPos + t.Cols) / t.Cols
	oldRows := t.LinesUsed

	if rows > t.LinesUsed {
		t.LinesUsed = rows
	}

	// Goto last row
	steps := oldRows - rpos
	if steps > 0 {
		t.MoveCursorDown(steps)
	}

	// Now clear every row and go up
	for i := 0; i < oldRows-1; i++ {
		t.ResetLine()
		t.MoveCursorUp(1)
	}
	t.ResetLine()
}

// ToNextRow handles newline after enter
func (t *Terminal) ToNextRow() {
	fmt.Printf("\r")
	fmt.Printf("\n")
	t.EraseInput()
}

// New creates new screen instance
func New(clear bool) *Terminal {
	term := &Terminal{
		LinesUsed: 1,
		Cols:      80,
		Rows:      24,
		PromptLen: 2,
		PromptStr: "$>",
		Buffer:    new(bytes.Buffer),
	}
	if clear {
		term.Clear()
		term.MoveCursor(0, 0)
	}
	return term
}

// ResetAll custom styles
func (t *Terminal) ResetAll() {
	fmt.Printf("\033[0m")
}

// ResetColor to default color
func (t *Terminal) ResetColor() {
	fmt.Printf("\033[32m")
}

// ResetLine clears line and put cursor at beginning of it
func (t *Terminal) ResetLine() {
	fmt.Printf("\r\033[K")
}

// Clear Terminal
func (t *Terminal) Clear() {
	fmt.Printf("\033[2J")
}
