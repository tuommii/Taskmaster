package job

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/tuommii/taskmaster/logger"
)

// Statuses
const (
	LOADED   = "LOADED"
	STARTING = "STARTING"
	RUNNING  = "RUNNING"
	STOPPED  = "STOPPED"
	FAILED   = "FAILED"
)

// Process represents runnable process
type Process struct {
	options
	IsForeground bool
	Instances    map[string]*Process
	Cmd          *exec.Cmd
	Started      time.Time
	Status       string
	Stdout       io.ReadCloser
	Stderr       io.ReadCloser
}

// Launch executes a task
func (p *Process) Launch(launchAutostartOnly bool) error {
	if err := p.validateLauch(launchAutostartOnly); err != nil {
		return err
	}
	p.prepare()
	oldMask := syscall.Umask(p.Umask)
	if err := p.execute(); err != nil {
		logger.Error(p.Name, p.Status, err)
		return err
	}
	p.killAfter()
	syscall.Umask(oldMask)
	p.clean()
	return nil
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

// Kill process
func (p *Process) Kill() error {
	if p.Status != RUNNING {
		return errors.New(p.Name + " wasn't running")
	}
	p.Status = STOPPED

	sig := killSignals[p.StopSignal]
	if sig == 0 {
		sig = syscall.SIGTERM
	}
	logger.Info("Killing with", p.StopSignal, sig)
	return p.Cmd.Process.Signal(sig)
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

func (p *Process) relaunch(err error) error {
	p.Status = FAILED
	p.Retries--
	if p.Retries > 0 && p.Retries < maxRetries {
		logger.Info("Trying launch", p.Name, "again...")
		p.execute()
	}
	return err
}

func (p *Process) setStarted() {
	// No delay
	if p.StartTime <= 0 {
		p.Status = RUNNING
		logger.Info(p.Name, p.Status)
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
		logger.Info(p.Name, "is consired started", p.Status)
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
		logger.Info(p.Name, "stopped")
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
	// if p.Status != RUNNING {
	// 	return
	// }
	go func() {
		err := p.Cmd.Wait()
		if err == nil {
			return
		}
		p.Status = STOPPED
		code := p.Cmd.ProcessState.ExitCode()
		if p.properExit(code) || code == -1 {
			logger.Info("Exited with proper code", code)
			if p.AutoRestart == "true" {
				logger.Info(p.Name, "Restarting...")
				p.Launch(false)
			}
			return
		}
		logger.Warning(p.Name, "Exited with wrong code:", code)
		if (p.AutoRestart == "unexpected" || p.AutoRestart == "true") && code != 127 {
			logger.Info(p.Name, "Restarting...")
			p.Launch(false)
		}
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
	p.Stdout, err = p.Cmd.StdoutPipe()
	if err != nil {
		logger.Error(err)
	}

	p.Stderr, err = p.Cmd.StderrPipe()
	if err != nil {
		logger.Error(err)
	}

	p.Cmd.Env = p.Env
	p.cwd(p.WorkingDir)
	go p.redirect(p.Stdout, p.OutputLog, os.Stdout)
	go p.redirect(p.Stderr, p.ErrorLog, os.Stderr)
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
	if stream == p.Stderr {
		which = "stderr"
	}
	p.Status = STOPPED
	logger.Info(p.Name, "writing", which, "stopped")
}

// SetForeground ...
func (p *Process) SetForeground(val bool) {
	p.IsForeground = val
}
