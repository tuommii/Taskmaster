package main

import (
	"flag"
	"fmt"
	"os"

	"miikka.xyz/debug"
	"miikka.xyz/tty"
)

func main() {

	debugFlag := flag.Bool("debug", false, "Write debug to file")
	flag.Parse()

	if *debugFlag {
		debug.Open()
		defer debug.Close()
	}

	// Dont edit this => restore terminal to same mode it was
	defaultMode, err := tty.GetMode()
	if err != nil {
		fmt.Println("Can't read file mode!", err)
		os.Exit(1)
	}
	// Restoring
	defer defaultMode.ApplyMode()

	activeMode := defaultMode
	activeMode.ToRaw()

	var b []byte = make([]byte, 5)
	pos := 0
	len := 0

	win := tty.New()
	win.Clear()
	win.MoveCursor(0, 0)
	// fmt.Println("Eka")
	// win.Reposition()
	// fmt.Println("Toka")
	// win.Reposition()
	// fmt.Println("Kolmas")
	// win.Reposition()
	// fmt.Println("Neljas")
	// win.Reposition()
	// fmt.Println("Viides")
	// win.Reposition()

	// PROMPT
	win.Buffer.WriteString(win.Prompt)
	fmt.Print(win.Buffer.String())
	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		win.KeyCode = code
		if code == 'x' {
			break
		} else if code == 186 {
			// fmt.Fprintf(buff, "\033[%dD", 1)
			win.MoveCursorLeft(1)
			pos--
			// win.MoveCursorLeft(1)
		} else if int(b[0]) >= 32 && int(b[0]) < 127 {
			if pos == len {
				win.Buffer.WriteRune(rune(b[0]))
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
			bytes := win.Buffer.Bytes()
			input := string(bytes[win.PromptLen:])

			if input == "miikka" {
				defaultMode.ApplyMode()

				fmt.Printf("Dmksajka\thaajaj\n\n\nsfdfdfsd\t\n")

				activeMode.ToRaw()
			}

			// fmt.Printf("\n\rINPUT WAS:%s\n\r", input)
			win.Buffer.Reset()
			win.Reposition()
			win.Buffer.WriteString(win.Prompt)
			fmt.Print(win.Buffer.String())
		}
		go debug.Write(win, win.Input, *debugFlag)
	}
}
