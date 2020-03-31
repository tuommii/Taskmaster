package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miikka.xyz/debug"
	"miikka.xyz/keyboard"
	"miikka.xyz/tty"
)

var usage = `hello

testi1
testi2
`

type command interface {
	run(args []string)
}

type helpCmd struct {
	usage string
}

func (cmd *helpCmd) run(args []string) {
	cmd.usage = "sasasa"
	fmt.Printf("HEELELELELELEP!")
}

func fail(code int, msg string, args ...interface{}) {
	if code == 0 {
		fmt.Fprintf(os.Stdout, msg+"\n", args...)
	} else {
		fmt.Fprintf(os.Stderr, msg+"\n", args...)
	}
	os.Exit(code)
}

// TODO: return interface
func parseInput(input string) command {
	var help helpCmd
	fmt.Println("")
	arr := strings.Split(input, " ")
	switch arr[0] {
	case "help":
		return &help
	case "exit":
		fmt.Println("Exit!")
		return &help
	default:
		fmt.Println("Unknown command!")
		return &help
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

	// Dont edit this, instead restore terminal to same mode
	// than it was when exiting.
	defaultMode, err := tty.GetMode()
	if err != nil {
		fmt.Println("Can't read file mode!", err)
		os.Exit(1)
	}
	defer defaultMode.ApplyMode()

	// Take copy of users terminal mode and set it to raw mode
	activeMode := defaultMode
	activeMode.ToRaw()

	// var b []byte = make([]byte, 5)
	pos := 0
	len := 0

	win := tty.New()
	win.Clear()
	win.MoveCursor(0, 0)

	win.Buffer.WriteString(win.Prompt)
	fmt.Print(win.Buffer.String())

	var code int
	for {
		code = keyboard.KeyPressed()
		if code == 'x' {
			break
		} else if code == 186 {
			win.MoveCursorLeft(1)
			pos--
			// win.MoveCursorLeft(1)
		} else if code >= 32 && code < 127 {
			if pos == len {
				win.Buffer.WriteRune(rune(code))
				win.ResetLine()
				fmt.Print(win.Buffer.String())
				win.Pos++
				win.InputLen++
			} else {
			}
		} else if code == 183 {
			win.MoveCursorUp(1)
		} else if code == 184 {
			win.MoveCursorDown(1)
		} else if code == 185 {
			win.MoveCursorRight(1)
		} else if code == 13 {
			defaultMode.ApplyMode()
			bytes := win.Buffer.Bytes()
			input := string(bytes[win.PromptLen:])

			cmd := parseInput(input)
			cmd.run(os.Args)
			clear(win)

			activeMode.ToRaw()
		}
		go debug.Write(win, win.Input, *debugFlag)
	}
}

func clear(win *tty.Terminal) {
	win.Buffer.Reset()
	win.Reposition()
	win.Buffer.WriteString(win.Prompt)
	fmt.Print(win.Buffer.String())
}
