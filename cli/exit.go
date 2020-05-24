package cli

import (
	"fmt"
)

// ExitCmd implements help command
var ExitCmd = &Command{
	Name:  "exit",
	Alias: "quit",
	Usage: "Quit program",
	Run:   exit,
}

func exit(cmd *Command, args []string) {
	fmt.Println("\n\nSERVER EXIT!")
	for _, arg := range args {
		fmt.Print("\nEXIT! ", arg)
	}
}
