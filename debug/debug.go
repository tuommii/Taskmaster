// Package debug Watch debug file: tail -f /tmp/taskmaster_debug
package debug

import (
	"fmt"
	"os"

	"taskmaster/tty"
)

const path = "/tmp/taskmaster_debug"

var (
	file *os.File
	err  error
)

// OpenFile opens file where data is written
func OpenFile() {
	file, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Debugging failed!")
	}
}

// Write writes data to file
func Write(win *tty.Terminal, input string, debug bool) {
	if !debug {
		return
	}
	file.WriteString(fmt.Sprintf("POS: [%d], INPUT_LEN: [%d], KEY_PRESSED: [%d]\n", win.Pos, win.InputLen, win.KeyCode))
}

// CloseFile closes file
func CloseFile() {
	file.Close()
}
