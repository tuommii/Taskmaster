package tty

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Termios represents termios struct
type Termios unix.Termios

// IsSupported returns if this is runned on valid OS
func IsSupported(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TCGETS)
	return err == nil
}

// RawMode configures setting to raw mode...
func (t *Termios) RawMode() {
	t.Iflag &^= unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON
	t.Oflag &^= unix.OPOST
	t.Lflag &^= unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN
	t.Cflag &^= unix.CSIZE | unix.PARENB
	t.Cflag |= unix.CS8
	// Block until 1 byte is read
	t.Cc[unix.VMIN] = 1
	t.Cc[unix.VTIME] = 0
	t.ApplyMode()
}

// ApplyMode activates current config to STDIN
func (t *Termios) ApplyMode() error {
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(os.Stdin.Fd()), uintptr(unix.TCSETS), uintptr(unsafe.Pointer(t)))
	if errno != 0 {
		return errno
	}
	return nil
}

// GetMode returns current config.
func GetMode() (Termios, error) {
	var t Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(os.Stdin.Fd()), uintptr(unix.TCGETS), uintptr(unsafe.Pointer(&t)))
	if errno != 0 {
		return t, errno
	}
	t.Ispeed &= unix.CBAUD | unix.CBAUDEX
	t.Ospeed &= unix.CBAUD | unix.CBAUDEX
	return t, nil
}
