package main

import (
	"bufio"
	"fmt"
	"os"

	"miikka.xyz/debug"
	"miikka.xyz/tty"
)

func readKey(reader *bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char
}

var str = ""

func main() {

	debug.Open()

	backup, err := tty.GetMode(os.Stdin)
	if err != nil {
		fmt.Println("Can't read file mode!", err)
		os.Exit(1)
	}
	// Restore terminal to same mode it was
	defer backup.UseTo(os.Stdin)

	activeMode := backup
	activeMode.ToRaw()
	activeMode.UseTo(os.Stdin)

	var b []byte = make([]byte, 5)

	win := tty.New()
	win.Clear()
	win.MoveCursor(0, 0)

	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		if code == 'x' {
			break
		} else if int(b[0]) >= int('A') && int(b[0]) < int('z') && code != 'y' {
			win.ResetLine()
			str += string(code)
			fmt.Print(str)
		} else if code == 183 {
			win.MoveCursorUp(1)
		} else if code == 184 {
			win.MoveCursorDown(1)
		} else if code == 185 {
			win.MoveCursorRight(1)
		} else if code == 186 {
			win.MoveCursorLeft(1)
		} else if code == 13 {
			win.MoveCursorDown(1)
			str = ""
			win.ResetLine()
		} else {
			fmt.Println(code)
		}
		go debug.Write(win, str)
	}
	debug.Close()
}
