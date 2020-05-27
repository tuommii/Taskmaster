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
	"syscall"
	"time"
)

// Statuses
const (
	STOPPED = iota
	STARTING
	RUNNING
	TIMEOUT
)

// Process represents runnable process
type Process struct {
	Name         string `json:"name"`
	Command      string `json:"command"`
	OutputLog    string `json:"stdout"`
	ErrorLog     string `json:"stderr"`
	WorkingDir   string `json:"workingDir"`
	Procs        int    `json:"procs"`
	StartDelay   int    `json:"startDelay"`
	StartRetries int    `json:"startRetries"`
	StopTime     int    `json:"stopTime"`
	StopSignal   string `json:"stopSignal"`
	Umask        int    `json:"umask"`
	Cmd          *exec.Cmd
	StartTime    time.Time
	Status       int
	stdout       io.ReadCloser
	stderr       io.ReadCloser
}

// LoadAll loads all jobs from config file
func LoadAll(path string) map[string]*Process {
	tasks := make(map[string]*Process)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error while opening config file: ", err)
	}
	err = json.Unmarshal([]byte(file), &tasks)
	for name, proc := range tasks {
		proc.Name = name
	}
	fmt.Printf("%+v\n", tasks)
	return tasks
}

// Launch executes a task
func (p *Process) Launch() {
	p.prepare()
	go p.redirect(p.stdout, p.OutputLog, os.Stdout)
	go p.redirect(p.stderr, p.ErrorLog, os.Stderr)
	oldMask := syscall.Umask(p.Umask)
	// TODO: use same techniue than kill after
	p.launchAfter()
	// p.Cmd.SysProcAttr.Ctty
	// Not creating goroutine if no delay
	p.killAfter()
	syscall.Umask(oldMask)
	go p.clean()
}

// Kill process
func (p *Process) Kill() error {
	return p.Cmd.Process.Kill()
	// Maybe ?
	// p.Cmd.Process.Release()
}

// TODO: Subject maybe means set running/started after x seconds
func (p *Process) launchAfter() {
	p.Cmd.Start()
	p.StartTime = time.Now()
	if p.StartDelay <= 0 {
		return
	}
	timeoutCh := time.After(time.Duration(p.StartDelay) * time.Second)
	go func() {
		<-timeoutCh
		fmt.Println(p.Name, "is consired started")
	}()
}

func (p *Process) killAfter() {
	fmt.Println("STOPTIME:", p.StopTime)
	if p.StopTime <= 0 {
		return
	}
	timeoutCh := time.After(time.Duration(p.StopTime) * time.Second)
	go func() {
		<-timeoutCh
		fmt.Println(p.Name, "timed out")
		p.Kill()
	}()
}

// clean process when ready
func (p *Process) clean() {
	// Wait until process is done
	err := p.Cmd.Wait()
	if err != nil {
		fmt.Println("Error while executing program:", p.Name, err)
	}
	// Maybe some use for p.Cmd.ProcessState later ?
	// No need to call Close() when using pipes ?
	// p.stdout.Close()
	// p.stderr.Close()
}

// prepare command for executiom
func (p *Process) prepare() {
	tokens := strings.Fields(p.Command)
	p.Cmd = exec.Command(tokens[0], tokens[1:]...)

	var err error
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
// TODO: when ready maybe use /dev/null
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
