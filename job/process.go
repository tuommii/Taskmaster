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

// Redirects standard stream to file. If path aint valid, then using alternative.
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
	// When ReadCloser is closed this will executed
	fmt.Println(p.Name, "stopped reading")
}

// Launch ...
func (p *Process) Launch() {
	tokens := strings.Fields(p.Command)
	p.Cmd = exec.Command(tokens[0], tokens[1:]...)

	// TODO: err checks
	p.stdout, _ = p.Cmd.StdoutPipe()
	p.stderr, _ = p.Cmd.StderrPipe()
	go p.redirect(p.stdout, p.OutputLog, os.Stdout)
	go p.redirect(p.stderr, p.ErrorLog, os.Stderr)

	// go func() {
	// 	s := bufio.NewScanner(p.stdout)
	// 	file, err := os.OpenFile(p.StdoutLog, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		file = os.Stdout
	// 		fmt.Println("FILE OPEN ERROR", err)
	// 	}
	// 	for s.Scan() {
	// 		fmt.Fprintln(file, s.Text())
	// 	}
	// 	// When ReadCloser is closed this will executed
	// 	fmt.Println(p.Name, "stopped reading")
	// }()

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
		p.stdout.Close()
		p.stderr.Close()
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
