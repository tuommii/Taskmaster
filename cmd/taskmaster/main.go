package main

import (
	"flag"
	"fmt"
	"os"

	"miikka.xyz/debug"
	"miikka.xyz/keyboard"
	"miikka.xyz/tty"
)

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

			handleInput(input)
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

func handleInput(input string) {
	if input == "miikka" {

		fmt.Printf("Dmksajka\thaajaj\n\n\nsfdfdfsd\t\n")

	}
}
