package job

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Process statuses
const (
	STOPPED = iota
	STARTING
	RUNNING
	TIMEOUT
)

// Holds all processes
var tasks = make(map[string]*Process)

// Process represents runnable process
type Process struct {
	Name         string `json:"name"`
	Command      string `json:"command"`
	OutputLog    string `json:"stdout"`
	ErrorLog     string `json:"stderr"`
	WorkingDir   string `json:"workingDir"`
	Procs        int    `json:"procs"`
	StartTime    int    `json:"startTime"`
	StartRetries int    `json:"startRetries"`
	StopTime     int    `json:"stopTime"`
	StopSignal   string `json:"stopSignal"`
	Cmd          *exec.Cmd
	Status       int
	stdout       io.ReadCloser
	stderr       io.ReadCloser
}

// LoadAll loads all jobs from config file
func LoadAll(path string) map[string]*Process {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error while opening config file: ", err)
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

// Kill process
func (p *Process) Kill() error {
	return p.Cmd.Process.Kill()
	// p.Cmd.Process.Release()
}

// clean process when ready
func (p *Process) clean() {
	// Wait until process is done
	err := p.Cmd.Wait()
	// p.Cmd.ProcessState
	if err != nil {
		fmt.Println("CLEAN:", err)
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
	p.cwd(p.WorkingDir)
}

// Change current working directory if path exists and is directory
func (p *Process) cwd(dir string) {
	var stat os.FileInfo
	var err error

	if stat, err = os.Stat(p.WorkingDir); err != nil {
		return
	}
	if stat.IsDir() {
		p.Cmd.Dir = p.WorkingDir
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
