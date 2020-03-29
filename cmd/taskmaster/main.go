package main

import (
	"bufio"
	"fmt"
	"os"

	"miikka.xyz/linedit"
)

func readKey(reader *bufio.Reader) (rune, int) {
	char, size, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char, size
}

func main() {
	backup, _ := linedit.Attr(os.Stdin)
	defer backup.Set(os.Stdin)
	tty := backup
	tty.Raw()
	tty.Set(os.Stdin)

	var b []byte = make([]byte, 5)

	for {
		n, _ := os.Stdin.Read(b)
		code := 0
		for i := 0; i < n; i++ {
			code += int(b[i])
		}
		if code == 'x' {
			break
		} else if int(b[0]) >= int('A') && int(b[0]) < int('z') {
			fmt.Print(string(code))
		} else if code >= 183 && code <= 186 {
			// cursor.MoveRight(10)
			fmt.Print(string(b))
		} else {
			fmt.Println(code)
		}
	}
}
