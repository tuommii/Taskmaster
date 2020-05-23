package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/tuommii/taskmaster/cli"
<<<<<<< HEAD
	"github.com/tuommii/taskmaster/pad"
=======
	"github.com/tuommii/taskmaster/debug"
	"github.com/tuommii/taskmaster/tty"
>>>>>>> 8335bd957d112732563dba8accfb1d937e192323
	"golang.org/x/crypto/ssh/terminal"
)

type state struct {
	history      []string
	historyMutex sync.RWMutex
	cols         int
	r            *bufio.Reader
	cursorRows   int
}

// State ...
type State struct {
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
	buf       []byte
}

func parseInput(input string) []string {
	// taskmaster.RealTimeExample()
	if len(input) == 0 {
		return nil
	}
	tokens := strings.Split(input, " ")
	return tokens
}

func runCommand(tokens []string, t *terminal.Terminal) {
	if len(tokens) == 0 {
		return
	}
	for _, cmd := range cli.Commands {
		if tokens[0] == cmd.Name || tokens[0] == cmd.Alias {
			cmd.Run(cmd, tokens[1:], t)
		}
	}
}

func main() {

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, oldState)

	s := &State{
		LinesUsed: 1,
		Cols:      80,
		Rows:      24,
		PromptLen: 2,
		PromptStr: "$>",
		Pos:       0,
		Buffer:    new(bytes.Buffer),
		buf:       make([]byte, 4096),
	}
	// win.Prompt("(taskmaster$>) ")
	// loop(win, &defaultMode, &activeMode, *debugFlag)
	var code int
	fmt.Print("$>")
	s.Pos = 2
	for {
		code = pad.KeyPressed()
		switch {
		case code == pad.Esc:
			return
		case pad.Printable(code):
			if s.Pos == s.InputLen+s.PromptLen {
				s.buf = append(s.buf, byte(code))
				// s.buf[s.Pos] = byte(code)
				s.Pos++
				s.InputLen++
				// Clear line
				// win.Buffer.WriteRune(rune(code))
				// win.ResetLine()
				// fmt.Print(win.Buffer.String())
				// win.Pos++
				// win.InputLen++
			} else {
				// Insert
			}
		case code == pad.Enter:
			// backup.ApplyMode()

			// fmt.Printf("\n")
			// bytes := win.Buffer.Bytes()
			// input := string(bytes[win.PromptLen:])
			// runCommand(parseInput(strings.Trim(string(s.buf), "\n")))
			s.buf = s.buf[:0]
			s.InputLen = 0
			s.LinesUsed = 1
			s.Pos = s.PromptLen
			fmt.Print("\n\r")
			fmt.Print(s.PromptStr)
			// clear(win)
			// curr.RawMode()
		case code == pad.Backspace:
			// win.Buffer.WriteRune(rune('\b'))
			// win.Buffer.WriteRune(rune(' '))
			// win.Buffer.WriteRune(rune('\b'))
			// win.Pos--
			// win.MoveCursorLeft(1)
		case code == pad.Left:
			// win.Pos--
			// win.MoveCursorLeft(1)
		case code == pad.Right:
		case code == pad.Up:
			s.buf = s.buf[:0]
			s.buf = []byte("miikka")
			// win.Pos++
			// win.MoveCursorRight(1)
		}
		fmt.Print("\r\033[K")
		fmt.Print(s.PromptStr, string(s.buf))
		// go debug.Write(win, win.Input, debugFlag)
	}

}

// func loop(win *tty.Terminal, backup *tty.Termios, curr *tty.Termios, debugFlag bool) {
// 	var code int
// 	for {
// 		code = pad.KeyPressed()
// 		switch {
// 		case code == pad.Esc:
// 			return
// 		case pad.Printable(code):
// 			if win.Pos == win.InputLen {
// 				win.Buffer.WriteRune(rune(code))
// 				win.ResetLine()
// 				fmt.Print(win.Buffer.String())
// 				win.Pos++
// 				win.InputLen++
// 			} else {
// 				buf := bytes.NewBuffer(win.Buffer.Bytes()[win.Pos:win.Pos])
// 				win.Buffer = buf
// 				buf.WriteRune(rune(code))
// 				// Insert
// 			}
// 		case code == pad.Enter:
// 			backup.ApplyMode()

// 			// fmt.Printf("\n")
// 			bytes := win.Buffer.Bytes()
// 			input := string(bytes[win.PromptLen:])
// 			runCommand(parseInput(strings.Trim(input, "\n")))
// 			clear(win)
// 			curr.RawMode()
// 		case code == pad.Backspace:
// 			// win.Buffer.WriteRune(rune('\b'))
// 			// win.Buffer.WriteRune(rune(' '))
// 			// win.Buffer.WriteRune(rune('\b'))
// 			win.Pos--
// 			win.MoveCursorLeft(1)
// 		case code == pad.Left:
// 			win.Pos--
// 			win.MoveCursorLeft(1)
// 		case code == pad.Right:
// 			win.Pos++
// 			win.MoveCursorRight(1)
// 		}

// 		go debug.Write(win, win.Input, debugFlag)
// 	}
// }

// func clear(win *tty.Terminal) {
// 	win.Buffer.Reset()
// 	win.ToNextRow()
// 	win.PrintPrompt()
// }
