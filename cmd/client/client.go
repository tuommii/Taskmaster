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

// ListenSignals and inform we are done
func (app *Client) ListenSignals() {
	sig := <-app.signals
	log.Println("signal received:", sig)
	app.done <- true
}

// send command that return job names to server
func (app *Client) getJobNames() {
	if app.conn == nil {
		return
	}
	fmt.Fprintf(app.conn, "job_names"+"\n")
	resp := make([]byte, 4096)
	n, err := app.conn.Read(resp)
	if err != nil {
		log.Println(err)
		return
	}
	// :n, otherwise last item len in array is width + rest of buffer
	app.term.SetJobNames(strings.Split(string(resp[:n]), "|"))
}

// ReadInput reads input until exit command or terminating signal
func (app *Client) ReadInput() {
	app.getJobNames()
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
			fmt.Fprintf(app.conn, input+"\n")
			terminal.Restore(0, app.oldState)
			_, err := app.conn.Read(reply)
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
	// var arr []string
	var result []string
	// Hack before commands are implemented
	// arr = append(arr, "help")
	// arr = append(arr, "h")
	// arr = append(arr, "status")
	// arr = append(arr, "st")
	// arr = append(arr, "reload")
	// arr = append(arr, "start")
	// arr = append(arr, "run")
	// arr = append(arr, "stop")
	// arr = append(arr, "exit")
	// arr = append(arr, "quit")
	// arr = append(arr, "fg")
	// arr = append(arr, "bg")

	splitted := strings.SplitN(input, " ", 2)
	if len(splitted) >= 2 {
		for _, name := range jobNames {
			if strings.HasPrefix(name, splitted[1]) {
				result = append(result, name)
			}
		}
		return result
	}

	for _, item := range commands {
		if strings.HasPrefix(item, input) {
			result = append(result, item)
		}
	}
	return result
}

// Quit restores terminal mode before exit
func (app *Client) Quit() {
	<-app.done
	terminal.Restore(0, app.oldState)
	log.Println("quit")
	os.Exit(1)
}
