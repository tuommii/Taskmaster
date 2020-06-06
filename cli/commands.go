package cli

import (
	"sort"
	"time"

	"github.com/tuommii/taskmaster/job"
)

type runnable func(tasks map[string]*job.Process, arg string) string

// Command ...
type Command struct {
	Runnable runnable
	Help     string
}

// Commands ...
var Commands = map[string]*Command{
	// Used for autocomplete
	"job_names": {
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
	"reload": nil,
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
	Commands["run"] = Commands["start"]
	Commands["st"] = Commands["status"]
	Commands["h"] = Commands["help"]
}

var notFound = " not found"

// CommandNames ...
func CommandNames() []string {
	var names []string
	for key := range Commands {
		names = append(names, key)
	}
	sort.Strings(names)
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
