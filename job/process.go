package job

import (
	"bufio"
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
	Done    chan bool
	Status  int
}

// LoadAll loads all jobs from config file
func LoadAll(path string) map[string]*Process {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error while opening config file: ", err)
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &tasks)
	return tasks
}

// Launch ...
func (p *Process) Launch() {
	tokens := strings.Fields(p.Command)
	p.Cmd = exec.Command(tokens[0], tokens[1:]...)

	p.Done = make(chan bool, 1)
	out, _ := p.Cmd.StdoutPipe()
	// if err != nil {
	// 	log.Println("pipe error", err)
	// }
	// defer out.Close()
	go func() {
		s := bufio.NewScanner(out)
		file, err := os.OpenFile(p.Stdout, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("FILE OPEN ERROR", err)
			return
		}
		for s.Scan() {
			fmt.Fprintln(file, s.Text())
		}
		fmt.Println(p.Name, "stopped reading")
		// defer out.Close()
	}()

	// errFile := p.GetStderr()
	// if errFile != nil {
	// 	defer errFile.Close()
	// }

	// outFile := p.GetStdout()
	// if outFile != nil {
	// 	defer outFile.Close()
	// }

	p.Cmd.Start()
	go func() {
		// Wait until process is done
		err := p.Cmd.Wait()
		if err != nil {
			fmt.Println(err)
		}
		out.Close()
		// Inform
		p.Done <- true
	}()
	// p.Done <- true
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
