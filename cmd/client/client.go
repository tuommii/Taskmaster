package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/tuommii/taskmaster/tty"
	"golang.org/x/crypto/ssh/terminal"
)

type client struct {
	signals  chan os.Signal
	done     chan bool
	oldState *terminal.State
	term     *tty.State
	conn     net.Conn
}

func create() *client {
	var err error
	client := &client{
		signals: make(chan os.Signal, 1),
		done:    make(chan bool, 1),
	}
	signal.Notify(client.signals, syscall.SIGINT, syscall.SIGTERM)
	client.oldState, err = terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	client.term = tty.New(4096)
	client.term.SetProposer(autocompleter)
	client.conn, _ = net.Dial("tcp", "127.0.0.1:4200")
	return client
}

func (client *client) listenSignals() {
	sig := <-client.signals
	log.Println("signal received:", sig)
	client.done <- true
}

func (client *client) getJobNames() []string {
	if client.conn == nil {
		return nil
	}
	fmt.Fprintf(client.conn, "job_names"+"\n")
	resp := make([]byte, 4096)
	n, err := client.conn.Read(resp)
	if err != nil {
		log.Println(err)
		return nil
	}
	// :n, otherwise last item len is width + rest of buffer
	names := strings.Split(string(resp[:n]), "|")
	sort.Strings(names)
	return names
}

func (client *client) readInput() {
	client.term.SetJobNames(client.getJobNames())
	for {
		reply := make([]byte, 1024)
		input := client.term.ReadKey(client.signals)
		terminal.Restore(0, client.oldState)
		switch {
		case input == "exit":
			client.done <- true
		case input != "":
			if client.conn == nil {
				fmt.Println("No connection to server...")
				terminal.MakeRaw(0)
				continue
			}
			fmt.Fprintf(client.conn, input+"\n")
			terminal.Restore(0, client.oldState)
			_, err := client.conn.Read(reply)
			if err != nil {
				log.Println("Error reading reply", err)
				terminal.MakeRaw(0)
				continue
			}
			fmt.Println(string(reply))
		}
		terminal.MakeRaw(0)
	}
}

func autocompleter(input string, commands []string, jobNames []string) []string {
	if suggestions := possibleJobs(input, jobNames); suggestions != nil {
		return suggestions
	}
	return possibleCommands(input, commands)
}

func possibleJobs(input string, jobNames []string) []string {
	splitted := strings.SplitN(input, " ", 2)
	if len(input) < 2 {
		return nil
	}

	var result []string
	for _, name := range jobNames {
		if strings.HasPrefix(name, splitted[1]) {
			result = append(result, name)
		}
	}
	return result
}

func possibleCommands(input string, commands []string) []string {
	var result []string
	for _, item := range commands {
		if strings.HasPrefix(item, input) {
			result = append(result, item)
		}
	}
	return result
}

func (client *client) quit() {
	<-client.done
	terminal.Restore(0, client.oldState)
	log.Println("quit")
	os.Exit(1)
}
