package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/tty"
	"golang.org/x/crypto/ssh/terminal"
)

const path = "/tmp/taskmaster_log"

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
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(file, time.Now().String()[:19]+" ", 0)

	logger.Println("Starting....")

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		logger.Fatal(err)
	}
	defer terminal.Restore(0, oldState)

	term := tty.New(4096)
	for {
		input := term.ReadKey()
		if input == "exit" {
			logger.Println("exit command")
			break
		}
		terminal.Restore(0, oldState)
		if input != "" {
			runCommand(parseInput(input))
		}
		terminal.MakeRaw(0)
	}
	logger.Println("quit")
}
