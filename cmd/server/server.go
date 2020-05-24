package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/tuommii/taskmaster"
	"golang.org/x/net/netutil"
)

func handleConnection(conn net.Conn, logger *log.Logger) {
	// read buffer from client after enter is hit
	bufferBytes, err := bufio.NewReader(conn).ReadBytes('\n')

	if err != nil {
		logger.Println("client left..")
		conn.Close()

		// escape recursion
		return
	}

	// convert bytes from buffer to string
	message := string(bufferBytes)
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

func main() {
	logger := taskmaster.Logger()
	logger.SetPrefix("SERVER: " + logger.Prefix() + " ")
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
