package main

import (
	"bufio"
	"flag"
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

func (s *server) launchAll() {
	for key, task := range s.tasks {
		fmt.Println(key, task)
		task.Launch()
	}
}

// reloader func, takes path to config and old processes
type reloaderF func(string, map[string]*job.Process)

func main() {
	configPath := flag.String("config", "./config.example.json", "path to config file")
	flag.Parse()
	// tasks := job.LoadAll(*configPath)
	s := newServer(*configPath, job.LoadAll(*configPath))
	s.listenSignals()
	s.launchAll()
	s.listenConnections()
}

func reloadConfig(configPath string, tasks map[string]*job.Process) {
	fmt.Println("RELOAD")
	// for key, task := range tasks {
	// 	fmt.Println("Killing and deleting", key)
	// 	task.Kill()
	// 	delete(tasks, key)
	// }
	// tasks = job.LoadAll(configPath)
	// for _, task := range tasks {
	// 	task.Launch()
	// }
	// fmt.Println("LOADED", len(tasks), "TASKS")
}

func (s *server) listenSignals() {
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent
	signalsCh := make(chan os.Signal, 1)
	// Passing no signals to Notify means that all
	// signals will be sent to the channel.
	signal.Notify(signalsCh)

	go func(tasks map[string]*job.Process) {
		for s := range signalsCh {
			switch {
			case s == syscall.SIGHUP:
				fmt.Println("RELOAD", len(tasks))
			case s == syscall.SIGTERM || s == syscall.SIGINT:
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
	conn.Write([]byte("you sent: " + "sended to client"))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn)
}
