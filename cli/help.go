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
	fmt.Println("HELP!")
	for _, arg := range args {
		fmt.Println("Help", arg)
	}
}
