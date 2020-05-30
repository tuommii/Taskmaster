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

	"github.com/tuommii/taskmaster"
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
		task.Launch()
	}
}

func (s *server) removeTasks() {
	for key, task := range s.tasks {
		fmt.Println("Killing and deleting", key)
		task.Kill()
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("client left..")
		conn.Close()
		// escape recursion
		return
	}
	msg := strings.Trim(string(data), "\n")
	taskmaster.RunCommand(taskmaster.ParseInput(msg))
	// get the remote address of the client
	clientAddr := conn.RemoteAddr().String()
	fmt.Println(msg, "from", clientAddr+"\n")
	conn.Write([]byte(msg + " sended to client"))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn)
}
