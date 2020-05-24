package cli

import (
	"fmt"
)

// StatusCmd implements help command
var StatusCmd = &Command{
	Name:  "status",
	Usage: "Show status",
	Alias: "stat",
	Run:   status,
}

// Help ...
func status(cmd *Command, args []string) {
	fmt.Println("\nSTATUS!")
	for _, arg := range args {
		fmt.Println("\nstatus ", arg)
	}
}
