package job

import (
	"bufio"
	"encoding/json"
	"errors"
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

const maxRetries = 10

// Process represents runnable process
type Process struct {
	Name       string `json:"name"`
	Command    string `json:"command"`
	OutputLog  string `json:"stdout"`
	ErrorLog   string `json:"stderr"`
	WorkingDir string `json:"workingDir"`
	Procs      int    `json:"procs"`
	// Time when process is consired started
	StartTime    int    `json:"startTime"`
	StartRetries int    `json:"startRetries"`
	StopTime     int    `json:"stopTime"`
	StopSignal   string `json:"stopSignal"`
	Umask        int    `json:"umask"`
	Retries      int    `json:"retries"`
	Cmd          *exec.Cmd
	Started      time.Time
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
	for name, process := range tasks {
		process.Name = name
		// TODO: check support with config reloading
		process.Status = STOPPED
	}
	return tasks
}

// Launch executes a task
func (p *Process) Launch() error {
	if p.Status != STOPPED {
		return errors.New("Can't launch started process")
	}
	p.Status = STARTING
	p.prepare()
	oldMask := syscall.Umask(p.Umask)
	p.launch()
	p.killAfter()
	syscall.Umask(oldMask)
	p.clean()
	return nil
}

// Kill process
func (p *Process) Kill() error {
	return p.Cmd.Process.Kill()
	// Maybe ?
	// p.Cmd.Process.Release()
}

// TODO: Subject maybe means set running/started after x seconds
func (p *Process) launch() {
	err := p.Cmd.Start()
	if err != nil {
		fmt.Println("exec error", p.Name, err)
		p.Status = STOPPED
		// Move down if retries + 1 is wanted
		p.Retries--
		if p.Retries > 0 && p.Retries < maxRetries {
			p.launch()
		}
	}
	p.Started = time.Now()
	if p.StartTime <= 0 {
		p.Status = RUNNING
		return
	}
	timeoutCh := time.After(time.Duration(p.StartTime) * time.Second)
	go func() {
		<-timeoutCh
		// Do not set running if execution has failed
		if p.Status == STARTING {
			p.Status = RUNNING
			fmt.Println(p.Name, "is consired started")
		}
	}()
}

func (p *Process) killAfter() {
	if p.StopTime <= 0 {
		return
	}
	// add timestart also
	timeoutCh := time.After(time.Duration(p.StopTime)*time.Second + time.Duration(p.StartTime)*time.Second)
	go func() {
		<-timeoutCh
		p.Status = STOPPED
		fmt.Println(p.Name, "stopped")
		p.Kill()
	}()
}

// clean process when ready
func (p *Process) clean() {
	// Wait until process is done
	if p.Status != RUNNING {
		return
	}
	go func() {
		err := p.Cmd.Wait()
		if err != nil {
			p.Status = STOPPED
			fmt.Println("Error while executing program:", p.Name, err)
		}
	}()
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
	go p.redirect(p.stdout, p.OutputLog, os.Stdout)
	go p.redirect(p.stderr, p.ErrorLog, os.Stderr)
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
		fmt.Println("Error while opening log file:", err)
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
