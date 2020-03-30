package main

import (
	"flag"
	"fmt"
	"os"

	"miikka.xyz/debug"
	"miikka.xyz/tty"
)

// For testing
func list() {
	fmt.Println("Eka")
	fmt.Println("Toka")
	fmt.Println("Kolmas")
	fmt.Println("Neljas")
	fmt.Println("Viides")
}

func main() {

	debugFlag := flag.Bool("debug", false, "Write debug to file")
	flag.Parse()

	if *debugFlag {
		debug.Open()
		defer debug.Close()
	}

	// Dont touch this => restore terminal to same mode it was
	defaultMode, err := tty.GetMode(os.Stdin)
	if err != nil {
		fmt.Println("Can't read file mode!", err)
		os.Exit(1)
	}
	// Restoring
	defer defaultMode.UseTo(os.Stdin)

	activeMode := defaultMode
	activeMode.ToRaw()
	activeMode.UseTo(os.Stdin)

	var b []byte = make([]byte, 5)

	win := tty.New()
	win.Clear()
	win.MoveCursor(0, 0)
	win.Redraw()
	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		if code == 'x' {
			break
		} else if int(b[0]) >= 32 && int(b[0]) < 127 {
			// win.Input = append(win.Input, string(code))
			// win.InputLen++
			win.Input += string(b[0])
			win.Pos++
			win.Redraw()
		} else if code == 183 {
			win.MoveCursorUp(1)
		} else if code == 184 {
			win.MoveCursorDown(1)
		} else if code == 185 {
			win.MoveCursorRight(1)
		} else if code == 186 {
			win.MoveCursorLeft(1)
		} else if code == 13 {
			// win.MoveCursorDown(1)
			// win.ResetLine()
			// fmt.Printf("Input was: [%s]\n", str)
			win.EraseInput()
			win.Redraw()
			win.Redraw()
		}
		go debug.Write(win, win.Input, *debugFlag)
	}
}
