package cli

import (
	"fmt"
	"sort"
	"time"

	"github.com/tuommii/taskmaster/job"
)

type runnable func(tasks map[string]*job.Process, arg string) string

// GetJobsCommand returns currently available jobs, needed for autocompletion
const GetJobsCommand = "suggestions"

// Command holds function and help message
type Command struct {
	Runnable runnable
	Help     string
}

// Commands holds all commands
var Commands = map[string]*Command{
	// Used for autocomplete
	GetJobsCommand: {
		Runnable: suggestions,
		Help:     "For internal use",
	},
	"help": {
		Runnable: help,
		Help:     "HELP",
	},
	"status": {
		Runnable: status,
		Help:     "STATUS",
	},
	"restart": {
		Runnable: restart,
		Help:     "RESTART",
	},
	"reload": {
		Runnable: nil,
		Help:     "RELOAD",
	},
	"start": {
		Runnable: start,
		Help:     "START",
	},
	"stop": {
		Runnable: stop,
		Help:     "STOP",
	},
	"uptime": {
		Runnable: uptime,
		Help:     "UPTIME",
	},
	"exit": nil,
	"quit": nil,
	"fg": {
		Runnable: fg,
		Help:     "FG",
	},
	"bg": {
		Runnable: bg,
		Help:     "BG",
	},
}

func init() {
	// Aliases
	Commands["run"] = Commands["start"]
	Commands["st"] = Commands["status"]
	Commands["h"] = Commands["help"]
}

var notFound = " not found"

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
	task, found := tasks[arg]
	if !found {
		return arg + notFound
	}
	err := task.Launch(false)
	if err != nil {
		return err.Error()
	}
	return arg + " STARTED"
}

func stop(tasks map[string]*job.Process, arg string) string {
	task, found := tasks[arg]
	if !found {
		return arg + notFound
	}
	fmt.Println("FOUND", task.Name)
	task.Kill()
	return arg + " STOPPED"
}

func status(tasks map[string]*job.Process, arg string) string {
	var res string
	for name, task := range tasks {
		res += name + " " + task.Status + "\n"
	}
	return res[:len(res)-1]
}

func fg(tasks map[string]*job.Process, arg string) string {
	task, found := tasks[arg]
	if !found {
		return arg + notFound
	}

	task.SetForeground(true)
	return "attached " + arg + " output to stdout"
}

func bg(tasks map[string]*job.Process, arg string) string {
	task, found := tasks[arg]
	if !found {
		return arg + notFound
	}
	task.SetForeground(false)
	return "deattached " + arg + " output from stdout"
}

func restart(tasks map[string]*job.Process, arg string) string {
	task, found := tasks[arg]
	if !found {
		return arg + notFound
	}
	task.Kill()
	task.Launch(false)
	return arg + " RESTARTED"
}

func uptime(tasks map[string]*job.Process, arg string) string {
	if task, found := tasks[arg]; found && !task.Started.IsZero() {
		return arg + " " + time.Since(task.Started).String()
	}
	return arg + notFound
}

// CommandNames returns all command names except command(s) for internal
func CommandNames() []string {
	var names []string
	for key := range Commands {
		if key == GetJobsCommand {
			continue
		}
		names = append(names, key)
	}
	sort.Strings(names)
	return names
}
