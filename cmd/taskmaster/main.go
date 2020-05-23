package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	ch := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		logger.Fatal(err)
	}

	term := tty.New(4096)

	go func() {
		sig := <-ch
		logger.Println("RECEIVED:", sig)
		done <- true
	}()

	go func() {
		for {
			input := term.ReadKey(ch)
			if input == "exit" {
				logger.Println("exit command")
				break
			}
			if input != "" {
				terminal.Restore(0, oldState)
				runCommand(parseInput(input))
			}
			terminal.MakeRaw(0)
		}
		done <- true
	}()

	<-done
	terminal.Restore(0, oldState)
	logger.Println("quit")
	os.Exit(1)
}
