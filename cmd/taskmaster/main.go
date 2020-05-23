package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/debug"
	"github.com/tuommii/taskmaster/pad"
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

func runCommand(tokens []string) {
	if len(tokens) == 0 {
		return
	}
	for _, cmd := range cli.Commands {
		if tokens[0] == cmd.Name || tokens[0] == cmd.Alias {
			cmd.Run(cmd, tokens[1:])
		}
	}
}

func main() {
	debug.OpenFile()
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
	s.buf = s.buf[:0]
	// win.Prompt("(taskmaster$>) ")
	// loop(win, &defaultMode, &activeMode, *debugFlag)
	var code int
	fmt.Print(s.PromptStr)
	for {
		code = pad.KeyPressed()
		switch {
		case code == pad.Esc:
			return
		case pad.Printable(code):
			if s.Pos == s.InputLen {
				s.buf = append(s.buf, byte(code))

				fmt.Print("\r\033[K")
				fmt.Print(s.PromptStr)
				fmt.Print(string(s.buf))

				s.Pos++
				s.InputLen++
				// s.buf[s.Pos] = byte(code)
				// Clear line
				// win.Buffer.WriteRune(rune(code))
				// win.ResetLine()
				// fmt.Print(win.Buffer.String())
				// win.Pos++
				// win.InputLen++
			} else {
				// make space for a new char

				s.buf = append(s.buf, 'M')
				copy(s.buf[s.Pos+1:], s.buf[s.Pos:])
				s.buf[s.Pos] = byte(code)

				fmt.Print("\r\033[K")
				fmt.Print(s.PromptStr)
				fmt.Print(string(s.buf))

				s.Pos++
				s.InputLen++

				// Move cursor
				fmt.Print("\r")
				for i := 0; i < s.Pos+s.PromptLen; i++ {
					fmt.Print("\033[1C")
				}

				// Insert
			}
			debug.Write(s.buf, s.Pos)
		case code == pad.Enter:
			// backup.ApplyMode()

			// fmt.Printf("\n")
			// bytes := win.Buffer.Bytes()
			// input := string(bytes[win.PromptLen:])
			terminal.Restore(0, oldState)
			runCommand(parseInput(strings.Trim(string(s.buf), "\n")))
			terminal.MakeRaw(0)
			s.buf = s.buf[:0]
			s.InputLen = 0
			s.LinesUsed = 1
			s.Pos = 0
			fmt.Print("\n\r")
			// fmt.Print("\r\033[K")
			fmt.Print(s.PromptStr)
			// fmt.Print(string(s.buf))

			// fmt.Print(s.PromptStr)
			// clear(win)
			// curr.RawMode()
		case code == pad.Backspace:
			// win.Buffer.WriteRune(r	une('\b'))
			// win.Buffer.WriteRune(rune(' '))
			// win.Buffer.WriteRune(rune('\b'))
			// win.Pos--
			// win.MoveCursorLeft(1)
		case code == pad.Left:
			// fmt.Fprint(os.Stdin, MoveLeft(1))
			// fmt.Print("DSDS")
			if s.Pos > 0 {
				s.Pos--
				fmt.Print("\033[1D")
				// fmt.Print(string(s.buf))
			}
		case code == pad.Right:
			// fmt.Print("DSDS")
			if s.Pos < s.InputLen {
				s.Pos++
				fmt.Print("\033[1C")
			}
			// win.MoveCursorLeft(1)
		case code == pad.Right:
		case code == pad.Up:
			s.buf = s.buf[:0]
			s.buf = []byte("miikka")
			s.Pos = 6
			s.InputLen = 6
			fmt.Print("\r\033[K")
			fmt.Print(s.PromptStr)
			fmt.Print(string(s.buf))

			// win.Pos++
			// win.MoveCursorRight(1)
		}
		// debug.Write(win, win.Input, debugFlag)
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

func MoveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}
