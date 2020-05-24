package job

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

var tasks = make(map[string]*Process)

// Process represents runnable process
type Process struct {
	Name    string    `json:"name"`
	Command string    `json:"command"`
	Cmd     *exec.Cmd `json:"cmd"`
}

// Load ...
func Load(path string) map[string]*Process {
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
	splitted := strings.Fields(p.Command)
	p.Cmd = exec.Command(splitted[0], splitted[1:]...)
	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	p.Cmd.Start()
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
