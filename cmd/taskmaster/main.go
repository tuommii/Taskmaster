package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"taskmaster/cli"
	"taskmaster/debug"
	"taskmaster/pad"
	"taskmaster/tty"
)

func parseInput(input string) {
	// taskmaster.Testi()
	if len(input) == 0 {
		return
	}
	tokens := strings.Split(input, " ")
	for _, cmd := range cli.Commands {
		if tokens[0] == cmd.Name {
			cmd.Run(cmd, tokens[1:])
		}
	}
}

func runCommand(args []string) {
	for _, name := range args {
		for _, cmd := range cli.Commands {
			if name == cmd.Name {
				// fmt.Printf(name)
				cmd.Run(cmd, []string{"Miikka"})
			}
		}
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

	win := tty.New(false)
	win.Prompt("$taskmaster>")

	var code int
	for {
		code = pad.KeyPressed()
		switch {
		// ESC
		case code == pad.Esc:
			return
		// IsPrintable
		case pad.Printable(code):
			if win.Pos == win.InputLen {
				win.Buffer.WriteRune(rune(code))
				win.ResetLine()
				fmt.Print(win.Buffer.String())
				win.Pos++
				win.InputLen++
			} else {
				// Insert
			}
		case code == pad.Enter:
			defaultMode.ApplyMode()
			bytes := win.Buffer.Bytes()
			input := string(bytes[win.PromptLen:])

			parseInput(input)
			clear(win)

			activeMode.ToRaw()
		}
		go debug.Write(win, win.Input, *debugFlag)
	}
}

func clear(win *tty.Terminal) {
	win.Buffer.Reset()
	win.ToNextRow()
	win.PrintPrompt()
}
