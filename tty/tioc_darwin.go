// +build darwin freebsd dragonfly openbsd netbsd

package tty

import "syscall"

const (
	getTermios = syscall.TIOCGETA
	setTermios = syscall.TIOCSETA
)
