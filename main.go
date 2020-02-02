package main

import (
	"bufio"
	"fmt"
	"os"

	"./config"
)

// Taskmaster holds all data
type Taskmaster struct {
	// struct
	config string
}

// NewTaskmanager returns new instance of Taskmaster
func NewTaskmanager(configFile string) *Taskmaster {
	return &Taskmaster{config: configFile}
}

func main() {
	config := config.LoadConfig("config.example.json")
	// fmt.Printf("%+v", config)
	fmt.Printf("%+v", config.Get("Hello2"))

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Fprintln(os.Stdin, line)
	}

}
