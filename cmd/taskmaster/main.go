package main

import (
	"bufio"
	"fmt"
	"os"

	"miikka.xyz/linedit"
)

func readKey(reader *bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println("Error reading key: ", err)
	}
	return char
}

func main() {
	backup, _ := linedit.Attr(os.Stdin)
	defer backup.Set(os.Stdin)
	tty := backup
	tty.Raw()
	tty.Set(os.Stdin)
	// cursor.MoveRight(10)

	reader := bufio.NewReader(os.Stdin)

	for {
		key := readKey(reader)
		if key == 'x' {
			break
		} else {
			fmt.Println(key)
		}
	}
}
