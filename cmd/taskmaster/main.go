package main

import (
	"log"
	"strings"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/tty"
	"golang.org/x/crypto/ssh/terminal"
)

func parseInput(input string) []string {
	// taskmaster.RealTimeExample()
	if len(input) == 0 {
		return nil
	}
	tokens := strings.Split(input, " ")
	return tokens
}

func runCommand(tokens []string) {
	if len(tokens) == 0 {
		return
	}
	for _, cmd := range cli.Commands {
		if tokens[0] == cmd.Name || tokens[0] == cmd.Alias {
			cmd.Run(cmd, tokens[1:])
		}
	}
}

func main() {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, oldState)

	term := tty.New(4096)
	for {
		input := term.ReadKey()
		if input == "exit" {
			break
		}
		terminal.Restore(0, oldState)
		if input != "" {
			runCommand(parseInput(input))
		}
		terminal.MakeRaw(0)
	}
}
