package cli

import (
	"os"
	"sync"
)

// Command ...
type Command struct {
	Run func(cmd *Command, args []string)
	// Shown on available commands list
	Name  string
	Usage string
	// Possible subcommands
	Commands []*Command
}

// Runnable test if command can be run
func (c *Command) Runnable() bool {
	return c.Run != nil
}

var atExitFuncs []func()
var exitStatus = 0
var exitMu sync.Mutex

// AtExit appends new function to be called when exit
func AtExit(f func()) {
	atExitFuncs = append(atExitFuncs, f)
}

func setExitStatus(n int) {
	exitMu.Lock()
	if exitStatus < n {
		exitStatus = n
	}
	exitMu.Unlock()
}

func GetExitStatus() int {
	return exitStatus
}

func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}
