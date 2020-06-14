package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/tuommii/taskmaster/cli"
	"github.com/tuommii/taskmaster/job"
	"golang.org/x/net/netutil"
)

type server struct {
	configPath string
	tasks      map[string]*job.Process
}

func newServer(configPath string, tasks map[string]*job.Process) *server {
	s := &server{
		configPath: configPath,
		tasks:      tasks,
	}
	return s
}

func (s *server) launchTasks() {
	for _, task := range s.tasks {
		if err := task.Launch(true); err != nil {
			fmt.Println(err)
		}
	}
}

func (s *server) removeTasks() {
	for key, task := range s.tasks {
		fmt.Println("Killing and deleting", key)
		err := task.Kill()
		if err != nil {
			log.Println(err)
		}
		delete(s.tasks, key)
	}
}

func (s *server) reloadConfig() {
	fmt.Println("reloading config...")
	s.removeTasks()
	s.tasks = job.LoadAll(s.configPath)
	s.launchTasks()
	fmt.Println("Loaded", len(s.tasks), "tasks")
}

func (s *server) listenSignals() {
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh)

	go func() {
		for sig := range signalsCh {
			switch {
			case sig == syscall.SIGHUP:
				s.reloadConfig()
			case sig == syscall.SIGTERM || sig == syscall.SIGINT:
				fmt.Printf("\nABORT!")
				os.Exit(0)
			default:
			}
		}
	}()
}

func (s *server) listenConnections() {
	l, err := net.Listen("tcp", ":4200")
	if err != nil {
		log.Fatal("LISTEN:", err)
	}
	defer l.Close()

	// Only one client at time allowed
	l = netutil.LimitListener(l, 1)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("ACCEPT", err)
		}
		fmt.Println("client connected")
		go s.handleConnection(conn)
	}
}

func parseUserInput(data []byte) (string, string) {
	msg := strings.Trim(string(data), "\n")
	fmt.Println(msg)
	input := strings.Split(msg, " ")
	cmd := input[0]
	var arg string
	if len(input) >= 2 {
		arg = input[1]
	}
	return cmd, arg
}

func (s *server) runCommand(cmd string, arg string, conn net.Conn) {
	// Special case, needs to know config path
	if cmd == "reload" {
		s.reloadConfig()
		conn.Write([]byte("config reloaded"))
		return
	}
	if command, found := cli.Commands[cmd]; found && command.Runnable != nil {
		conn.Write([]byte(command.Runnable(s.tasks, arg)))
		return
	}
	conn.Write([]byte("unknown command"))
}

func (s *server) handleConnection(conn net.Conn) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("client left..")
		conn.Close()
		// escape recursion
		return
	}

	cmd, arg := parseUserInput(data)
	s.runCommand(cmd, arg, conn)
	// recursive func to handle io.EOF for random disconnects
	s.handleConnection(conn)
}
