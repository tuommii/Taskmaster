package tty

import "os"

// Keycodes
const (
	CtrlC     = 3
	CtrlD     = 4
	Tab       = 9
	Enter     = 13
	Esc       = 27
	Left      = 186
	Up        = 183
	Right     = 185
	Down      = 184
	Backspace = 127
	Delete    = 295
)

var buffer = make([]byte, 5)

func isPrintable(code int) bool {
	if code >= 32 && code < 127 {
		return true
	}
	return false
}

// keyPressed returns code for a pressed key
func keyPressed() int {
	var code int
	n, _ := os.Stdin.Read(buffer)
	for i := 0; i < n; i++ {
		code += int(buffer[i])
	}
	return code
}
