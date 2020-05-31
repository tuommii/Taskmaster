package cli

import (
	"fmt"

	"github.com/tuommii/taskmaster/job"
)

// HelpCmd implements help command
var HelpCmd = &Command{
	Name:  "help",
	Alias: "h",
	Usage: "Show help",
	Run:   help,
}

func help(cmd *Command, args []string, tasks map[string]*job.Process) {
	fmt.Println("\n\nHELP!")
	for _, arg := range args {
		fmt.Print("\nHelp ", arg)
	}
}
