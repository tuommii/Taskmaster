package cli

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// HelpCmd implements help command
var HelpCmd = &Command{
	Name:  "help",
	Usage: "Show help",
	Run:   help,
}

func help(cmd *Command, args []string, t *terminal.Terminal) {
	t.Write([]byte("\nHELP!"))
	fmt.Print("\n\nHELP!")
	fmt.Fprint(os.Stdin, "HELP!")
	for _, arg := range args {
		fmt.Print("\nHelp ", arg)
	}
}
