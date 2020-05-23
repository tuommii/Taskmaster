package cli

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

// StatusCmd implements help command
var StatusCmd = &Command{
	Name:  "status",
	Usage: "Show status",
	Alias: "stat",
	Run:   status,
}

// Help ...
func status(cmd *Command, args []string, t *terminal.Terminal) {
	fmt.Print("\nSTATUS!")
	for _, arg := range args {
		fmt.Print("\nstatus ", arg)
	}
}
