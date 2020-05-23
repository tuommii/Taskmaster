// Package debug Watch debug file: tail -f /tmp/taskmaster_debug
package debug

import (
	"fmt"
	"os"
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
func Write(buf []byte, pos int) {
	file.WriteString(fmt.Sprintf("BUF: [%s], POS: [%d]\n", string(buf), pos))
}

// CloseFile closes file
func CloseFile() {
	file.Close()
}
