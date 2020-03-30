// Package tty provides cursor and screen
// manipulation operations for VT-100 terminals
package tty

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Termios represents termios struct
type Termios unix.Termios

// IsTerminal returns if this is runned on valid tty
func IsTerminal(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

// ToRaw configures setting to raw mode...
func (t *Termios) ToRaw() {
	t.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	// t.Iflag &^= BRKINT |unix.ISTRIP | ICRNL | IXON // Stevens RAW
	t.Oflag &^= unix.OPOST
	t.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.Cflag &^= unix.CSIZE | unix.PARENB
	t.Cflag |= unix.CS8
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
}

// UseToterminal t attributes on file.
func (t *Termios) UseTo(file *os.File) error {
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(unix.TCSETS), uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// GetMode (terminal related) attributes from file.
func GetMode(file *os.File) (Termios, error) {
	var t Termios
	fd := file.Fd()
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(unix.TCGETS), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return t, errno
	}
	t.Ispeed &= unix.CBAUD | unix.CBAUDEX
	t.Ospeed &= unix.CBAUD | unix.CBAUDEX
	return t, nil
}

// Terminal represents terminal window
type Terminal struct {
	// Cursor x position
	X      int
	Width  int
	Height int
}

// New creates new screen instance
func New() *Terminal {
	term := &Terminal{}
	return term
}

// ResetAll custom styles
func (s *Terminal) ResetAll() {
	fmt.Printf("\033[0m")
}

// ResetColor to default color
func (s *Terminal) ResetColor() {
	fmt.Printf("\033[32m")
}

// ResetLine clears line and put cursor at beginning of it
func (s *Terminal) ResetLine() {
	fmt.Printf("\r\033[K")
}

// MoveCursor to given position
func (s *Terminal) MoveCursor(x int, y int) {
	fmt.Printf("\033[%d;%dH", y, x)
}

// MoveCursorRight steps
func (s *Terminal) MoveCursorRight(step int) {
	fmt.Printf("\033[%dC", step)
}

// MoveCursorLeft ...
func (s *Terminal) MoveCursorLeft(step int) {
	fmt.Printf("\033[%dD", step)
}

// MoveCursorDown steps
func (s *Terminal) MoveCursorDown(step int) {
	fmt.Printf("\033[%dB", step)
}

// MoveCursorUp steps
func (s *Terminal) MoveCursorUp(step int) {
	fmt.Printf("\033[%dA", step)
}

// Clear Terminal
func (s *Terminal) Clear() {
	fmt.Printf("\033[2J")
}
