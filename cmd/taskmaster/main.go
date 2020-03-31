package main

import (
	"bytes"
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
	pos := 0

	win := tty.New()
	win.Clear()
	fmt.Println("Eka")
	win.Reposition()
	fmt.Println("Toka")
	win.Reposition()
	fmt.Println("Kolmas")
	win.Reposition()
	fmt.Println("Neljas")
	win.Reposition()
	fmt.Println("Viides")
	win.Reposition()

	buff := new(bytes.Buffer)

	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		if code == 'x' {
			break
		} else if code == 186 {
			// fmt.Fprintf(buff, "\033[%dD", 1)
			win.MoveCursorLeft(1)
			pos--
			// win.MoveCursorLeft(1)
		} else if int(b[0]) >= 32 && int(b[0]) < 127 {
			buff.WriteRune(rune(b[0]))
			win.ResetLine()
			fmt.Print(buff.String())
			pos++
			// win.Input = append(win.Input, string(code))
			// win.InputLen++
			// win.Input += string(b[0])
			// win.Pos++
			// win.Redraw()
		} else if code == 183 {
			win.MoveCursorUp(1)
		} else if code == 184 {
			win.MoveCursorDown(1)
		} else if code == 185 {
			win.MoveCursorRight(1)
		} else if code == 13 {
			buff.Reset()
			win.Reposition()
		}
		go debug.Write(win, win.Input, *debugFlag)
	}
}
