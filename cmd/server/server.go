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
		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn net.Conn) {
	var res string
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		fmt.Println("client left..")
		conn.Close()
		// escape recursion
		return
	}
	msg := strings.Trim(string(data), "\n")
	fmt.Println(msg)
	// cli.RunCommand(cli.ParseInput(msg), s.tasks)
	// get the remote address of the client
	// clientAddr := conn.RemoteAddr().String()
	// fmt.Println(msg, "from", clientAddr+"\n")

	if msg == "status" {
		res = ""
		for _, task := range s.tasks {
			res += task.Name + ",  " + task.Status + "\n"
		}
		conn.Write([]byte(res + "\n"))
	} else {
		conn.Write([]byte("FROM SERVER: " + msg + "\n"))
	}

	if msg == "fg" {
		fmt.Println("FOREGROUND")
		s.tasks["REALTIME"].SetForeground(true)
	} else if msg == "bg" {
		fmt.Println("BACKGROUND")
		s.tasks["REALTIME"].SetForeground(false)
	}

	// recursive func to handle io.EOF for random disconnects
	s.handleConnection(conn)
}
