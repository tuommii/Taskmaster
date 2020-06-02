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
		// TODO: handle err
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

// Hot-reload config
func (s *server) reloadConfig() {
	fmt.Println("reloading config...")
	s.removeTasks()
	s.tasks = job.LoadAll(s.configPath)
	s.launchTasks()
	fmt.Println("Loaded", len(s.tasks), "tasks")
}

func (s *server) getJobNames() string {
	var names string
	for name := range s.tasks {
		names += name + "|"
	}
	return names[:len(names)-1]
}

func (s *server) listenSignals() {
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent
	signalsCh := make(chan os.Signal, 1)
	// Passing no signals to Notify means that all
	// signals will be sent to the channel.
	signal.Notify(signalsCh)

	go func(tasks map[string]*job.Process) {
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
	}(s.tasks)
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
		go s.handleConnection(conn)
	}
}

func (s *server) jobFound(name string) bool {
	if _, found := s.tasks[name]; found {
		return true
	}
	return false
}

func parseUserInput(data []byte) (string, string) {
	msg := strings.Trim(string(data), "\n")
	// TODO: remove
	fmt.Println(msg)
	input := strings.Split(msg, " ")
	cmd := input[0]
	var arg string
	if len(input) >= 2 {
		arg = input[1]
	}
	return cmd, arg
}

// TODO: make this DRY
func (s *server) switchCommand(cmd string, arg string, conn net.Conn) {
	switch {
	case cmd == "job_names":
		conn.Write([]byte(s.getJobNames()))
	case cmd == "help" || cmd == "h":
		conn.Write([]byte("help cmd"))
	case cmd == "status" || cmd == "st":
		conn.Write([]byte("status todo"))
	case cmd == "start" || cmd == "run":
		if !s.jobFound(arg) {
			conn.Write([]byte("job not found"))
			break
		}
		s.tasks[arg].Launch(false)
		conn.Write([]byte(arg + " started"))
	case cmd == "stop":
		if !s.jobFound(arg) {
			conn.Write([]byte("job not found"))
			break
		}
		s.tasks[arg].Kill()
		conn.Write([]byte(arg + " stopped"))
	case cmd == "restart":
		conn.Write([]byte("restart"))
	case cmd == "exit" || cmd == "quit":
		conn.Write([]byte("exit or quit"))
	case cmd == "fg":
		if !s.jobFound(arg) {
			conn.Write([]byte("job not found"))
			break
		}
		s.tasks[arg].SetForeground(true)
		conn.Write([]byte("attached " + arg + " output to stdout"))
	case cmd == "bg":
		if !s.jobFound(arg) {
			conn.Write([]byte("job not found"))
			break
		}
		// TODO: maybe validations
		s.tasks[arg].SetForeground(false)
		conn.Write([]byte("deattached " + arg + " output from stdout"))
	default:
		conn.Write([]byte("server received: " + cmd))
	}

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
	s.switchCommand(cmd, arg, conn)
	// recursive func to handle io.EOF for random disconnects
	s.handleConnection(conn)
}
