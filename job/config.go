package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
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
	Procs int `json:"instances"`
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

var validators = []func(*Process) bool{
	func(p *Process) bool { return p.validateName() },
	func(p *Process) bool { return p.validateStartTime() },
	func(p *Process) bool { return p.validateStopTime() },
	func(p *Process) bool { return p.validateProcs() },
	func(p *Process) bool { return p.validateRetries() },
}

// LoadAll loads all jobs from config file
func LoadAll(path string) map[string]*Process {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error while opening config file: ", err)
	}
	tasks := make(map[string]*Process)
	err = json.Unmarshal([]byte(file), &tasks)
	if err != nil {
		fmt.Println("Error while loading config file", err)
		os.Exit(1)
	}
	initTasks(tasks)
	return tasks
}

func initTasks(tasks map[string]*Process) {
	copies := make(map[string]*Process)
	for name, task := range tasks {
		// Names are keys in json-file so they must be set
		task.Name = name
		task.Status = LOADED
		if err := validateConfig(task); err != nil {
			fmt.Println(err)
			delete(tasks, name)
			continue
		}
		createCopies(task, copies)
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks given. Exiting...")
		os.Exit(1)
	}
	// merge
	for k, v := range copies {
		tasks[k] = v
	}
}

func createCopies(task *Process, copies map[string]*Process) {
	for i := 0; i < task.Procs-1; i++ {
		var copy Process
		copy = *task
		copy.Name += strconv.Itoa(i + 2)
		copies[copy.Name] = &copy
	}
}

func validateConfig(task *Process) error {
	for i := 0; i < len(validators); i++ {
		if fine := validators[i](task); !fine {
			return errors.New("Invalid config for: " + task.Name)
		}
	}
	return nil
}

func (p *Process) validateName() bool {
	nameLen := len(p.Name)
	if nameLen < 1 || nameLen > 32 || !alphaOnly(p.Name) {
		return false
	}
	return true
}

func (p *Process) validateStartTime() bool {
	return p.StartTime >= 0
}

func (p *Process) validateStopTime() bool {
	return p.StopTime >= 0
}

func (p *Process) validateProcs() bool {
	if p.Procs > 4 {
		return false
	}
	return true
}

func (p *Process) validateRetries() bool {
	if p.Retries > 4 {
		return false
	}
	return true
}

func alphaOnly(str string) bool {
	for _, c := range str {
		if !((c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')) {
			return false
		}
	}
	return true
}
