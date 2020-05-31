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
	tasks      job.Tasks
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
		if err := task.Launch(); err != nil {
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
	fmt.Println("reload config...")
	s.removeTasks()
	s.tasks = job.LoadAll(s.configPath)
	s.launchTasks()
	fmt.Println("Loaded", len(s.tasks), "tasks")
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

func (s *server) handleConnection(conn net.Conn) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("client left..")
		conn.Close()
		// escape recursion
		return
	}
	msg := strings.Trim(string(data), "\n")
	input := strings.Split(msg, " ")
	fmt.Println(msg)

	// cli.RunCommand(cli.ParseInput(msg), s.tasks)
	// get the remote address of the client
	// clientAddr := conn.RemoteAddr().String()
	// fmt.Println(msg, "from", clientAddr+"\n")
	cmd := input[0]
	var arg string
	if len(input) >= 2 {
		arg = input[1]
	}
	switch {
	case cmd == "help" || cmd == "h":
		conn.Write([]byte("help cmd"))
	case cmd == "status" || cmd == "st":
		conn.Write([]byte(s.tasks.Status()))
	case cmd == "start" || cmd == "run":
		if s.jobFound(arg) {
			if s.tasks[arg].Status == job.LOADED || s.tasks[arg].Status == job.STOPPED {
				s.tasks[arg].Status = job.STOPPED
				s.tasks[arg].Launch()
				conn.Write([]byte("started"))
			}
		} else {
			conn.Write([]byte("job not found"))
		}
		break
	case cmd == "stop":
		if s.jobFound(arg) {
			s.tasks[arg].Kill()
			conn.Write([]byte("placeholder"))
		} else {
			conn.Write([]byte("job not found"))
		}
		break
	case cmd == "restart":
		conn.Write([]byte("restart"))
	case cmd == "exit" || cmd == "quit":
		conn.Write([]byte("exit or quit"))
	case cmd == "fg":
		if s.jobFound(arg) {
			s.tasks[arg].SetForeground(true)
			conn.Write([]byte("fore"))
		} else {
			conn.Write([]byte("job not found"))
		}
		break
	case cmd == "bg":
		if s.jobFound(arg) {
			if s.tasks[arg].Status != job.RUNNING {
				conn.Write([]byte("aint running"))
			} else {
				s.tasks[arg].SetForeground(false)
				conn.Write([]byte("foreground"))
			}
		} else {
			conn.Write([]byte("job not found"))
		}
		break
	default:
		conn.Write([]byte("server received: " + cmd))
	}

	// recursive func to handle io.EOF for random disconnects
	s.handleConnection(conn)
}
