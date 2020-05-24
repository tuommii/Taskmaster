package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var actions = make(map[string]string)

func init() {
	actions["help"] = "HELP!"
}

func main() {
	l, err := net.Listen("tcp", ":4200")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		data = strings.TrimSpace(string(data))
		if data == "exit" {
			fmt.Println("EXIT command")
			return
		}
		fmt.Println("->", string(actions[data]))
		conn.Write([]byte("from server"))
	}
}
