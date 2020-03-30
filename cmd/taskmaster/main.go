package main

import (
	"bufio"
	"fmt"
	"os"

	"miikka.xyz/linedit"
	"miikka.xyz/screen"
)

// Watch debug file: tail -f /tmp/taskmaster_debug
func debug(win *screen.Terminal, input string) {
	file, err := os.OpenFile("/tmp/taskmaster_debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Debugging failed!")
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("POS: %d, STR: %s\n", win.X, input))
}

func readKey(reader *bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char
}

var str = ""

func main() {
	backup, _ := linedit.Attr(os.Stdin)
	// Restore terminal to same mode it was
	defer backup.Set(os.Stdin)

	tty := backup
	tty.Raw()
	tty.Set(os.Stdin)

	var b []byte = make([]byte, 5)

	win := screen.New()
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
		debug(win, str)
	}
}
