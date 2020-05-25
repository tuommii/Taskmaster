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
	"time"

	"github.com/tuommii/taskmaster"
	"github.com/tuommii/taskmaster/job"
	"golang.org/x/net/netutil"
)

func handleConnection(conn net.Conn, logger *log.Logger) {
	// read buffer from client after enter is hit
	data, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		fmt.Println("client left..")
		conn.Close()

		// escape recursion
		return
	}

	// convert bytes from buffer to string
	message := strings.Trim(string(data), "\n")
	taskmaster.RunCommand(taskmaster.ParseInput(message))
	// get the remote address of the client
	clientAddr := conn.RemoteAddr().String()
	// format a response
	response := fmt.Sprintf(message + " from " + clientAddr + "\n")

	// have server print out important information
	fmt.Println(response)

	// let the client know what happened
	conn.Write([]byte("you sent: " + response))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn, logger)
}

func main() {

	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent
	signalsCh := make(chan os.Signal, 1)
	// Passing no signals to Notify means that
	// all signals will be sent to the channel.
	// signal.Notify(signalsCh)
	signal.Ignore(syscall.SIGCHLD)
	go func() {
		s := <-signalsCh
		fmt.Println("GOT SIGNAL", s)
	}()

	logger := taskmaster.Logger()
	logger.SetPrefix("SERVER: " + logger.Prefix())

	pathFlag := flag.String("config", "./config.example.json", "path to config file")
	flag.Parse()
	tasks := job.LoadAll(*pathFlag)

	for key, task := range tasks {
		fmt.Println(key, task)
		task.Launch()
	}

	// Simulate killing
	time.Sleep(time.Second * 2)
	fmt.Println("2sec")
	signalsCh <- syscall.SIGALRM
	for _, task := range tasks {
		// task.Done <- true
		// task.Cmd.Process.Kill()
		task.Kill()
	}

	l, err := net.Listen("tcp", ":4200")
	if err != nil {
		logger.Fatal("LISTEN:", err)
	}
	defer l.Close()

	// Only one client at time allowed
	l = netutil.LimitListener(l, 1)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Fatal("ACCEPT", err)
		}
		go handleConnection(conn, logger)
	}
}
