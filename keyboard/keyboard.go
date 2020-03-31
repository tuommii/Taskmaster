package keyboard

import "os"

// Keycodes
const (
	Tab       = 9
	Enter     = 10
	Esc       = 27
	Left      = 186
	Up        = 183
	Right     = 185
	Down      = 184
	Backspace = 127
	Delete    = 295
)

var buffer = make([]byte, 5)

// KeyPressed returns code for pressed key
func KeyPressed() int {
	var code int
	n, _ := os.Stdin.Read(buffer)
	for i := 0; i < n; i++ {
		code += int(buffer[i])
	}
	return code
}
