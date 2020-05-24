package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/tuommii/taskmaster"
	"github.com/tuommii/taskmaster/config"
	"github.com/tuommii/taskmaster/job"
	"golang.org/x/net/netutil"
)

func handleConnection(conn net.Conn, logger *log.Logger) {
	// read buffer from client after enter is hit
	data, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		logger.Println("client left..")
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
	logger.Println(response)

	// let the client know what happened
	conn.Write([]byte("you sent: " + response))

	// recursive func to handle io.EOF for random disconnects
	handleConnection(conn, logger)
}

func launch(executable string, args []string) {
	cmd := exec.Command(executable, args...)
	stdout, err := cmd.StdoutPipe()
	// stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	buf := bufio.NewReader(stdout) // Notice that this is not in a loop
	num := 1
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			return
		}
		// if num > 3 {
		// 	return
		// }
		num++
		fmt.Println(string(line))
	}
}

func main() {
	logger := taskmaster.Logger()
	logger.SetPrefix("SERVER: " + logger.Prefix())

	pathFlag := flag.String("config", "./config.example.json", "path to config file")
	flag.Parse()
	conf := config.LoadConfig(*pathFlag)

	for key, entry := range conf.Entries {
		p := job.Process{Name: key, Command: entry.Command}
		p.Launch()
	}

	l, err := net.Listen("tcp", ":4200")
	if err != nil {
		logger.Fatal("LISTEN:", err)
	}
	defer l.Close()

	l = netutil.LimitListener(l, 1)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Fatal("ACCEPT", err)
		}
		go handleConnection(conn, logger)

		// data = strings.TrimSpace(string(data))
		// taskmaster.RunCommand(taskmaster.ParseInput(data))
		// if data != "" {
		// 	fmt.Println("->", data)
		// }
		// conn.Write([]byte("from server"))
	}
}
