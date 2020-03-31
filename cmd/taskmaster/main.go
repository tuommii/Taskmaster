package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"miikka.xyz/cli"
	"miikka.xyz/debug"
	"miikka.xyz/keyboard"
	"miikka.xyz/tty"
)

var helpCommand = &cli.Command{Name: "help"}
var statusCommand = &cli.Command{Name: "status"}

func help(cmd *cli.Command, args []string) {
	fmt.Println(cmd.Name)
	fmt.Println("Etkö tiedä mitä auttaminen on!")
	for _, arg := range args {
		fmt.Println(arg)
	}
}

func init() {
	helpCommand.Run = help
	statusCommand.Run = help
	cli.Commands = []*cli.Command{
		helpCommand,
		statusCommand,
	}
}

// TODO: return interface
func parseInput(input string) {
	fmt.Println("")
	arr := strings.Split(input, " ")
	for _, name := range arr {
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

	// var b []byte = make([]byte, 5)
	win := tty.New()
	win.Clear()
	win.MoveCursor(0, 0)

	win.Buffer.WriteString(win.Prompt)
	fmt.Print(win.Buffer.String())

	var code int
	for {
		code = keyboard.KeyPressed()
		switch {
		// ESC
		case code == keyboard.Esc:
			return
		// IsPrintable
		case keyboard.IsPrintable(code):
			if win.Pos == win.InputLen {
				win.Buffer.WriteRune(rune(code))
				win.ResetLine()
				fmt.Print(win.Buffer.String())
				win.Pos++
				win.InputLen++
			} else {
				// Insert
			}
		case code == keyboard.Enter:
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
	win.Reposition()
	win.Buffer.WriteString(win.Prompt)
	fmt.Print(win.Buffer.String())
}
