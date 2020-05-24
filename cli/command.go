package cli

import (
	"os"
	"sync"
)

// Commands holds all commands
var Commands []*Command

// Command represents command
type Command struct {
	Run func(cmd *Command, args []string)
	// Shown on available commands list
	Name string
	// Another string that runs same command
	Alias string
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

// GetExitStatus return exit status
func GetExitStatus() int {
	return exitStatus
}

// Exit calls all exit functions and exits with given status
func Exit() {
	for _, f := range atExitFuncs {
		f()
	}
	os.Exit(exitStatus)
}

// TODO: This could also be dynamic
// Add all commands to slice
func init() {
	Commands = append(Commands, HelpCmd)
	Commands = append(Commands, StatusCmd)
	Commands = append(Commands, StartCmd)
	Commands = append(Commands, ExitCmd)
}
