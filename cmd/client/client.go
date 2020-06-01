package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tuommii/taskmaster/tty"
	"golang.org/x/crypto/ssh/terminal"
)

// Client is wrapper for client data
type Client struct {
	signals  chan os.Signal
	done     chan bool
	oldState *terminal.State
	term     *tty.State
	conn     net.Conn
}

// Create new app (taskmaster)
func Create() *Client {
	var err error
	client := &Client{
		signals: make(chan os.Signal, 1),
		done:    make(chan bool, 1),
	}
	signal.Notify(client.signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	client.oldState, err = terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	client.term = tty.New(4096)

	// autocompletion
	client.term.SetProposer(func(input string, jobNames []string) []string {
		var arr []string
		var result []string
		arr = append(arr, "help")
		arr = append(arr, "h")
		arr = append(arr, "status")
		arr = append(arr, "st")
		arr = append(arr, "reload")
		arr = append(arr, "start")
		arr = append(arr, "run")
		arr = append(arr, "stop")
		arr = append(arr, "exit")
		arr = append(arr, "quit")

		splitted := strings.SplitN(input, " ", 2)
		if len(splitted) >= 2 {
			// var names []string
			// names = append(names, "realtime")
			// names = append(names, "failing")
			// names = append(names, "test")
			for _, name := range jobNames {
				if strings.HasPrefix(name, splitted[1]) {
					result = append(result, name)
				}
			}
			return result
		}

		for _, item := range arr {
			if strings.HasPrefix(item, input) {
				result = append(result, item)
			}
		}
		return result
	})

	client.conn, err = net.Dial("tcp", "127.0.0.1:4200")
	if err != nil {
		// TODO: Messes ui (raw/normal)
		// log.Println(err)
	}
	return client
}

// ListenSignals ...
func (app *Client) ListenSignals() {
	sig := <-app.signals
	log.Println("RECEIVED:", sig)
	app.done <- true
}

// ReadInput reads input until exit command or terminating signal
func (app *Client) ReadInput() {
	fmt.Fprintf(app.conn, "secret_command_for_suggestions"+"\n")
	resp := make([]byte, 4096)
	n, err := app.conn.Read(resp)
	if err != nil {
		log.Println(err)
	}
	names := strings.Split(string(resp[:n]), "|")
	app.term.SetJobNames(names)
	for {
		reply := make([]byte, 1024)
		input := app.term.ReadKey(app.signals)
		terminal.Restore(0, app.oldState)
		switch {
		case input == "exit":
			app.done <- true
		case input != "":
			if app.conn == nil {
				fmt.Println("No connection to server...")
				terminal.MakeRaw(0)
				continue
			}
			// Send to server
			fmt.Fprintf(app.conn, input+"\n")
			terminal.Restore(0, app.oldState)
			_, err := app.conn.Read(reply)
			if err != nil {
				log.Println("Error reading reply", err)
				terminal.MakeRaw(0)
				continue
			}
			// Print server response
			fmt.Println(string(reply))
		}
		terminal.MakeRaw(0)
	}
}

// Quit restores terminal mode before exit
func (app *Client) Quit() {
	<-app.done
	terminal.Restore(0, app.oldState)
	log.Println("quit")
	os.Exit(1)
}
