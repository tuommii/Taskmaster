package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Process statuses
const (
	STARTING = iota
	RUNNING
	STOPPED
	TIMEOUT
)

// Holds all processes
var tasks = make(map[string]*Process)

// Process represents runnable process
type Process struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Cmd     *exec.Cmd
	Status  int
}

// LoadAll jobs from config file
func LoadAll(path string) map[string]*Process {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error while opening config file: ", err)
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &tasks)
	return tasks
}

// GetStdout ... Remember close file
func (p *Process) GetStdout() *os.File {
	if p.Stdout == "" {
		fmt.Println("out was empty")
		return nil
	}
	file, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	p.Cmd.Stdout = file
	return file
}

// GetStderr ... Remember close file
func (p *Process) GetStderr() *os.File {
	if p.Stderr == "" {
		fmt.Println("err was empty")
		return nil
	}
	file, err := os.OpenFile(p.Stderr, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	p.Cmd.Stderr = file
	return file
}

// Launch ...
func (p *Process) Launch() {
	tokens := strings.Fields(p.Command)
	p.Cmd = exec.Command(tokens[0], tokens[1:]...)

	errFile := p.GetStderr()
	if errFile != nil {
		defer errFile.Close()
	}

	outFile := p.GetStdout()
	if outFile != nil {
		defer outFile.Close()
	}

	p.Cmd.Start()
	// stdout, err := p.Cmd.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// buf := bufio.NewReader(stdout) // Notice that this is not in a loop
	// num := 1
	// for {
	// 	line, _, err := buf.ReadLine()
	// 	if err != nil {
	// 		return
	// 	}
	// 	// if num > 3 {
	// 	// 	return
	// 	// }
	// 	num++
	// 	fmt.Println(string(line))
	// }
}
