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

func Write(win *tty.Terminal, input string, debug bool) {
	if !debug {
		return
	}
	file.WriteString(fmt.Sprintf("POS: [%d], INPUT_LEN: [%d], KEY_PRESSED: [%d]\n", win.Pos, win.InputLen, win.KeyCode))
}

func Close() {
	file.Close()
}
