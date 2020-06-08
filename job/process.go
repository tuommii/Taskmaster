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
	LOADED   = "LOADED"
	STARTING = "STARTING"
	RUNNING  = "RUNNING"
	STOPPED  = "STOPPED"
	FAILED   = "FAILED"
)

// Even config file has more, this is max
const maxRetries = 10

type options struct {
	Name string `json:"name"`
	// Command to execute, with arguments
	Command string `json:"command"`
	// Log files
	OutputLog string `json:"stdout"`
	ErrorLog  string `json:"stderr"`
	// Run command when config is loaded
	AutoStart  bool   `json:"autostart"`
	WorkingDir string `json:"workingDir"`
	// How many instances is launched
	Procs int `json:"procs"`
	// Time when process is consired started
	StartTime int `json:"startTime"`
	// After StopTime task quits. Counted from StartTime
	StopTime   int    `json:"stopTime"`
	StopSignal string `json:"stopSignal"`
	Umask      int    `json:"umask"`
	// Max tries to start a task
	Retries int `json:"retries"`
	// If process exits in any other way than whit stop request
	ExitCodes []int `json:"exitCodes"`
}

// Process represents runnable process
type Process struct {
	options
	IsForeground bool
	Cmd          *exec.Cmd
	Started      time.Time
	Status       string
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
	for name, task := range tasks {
		task.Name = name
		task.Status = LOADED
	}
	return tasks
}

func (p *Process) validateLauch(launchAutostartOnly bool) error {
	if launchAutostartOnly == true && p.Status == LOADED && p.AutoStart == false {
		return errors.New(p.Name + " loaded, but not started")
	}
	if p.Status == RUNNING {
		return errors.New("Can't launch a already started process")
	}
	return nil
}

// Launch executes a task
func (p *Process) Launch(launchAutostartOnly bool) error {
	if err := p.validateLauch(launchAutostartOnly); err != nil {
		return err
	}
	p.prepare()
	oldMask := syscall.Umask(p.Umask)
	if err := p.execute(); err != nil {
		fmt.Println(p.Name, p.Status, err)
		return err
	}
	p.killAfter()
	syscall.Umask(oldMask)
	p.clean()
	return nil
}

// Kill process
func (p *Process) Kill() error {
	if p.Status != RUNNING {
		return errors.New(p.Name + " wasn't running")
	}
	p.Status = STOPPED
	return p.Cmd.Process.Kill()
}

func (p *Process) execute() error {
	if err := p.Cmd.Start(); err != nil {
		if err := p.relaunch(err); err != nil {
			return err
		}
	}
	p.setStarted()
	return nil
}

// Try launch as long as retries left
func (p *Process) relaunch(err error) error {
	p.Status = FAILED
	p.Retries--
	if p.Retries > 0 && p.Retries < maxRetries {
		fmt.Println("Trying launch", p.Name, "again...")
		p.execute()
	}
	return err
}

func (p *Process) setStarted() {
	// No delay
	if p.StartTime <= 0 {
		p.Status = RUNNING
		fmt.Println(p.Name, p.Status)
		return
	}
	// Delay
	timeoutCh := time.After(time.Duration(p.StartTime) * time.Second)
	go func() {
		<-timeoutCh
		// Do not set running if execution has failed
		if p.Status != STARTING {
			return
		}
		p.Status = RUNNING
		p.Started = time.Now()
		fmt.Println(p.Name, "is consired started", p.Status)
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
		if err := p.Kill(); err != nil {
			return
		}
		fmt.Println(p.Name, "stopped")
	}()
}

func (p *Process) properExit(code int) bool {
	for _, val := range p.ExitCodes {
		if val == code {
			return true
		}
	}
	return false
}

// clean process when ready
func (p *Process) clean() {
	if p.Status != RUNNING {
		return
	}
	go func() {
		err := p.Cmd.Wait()
		if err == nil {
			return
		}
		p.Status = STOPPED
		code := p.Cmd.ProcessState.ExitCode()
		if p.properExit(code) {
			fmt.Println("EXITED WITH PROPER CODE:", code)
			return
		}
		fmt.Println("EXIT WITH WRONG CODE:", code)
	}()
	// No need to call Close() when using pipes ?
	// p.stdout.Close()
	// p.stderr.Close()
}

func (p *Process) prepare() {
	p.Status = STARTING
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

// Change current working directory if path exists and is a directory
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
	}
	for s.Scan() {
		if p.IsForeground {
			fmt.Fprintln(os.Stdout, s.Text())
		}
		fmt.Fprintln(file, s.Text())
	}
	// When stream is closed this will executed
	which := "stdout"
	if stream == p.stderr {
		which = "stderr"
	}
	fmt.Println(p.Name, "writing", which, "stopped")
}

// SetForeground ...
func (p *Process) SetForeground(val bool) {
	p.IsForeground = val
}

// TODO:
func (p *Process) launchPool() {
}
