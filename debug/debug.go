// Package debug Watch debug file: tail -f /tmp/taskmaster_debug
package debug

import (
	"fmt"
	"os"

	"miikka.xyz/tty"
)

var file *os.File
var err error

// Open ...
func Open() {
	file, err = os.OpenFile("/tmp/taskmaster_debug", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Debugging failed!")
	}
}

func Write(win *tty.Terminal, input string) {
	file.WriteString(fmt.Sprintf("Hello form debug! POS: %d, STR: %s\n", win.X, input))
}

func Close() {
	file.Close()
}
