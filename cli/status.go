package cli

import (
	"fmt"

	"github.com/tuommii/taskmaster/job"
)

// StatusCmd implements help command
var StatusCmd = &Command{
	Name:  "status",
	Usage: "Show status",
	Alias: "st",
	Run:   status,
}

// Help ...
func status(cmd *Command, args []string, tasks map[string]*job.Process) {
	fmt.Println("Listing job statuses")
	for _, task := range tasks {
		fmt.Println(task.Name, task.Status)
	}
}
