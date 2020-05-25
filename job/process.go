package job

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	Name      string `json:"name"`
	Command   string `json:"command"`
	OutputLog string `json:"stdout"`
	ErrorLog  string `json:"stderr"`
	Cmd       *exec.Cmd
	Status    int
	stdout    io.ReadCloser
	stderr    io.ReadCloser
}

// LoadAll loads all jobs from config file
func LoadAll(path string) map[string]*Process {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error while opening config file: ", err)
		panic(err)
	}
	err = json.Unmarshal([]byte(file), &tasks)
	fmt.Printf("%+v\n", tasks)
	return tasks
}

// Launch ...
func (p *Process) Launch() {
	p.prepare()
	go p.redirect(p.stdout, p.OutputLog, os.Stdout)
	go p.redirect(p.stderr, p.ErrorLog, os.Stderr)
	p.Cmd.Start()
	go p.clean()
}

// clean process when ready
func (p *Process) clean() {
	// Wait until process is done
	err := p.Cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}
	// No need to call Close() when using pipes ?
	// p.stdout.Close()
	// p.stderr.Close()
}

// prepare command for executiom
func (p *Process) prepare() {
	tokens := strings.Fields(p.Command)
	p.Cmd = exec.Command(tokens[0], tokens[1:]...)

	var err error
	// TODO: err checks
	p.stdout, err = p.Cmd.StdoutPipe()
	if err != nil {
		fmt.Println("PIPE:", err)
	}
	p.stderr, err = p.Cmd.StderrPipe()
	if err != nil {
		fmt.Println("PIPE", err)
	}
}

// redirect standard stream to file. If path wasn't valid, then using alternative.
func (p *Process) redirect(stream io.ReadCloser, path string, alternative *os.File) {
	s := bufio.NewScanner(stream)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		file = alternative
		fmt.Println("FILE OPEN ERROR", err)
	}
	for s.Scan() {
		fmt.Fprintln(file, s.Text())
	}
	// When stream is closed this will executed
	which := "STDOUT"
	if stream == p.stderr {
		which = "STDERR"
	}
	fmt.Println(p.Name, which, "stopped")
}
