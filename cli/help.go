package cli

import "fmt"

// HelpCmd implements help command
var HelpCmd = &Command{
	Name:  "help",
	Usage: "Show help",
	Run:   help,
}

// Help ...
func help(cmd *Command, args []string) {
	fmt.Print("\nHELP!")
	for _, arg := range args {
		fmt.Print("\nHelp ", arg)
	}
}
