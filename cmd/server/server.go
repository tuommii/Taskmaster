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

func listenSignals() {
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent
	signalsCh := make(chan os.Signal, 1)
	signal.Notify(signalsCh)

	go func() {
		// Must be in loop or otherwise config can be reloaded only once
		for s := range signalsCh {
			switch {
			case s == syscall.SIGHUP:
				fmt.Println("RELOAD CONFIG!")
			case s == syscall.SIGINT:
				fmt.Printf("\nABORT!")
				os.Exit(0)
			case s == syscall.SIGTERM:
				fmt.Printf("\nABORT!")
				os.Exit(0)
			default:
			}
		}
	}()
}

func main() {
	listenSignals()
	// // We must use a buffered channel or risk missing the signal
	// // if we're not ready to receive when the signal is sent
	// signalsCh := make(chan os.Signal, 1)
	// // Passing no signals to Notify means that
	// // all signals will be sent to the channel.
	// signal.Notify(signalsCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGABRT)
	// // signal.Ignore(syscall.SIGCHLD)

	// go func() {
	// 	// Must be in loop or otherwise config can be reloaded only once
	// 	for {
	// 		s := <-signalsCh
	// 		if s == syscall.SIGHUP {
	// 			fmt.Println("RELOAD CONFIG!")
	// 		}
	// 		if s == syscall.SIGABRT || s == syscall.SIGTERM {
	// 			os.Exit(0)
	// 		}
	// 	}
	// }()

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

func handleConnection(conn net.Conn, logger *log.Logger) {
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
	handleConnection(conn, logger)
}
