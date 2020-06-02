package cli

import (
	"github.com/tuommii/taskmaster/job"
)

type runnable func(tasks map[string]*job.Process, arg string) string

// Commands hold all commands
var Commands = map[string]runnable{
	// Used for autocomplete
	"job_names": suggestions,
	"help":      help,
	"h":         help,
	"status":    status,
	"st":        status,
	"reload":    nil,
	"start":     start,
	"run":       start,
	"stop":      stop,
	"exit":      nil,
	"quit":      nil,
	"fg":        fg,
	"bg":        bg,
}

// CommandNames ...
func CommandNames() []string {
	var names []string
	for key := range Commands {
		names = append(names, key)
	}
	return names
}

func suggestions(tasks map[string]*job.Process, arg string) string {
	var names string
	for name := range tasks {
		names += name + "|"
	}
	return names[:len(names)-1]
}

func help(tasks map[string]*job.Process, arg string) string {
	return "HELP command!"
}

func start(tasks map[string]*job.Process, arg string) string {
	if task, found := tasks[arg]; found {
		task.Launch(false)
		return arg + " STARTED"
	}
	return arg + " NOT FOUND"
}

func stop(tasks map[string]*job.Process, arg string) string {
	if task, found := tasks[arg]; found {
		task.Kill()
		return arg + " STOPPED"
	}
	return arg + " NOT FOUND"
}

func status(tasks map[string]*job.Process, arg string) string {
	var res string
	for name, task := range tasks {
		res += name + " " + task.Status + "\n"
	}
	return res[:len(res)-1]
}

func fg(tasks map[string]*job.Process, arg string) string {
	if task, found := tasks[arg]; found {
		task.SetForeground(true)
		return "attached " + arg + " output to stdout"
	}
	return arg + " NOT FOUND"
}

func bg(tasks map[string]*job.Process, arg string) string {
	if task, found := tasks[arg]; found {
		task.SetForeground(false)
		return "deattached " + arg + " output from stdout"
	}
	return arg + " NOT FOUND"
}
