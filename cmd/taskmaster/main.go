package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/debug"
	"github.com/tuommii/taskmaster/pad"
	"github.com/tuommii/taskmaster/tty"
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

func init() {
	if !tty.IsSupported(os.Stdin.Fd()) {
		fmt.Print("OS not supported!")
		os.Exit(1)
	}
}

func main() {

	// Debug to file if flags is set. In Makefile this flag is present
	debugFlag := flag.Bool("debug", false, "Write debug to file")
	flag.Parse()
	if *debugFlag {
		debug.OpenFile()
		defer debug.CloseFile()
	}

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, oldState)

	// Dont edit this, instead restore terminal to same mode
	// than it was when exiting.
	// defaultMode, err := tty.GetMode()
	// if err != nil {
	// 	fmt.Println("Can't read file mode!", err)
	// 	os.Exit(1)
	// }
	// defer defaultMode.ApplyMode()

	// Take copy of users terminal mode and set it to raw mode
	// activeMode := defaultMode
	// activeMode.RawMode()

	s := &State{
		LinesUsed: 1,
		Cols:      80,
		Rows:      24,
		PromptLen: 2,
		PromptStr: "$>",
		Buffer:    new(bytes.Buffer),
		buf:       make([]byte, 4096),
	}
	// win.Prompt("(taskmaster$>) ")
	// loop(win, &defaultMode, &activeMode, *debugFlag)
	var code int
	for {
		code = pad.KeyPressed()
		switch {
		case code == pad.Esc:
			return
		case pad.Printable(code):
			if s.Pos == s.InputLen {
				s.buf = append(s.buf, byte(code))
				// s.buf[s.Pos] = byte(code)
				s.Pos++
				s.InputLen++
				// Clear line
				fmt.Print("\r\033[K")
				fmt.Print(string(s.buf))
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
			s.Pos = 0
			s.InputLen = 0
			s.LinesUsed = 1
			fmt.Print("\n\r")
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
			// win.Pos++
			// win.MoveCursorRight(1)
		}

		// go debug.Write(win, win.Input, debugFlag)
	}

}

func loop(win *tty.Terminal, backup *tty.Termios, curr *tty.Termios, debugFlag bool) {
	var code int
	for {
		code = pad.KeyPressed()
		switch {
		case code == pad.Esc:
			return
		case pad.Printable(code):
			if win.Pos == win.InputLen {
				win.Buffer.WriteRune(rune(code))
				win.ResetLine()
				fmt.Print(win.Buffer.String())
				win.Pos++
				win.InputLen++
			} else {
				buf := bytes.NewBuffer(win.Buffer.Bytes()[win.Pos:win.Pos])
				win.Buffer = buf
				buf.WriteRune(rune(code))
				// Insert
			}
		case code == pad.Enter:
			backup.ApplyMode()

			// fmt.Printf("\n")
			bytes := win.Buffer.Bytes()
			input := string(bytes[win.PromptLen:])
			runCommand(parseInput(strings.Trim(input, "\n")))
			clear(win)
			curr.RawMode()
		case code == pad.Backspace:
			// win.Buffer.WriteRune(rune('\b'))
			// win.Buffer.WriteRune(rune(' '))
			// win.Buffer.WriteRune(rune('\b'))
			win.Pos--
			win.MoveCursorLeft(1)
		case code == pad.Left:
			win.Pos--
			win.MoveCursorLeft(1)
		case code == pad.Right:
			win.Pos++
			win.MoveCursorRight(1)
		}

		go debug.Write(win, win.Input, debugFlag)
	}
}

func clear(win *tty.Terminal) {
	win.Buffer.Reset()
	win.ToNextRow()
	win.PrintPrompt()
}
