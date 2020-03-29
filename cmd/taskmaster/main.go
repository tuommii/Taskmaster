package main

import (
	"bufio"
	"fmt"
	"os"

	"miikka.xyz/linedit"
)

// Reset all custom styles
const resetAll = "\033[0m"

// Reset to default color
const resetColor = "\033[32m"

// Clear line and put cursor at beginning of it
const resetLine = "\r\033[K"

// Clear screen
func Clear() {
	fmt.Printf("\033[2J")
}

// MoveCursor to given position
func MoveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// MoveCursorForward relative the current position
func MoveCursorForward(bias int) {
	fmt.Printf("\033[%dC", bias)
}

// MoveCursorBackward ...
func MoveCursorBackward(bias int) {
	fmt.Printf("\033[%dD", bias)
}

// Move cursor down relative the current position
func MoveCursorDown(bias int) {
	fmt.Printf("\033[%dB", bias)
}

// Move cursor up relative the current position
func MoveCursorUp(bias int) {
	fmt.Printf("\033[%dA", bias)
}

func readKey(reader *bufio.Reader) (rune, int) {
	char, size, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char, size
}

var str = ""

func main() {
	backup, _ := linedit.Attr(os.Stdin)
	defer backup.Set(os.Stdin)
	tty := backup
	tty.Raw()
	tty.Set(os.Stdin)

	var b []byte = make([]byte, 5)

	Clear()
	// MoveCursor(0, 0)
	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		if code == 'x' {
			break
		} else if int(b[0]) >= int('A') && int(b[0]) < int('z') && code != 'y' {
			fmt.Print(resetLine)
			str += string(code)
			fmt.Print(str)
		} else if code >= 183 && code <= 186 {
			// cursor.MoveRight(10)
			fmt.Print(resetLine)
			MoveCursorForward(10)
			// fmt.Print(string(b))
		} else if code == 'y' {
			MoveCursor(0, 0)
		} else if code == 13 {
			MoveCursorDown(1)
			str = ""
			fmt.Print(resetLine)
		} else {
			fmt.Println(code)
		}
	}
}
