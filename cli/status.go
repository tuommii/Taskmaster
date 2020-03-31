package cli

import "fmt"

// StatusCmd implements help command
var StatusCmd = &Command{
	Name:  "status",
	Usage: "Show status",
	Run:   status,
}

// Help ...
func status(cmd *Command, args []string) {
	fmt.Println("STATUS!")
	for _, arg := range args {
		fmt.Println("status", arg)
	}
}
